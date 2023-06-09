package pgxapi

import (
	"context"
	"errors"
	"fmt"
	"log"
	"net"
	"os"
	"strconv"
	"strings"
	"time"

	"github.com/cenkalti/backoff/v4"
	"github.com/jackc/pgconn"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/jackc/pgx/v5/tracelog"
	"github.com/rickb777/where/quote"
)

// MustConnectEnv is as per ConnectEnv but with a fatal termination on error.
func MustConnectEnv(ctx context.Context, lgr tracelog.Logger, logLevel tracelog.LogLevel, tries int) SqlDB {
	db, err := ConnectEnv(ctx, lgr, logLevel, tries)
	if err != nil {
		log.Fatalf("%v\n", err)
	}

	return db
}

// ConnectEnv connects to the PostgreSQL server using environment variables:
// PGHOST, PGPORT, PGUSER, PGPASSWORD, PGDATABASE, PGCONNECT_TIMEOUT,
// PGSSLMODE, PGSSLKEY, PGSSLCERT, PGSSLROOTCERT.
// Also available are PGQUOTE, DB_MAX_CONNECTIONS, DB_CONNECT_DELAY and DB_CONNECT_TIMEOUT.
// Use PGQUOTE to set "ansi", "mysql" or "none" as the policy for quoting identifiers (the default
// is none).
func ConnectEnv(ctx context.Context, lgr tracelog.Logger, logLevel tracelog.LogLevel, tries int) (SqlDB, error) {
	poolConfig := ParseEnvConfig()
	poolConfig.ConnConfig.Tracer = &tracelog.TraceLog{Logger: lgr, LogLevel: logLevel}
	quoter := quote.PickQuoter(os.Getenv("PGQUOTE"))
	return Connect(ctx, poolConfig, quoter, tries)
}

// ParseEnvConfig creates connection pool config information based on environment variables:
// PGHOST, PGPORT, PGUSER, PGPASSWORD, PGDATABASE, PGCONNECT_TIMEOUT,
// PGSSLMODE, PGSSLKEY, PGSSLCERT, PGSSLROOTCERT.
// Also available are DB_URL, DB_MAX_CONNECTIONS, DB_CONNECT_DELAY and DB_CONNECT_TIMEOUT.
func ParseEnvConfig() *pgxpool.Config {
	dbUrl := os.Getenv("DB_URL")
	if dbUrl != "" && !strings.HasPrefix("postgres://", dbUrl) {
		dbUrl = "postgres://" + dbUrl
	}
	// conveniently, if dbUrl is blank, ParseConfig checks the 'standard' environment variables
	// including PGHOST, PGPORT, PGDATABASE, PGUSER, etc
	// (https://pkg.go.dev/github.com/jackc/pgconn#ParseConfig)
	config, err := pgxpool.ParseConfig(dbUrl)
	if err != nil {
		log.Fatalf("Unable to parse environment: %v\n", err)
	}

	return config
}

// MustConnect is as per Connect but with a fatal termination on error.
func MustConnect(ctx context.Context, config *pgxpool.Config, quoter quote.Quoter, tries int) SqlDB {
	db, err := Connect(ctx, config, quoter, tries)
	if err != nil {
		log.Fatalf("%v\n", err)
	}

	return db
}

// Connect opens a database connection and pings the server.
// If the connection fails, it is retried using an exponential backoff.
// the maximum number of (re-)tries can be specified; if this is zero, there is no limit.
func Connect(ctx context.Context, config *pgxpool.Config, quoter quote.Quoter, tries int) (SqlDB, error) {
	logger := config.ConnConfig.Tracer.(*tracelog.TraceLog).Logger
	logger.Log(ctx, tracelog.LogLevelInfo, "DB connection",
		map[string]interface{}{
			"host":     config.ConnConfig.Host,
			"port":     config.ConnConfig.Port,
			"user":     config.ConnConfig.User,
			"database": config.ConnConfig.Database},
	)

	pool, err := createConnectionPool(ctx, logger, config, tries)
	if err != nil {
		return nil, fmt.Errorf("%w - unable to connect to the database.", err)
	}

	// ping the connection using an empty statement
	_, err = pool.Exec(ctx, ";")
	if err != nil {
		return nil, fmt.Errorf("%w - unable to communicate with to the database.", err)
	}

	return WrapDB(pool, logger, quoter), nil
}

//-------------------------------------------------------------------------------------------------

const (
	// https://www.postgresql.org/docs/current/static/errcodes-appendix.html
	// This error code is seen sometimes when the database is starting.
	psqlCannotConnectNow = "57P03"
)

func createConnectionPool(ctx context.Context, lgr tracelog.Logger, config *pgxpool.Config, tries int) (*pgxpool.Pool, error) {
	backOff := backoff.NewExponentialBackOff()
	backOff.MaxElapsedTime = osGetenvDuration("DB_CONNECT_TIMEOUT", 0)

	var pool *pgxpool.Pool
	var err error

	dbConnectDelay := osGetenvDuration("DB_CONNECT_DELAY", 0)
	if dbConnectDelay > 0 {
		lgr.Log(ctx, tracelog.LogLevelInfo, "Waiting to connect to Postgres.", nil)
		time.Sleep(dbConnectDelay)
	}

	lgr.Log(ctx, tracelog.LogLevelInfo, "Connecting to Postgres.", nil)

	// Construct a connection pool, retry until a connection pool can be established
	err = backoff.RetryNotify(
		func() error {
			tries--
			pool, err = pgxpool.NewWithConfig(ctx, config)
			if err != nil {
				if tries == 0 {
					return backoff.Permanent(err) // no more tries
				}

				e0 := errors.Unwrap(err)
				if e0 == nil {
					e0 = err
				}
				e1, ok := e0.(*pgconn.PgError)
				if ok && e1.Code == psqlCannotConnectNow {
					// Retry connection on this specific error
					return e1
				}
				e2, ok := e0.(*net.OpError)
				if ok && e2.Op == "dial" {
					// Retry connection on this specific error
					return e2
				}
				// Other errors are considered permanent
				return backoff.Permanent(err)
			}
			return nil
		},
		backOff,
		func(err error, next time.Duration) { notify(ctx, lgr, err, next) },
	)

	if err != nil {
		return nil, err
	}

	defer func() {
		if err != nil {
			pool.Close()
		}
	}()

	lgr.Log(ctx, tracelog.LogLevelInfo, "Connected to Postgres successfully.", nil)
	return pool, nil
}

func notify(ctx context.Context, lgr tracelog.Logger, err error, next time.Duration) {
	lgr.Log(ctx, tracelog.LogLevelWarn, "Failed to create Postgres connection pool",
		map[string]interface{}{
			"error":    err,
			"retry_in": next.Truncate(time.Millisecond),
		})
}

func osGetenvDuration(name string, deflt time.Duration) time.Duration {
	duration := os.Getenv(name)
	d, err := time.ParseDuration(duration)
	if err == nil {
		return d
	}

	// pgx allows some durations to be provided as plain integers (i.e. seconds)
	s, err := strconv.Atoi(duration)
	if err == nil {
		return time.Duration(s) * time.Second
	}

	return deflt
}
