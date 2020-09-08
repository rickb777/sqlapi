package pgxapi

import (
	"context"
	"log"
	"net"
	"os"
	"strconv"
	"syscall"
	"time"

	"github.com/cenkalti/backoff/v3"
	"github.com/jackc/pgx"
	"github.com/pkg/errors"
	"github.com/rickb777/where/quote"
)

// MustConnectEnv is as per ConnectEnv but with a fatal termination on error.
func MustConnectEnv(lgr pgx.Logger, logLevel pgx.LogLevel) SqlDB {
	db, err := ConnectEnv(lgr, logLevel)
	if err != nil {
		log.Fatalf("%v\n", err)
	}

	return db
}

// ConnectEnv connects to the PostgreSQL server using environment variables:
// PGHOST, PGPORT, PGUSER, PGPASSWORD, PGDATABASE, PGCONNECT_TIMEOUT,
// PGSSLMODE, PGSSLKEY, PGSSLCERT, PGSSLROOTCERT.
// Also available are DB_MAX_CONNECTIONS, DB_CONNECT_DELAY and DB_CONNECT_TIMEOUT.
func ConnectEnv(lgr pgx.Logger, logLevel pgx.LogLevel) (SqlDB, error) {
	poolConfig := ParseEnvConfig()
	poolConfig.Logger = lgr
	poolConfig.LogLevel = logLevel
	return Connect(poolConfig)
}

// ParseEnvConfig creates connection pool config information based on environment variables:
// PGHOST, PGPORT, PGUSER, PGPASSWORD, PGDATABASE, PGCONNECT_TIMEOUT,
// PGSSLMODE, PGSSLKEY, PGSSLCERT, PGSSLROOTCERT.
// Also available are DB_MAX_CONNECTIONS, DB_CONNECT_DELAY and DB_CONNECT_TIMEOUT.
func ParseEnvConfig() pgx.ConnPoolConfig {
	setDefaultEnvValues()
	config, err := pgx.ParseEnvLibpq()
	if err != nil {
		log.Fatalf("Unable to parse environment: %v\n", err)
	}

	maxConnections, _ := strconv.Atoi(os.Getenv("DB_MAX_CONNECTIONS"))
	// if maxConnections == 0, dbx later changes this to its default (5)

	poolConfig := pgx.ConnPoolConfig{
		ConnConfig:     config,
		MaxConnections: maxConnections,
		AcquireTimeout: osGetEnvDuration("PGCONNECT_TIMEOUT", time.Second),
	}
	return poolConfig
}

func setDefaultEnvValues() {
	requireEnv("PGHOST", "localhost")
	requireEnv("PGPORT", "5432")
	requireEnv("PGDATABASE", "postgres")
	requireEnv("PGUSER", "postgres")
	requireEnv("PGPASSWORD", "psql")
	requireEnv("PGSSLMODE", "prefer")
	requireEnv("DB_MAX_CONNECTIONS", "100")
	requireEnv("DB_CONNECT_DELAY", "1ms")   // doesn't attempt connection until after this delay
	requireEnv("DB_CONNECT_TIMEOUT", "10s") // app aborts after this time; 0 is infinite
}

// MustConnect is as per Connect but with a fatal termination on error.
func MustConnect(config pgx.ConnPoolConfig) SqlDB {
	db, err := Connect(config)
	if err != nil {
		log.Fatalf("%v\n", err)
	}

	return db
}

// Connect opens a database connection and pings the server.
func Connect(config pgx.ConnPoolConfig) (SqlDB, error) {
	// set default behaviour to be unquoted identifiers (instead of ANSI SQL quote marks)
	quote.DefaultQuoter = quote.NoQuoter

	config.Logger.Log(pgx.LogLevelInfo, "DB connection",
		map[string]interface{}{"host": config.Host, "port": config.Port, "user": config.User, "database": config.Database},
	)

	pool, err := createConnectionPool(config.Logger, config)
	if err != nil {
		return nil, errors.Wrap(err, "Unable to connect to the database.\n")
	}

	// ping the connection using an empty statement
	_, err = pool.ExecEx(context.Background(), ";", nil)
	if err != nil {
		return nil, errors.Wrap(err, "Unable to communicate with to the database.\n")
	}

	return WrapDB(pool, config.Logger), nil
}

//-------------------------------------------------------------------------------------------------

const (
	// https://www.postgresql.org/docs/current/static/errcodes-appendix.html
	// This error code is seen sometimes when the database is starting.
	psqlCannotConnectNow = "57P03"
)

func createConnectionPool(lgr pgx.Logger, config pgx.ConnPoolConfig) (*pgx.ConnPool, error) {
	backOff := backoff.NewExponentialBackOff()
	backOff.MaxElapsedTime = osGetEnvDuration("DB_CONNECT_TIMEOUT", 0)

	var pool *pgx.ConnPool
	var err error

	dbConnectDelay := osGetEnvDuration("DB_CONNECT_DELAY", 0)
	if dbConnectDelay > 0 {
		lgr.Log(pgx.LogLevelInfo, "Waiting to connect to Postgres.", nil)
		time.Sleep(dbConnectDelay)
	}

	lgr.Log(pgx.LogLevelInfo, "Connecting to Postgres.", nil)

	// Construct a connection pool, retry until a connection pool can be established
	err = backoff.RetryNotify(
		func() error {
			pool, err = pgx.NewConnPool(config)
			if err != nil {
				e1, ok := err.(pgx.PgError)
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
		func(err error, next time.Duration) { notify(lgr, err, next) },
	)

	if err != nil {
		return nil, err
	}

	defer func() {
		if err != nil {
			pool.Close()
		}
	}()

	lgr.Log(pgx.LogLevelInfo, "Connected to Postgres successfully.", nil)
	return pool, nil
}

func notify(lgr pgx.Logger, err error, next time.Duration) {
	lgr.Log(pgx.LogLevelWarn, "Failed to create Postgres connection pool",
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

func osGetEnvDuration(name string, deflt time.Duration) time.Duration {
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
