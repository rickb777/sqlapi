package pgxapi

import (
	"context"
	"fmt"
	"log"
	"net"
	"os"
	"strconv"
	"strings"
	"syscall"
	"time"

	"github.com/jackc/pgconn"
	"github.com/jackc/pgx/v4/pgxpool"

	"github.com/cenkalti/backoff/v3"
	"github.com/jackc/pgx/v4"
	"github.com/rickb777/where/quote"
)

// MustConnectEnv is as per ConnectEnv but with a fatal termination on error.
func MustConnectEnv(ctx context.Context, lgr pgx.Logger, logLevel pgx.LogLevel) SqlDB {
	db, err := ConnectEnv(ctx, lgr, logLevel)
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
func ConnectEnv(ctx context.Context, lgr pgx.Logger, logLevel pgx.LogLevel) (SqlDB, error) {
	poolConfig := ParseEnvConfig()
	poolConfig.ConnConfig.Logger = lgr
	poolConfig.ConnConfig.LogLevel = logLevel
	quoter := quote.PickQuoter(os.Getenv("PGQUOTE"))
	return Connect(ctx, poolConfig, quoter)
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
func MustConnect(ctx context.Context, config *pgxpool.Config, quoter quote.Quoter) SqlDB {
	db, err := Connect(ctx, config, quoter)
	if err != nil {
		log.Fatalf("%v\n", err)
	}

	return db
}

// Connect opens a database connection and pings the server.
func Connect(ctx context.Context, config *pgxpool.Config, quoter quote.Quoter) (SqlDB, error) {
	config.ConnConfig.Logger.Log(ctx, pgx.LogLevelInfo, "DB connection",
		map[string]interface{}{
			"host":     config.ConnConfig.Host,
			"port":     config.ConnConfig.Port,
			"user":     config.ConnConfig.User,
			"database": config.ConnConfig.Database},
	)

	pool, err := createConnectionPool(ctx, config.ConnConfig.Logger, config)
	if err != nil {
		return nil, fmt.Errorf("%w - unable to connect to the database.", err)
	}

	// ping the connection using an empty statement
	_, err = pool.Exec(ctx, ";")
	if err != nil {
		return nil, fmt.Errorf("%w - unable to communicate with to the database.", err)
	}

	return WrapDB(pool, config.ConnConfig.Logger, quoter), nil
}

//-------------------------------------------------------------------------------------------------

const (
	// https://www.postgresql.org/docs/current/static/errcodes-appendix.html
	// This error code is seen sometimes when the database is starting.
	psqlCannotConnectNow = "57P03"
)

func createConnectionPool(ctx context.Context, lgr pgx.Logger, config *pgxpool.Config) (*pgxpool.Pool, error) {
	backOff := backoff.NewExponentialBackOff()
	backOff.MaxElapsedTime = osGetenvDuration("DB_CONNECT_TIMEOUT", 0)

	var pool *pgxpool.Pool
	var err error

	dbConnectDelay := osGetenvDuration("DB_CONNECT_DELAY", 0)
	if dbConnectDelay > 0 {
		lgr.Log(ctx, pgx.LogLevelInfo, "Waiting to connect to Postgres.", nil)
		time.Sleep(dbConnectDelay)
	}

	lgr.Log(ctx, pgx.LogLevelInfo, "Connecting to Postgres.", nil)

	// Construct a connection pool, retry until a connection pool can be established
	err = backoff.RetryNotify(
		func() error {
			pool, err = pgxpool.ConnectConfig(ctx, config)
			if err != nil {
				e1, ok := err.(*pgconn.PgError)
				if ok && e1.Code == psqlCannotConnectNow {
					// Retry connection on this specific error
					return e1
				}
				e2, ok := err.(*net.OpError)
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

	lgr.Log(ctx, pgx.LogLevelInfo, "Connected to Postgres successfully.", nil)
	return pool, nil
}

func notify(ctx context.Context, lgr pgx.Logger, err error, next time.Duration) {
	lgr.Log(ctx, pgx.LogLevelWarn, "Failed to create Postgres connection pool",
		map[string]interface{}{
			"error":    err,
			"retry_in": next.Truncate(time.Millisecond),
		})
}

func requireEnv(name, defaultValue string) {
	_, exists := syscall.Getenv(name)
	if !exists {
		err := syscall.Setenv(name, defaultValue)
		if err != nil {
			log.Fatalf("Failed to set %q=%q; %v\n", name, defaultValue, err)
		}
	}
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
