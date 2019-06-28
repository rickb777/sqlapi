package pgxapi

import (
	"context"
	"log"
	"net"
	"os"
	"syscall"
	"time"

	"github.com/cenkalti/backoff"
	"github.com/jackc/pgx"
	"github.com/pkg/errors"
	"github.com/rickb777/where/quote"
)

func ConnectEnv(lgr pgx.Logger, logLevel pgx.LogLevel) SqlDB {
	setDefaultEnvValues()
	config, err := pgx.ParseEnvLibpq()
	if err != nil {
		log.Fatalf("Unable to parse environment: %v\n", err)
	}

	return Connect(config, lgr, logLevel)
}

func setDefaultEnvValues() {
	requireEnv("PGHOST", "localhost")
	requireEnv("PGPORT", "5432")
	requireEnv("PGDATABASE", "postgres")
	requireEnv("PGUSER", "postgres")
	requireEnv("PGPASSWORD", "psql")
	requireEnv("PGSSLMODE", "prefer")
	requireEnv("DB_CONNECT_DELAY", "1ms") // doesn't attempt connection until after this delay
	requireEnv("DB_CONNECT_TIMEOUT", "0") // app aborts after this time; 0 is infinite
}

func Connect(config pgx.ConnConfig, lgr pgx.Logger, logLevel pgx.LogLevel) SqlDB {
	db, err := doConnect(config, lgr, logLevel)
	if err != nil {
		log.Fatalf("%v\n", err)
	}
	return db
}

// Connect opens a database connection and pings the server. Any failure is fatal (i.e.
// the app terminates).
func doConnect(config pgx.ConnConfig, lgr pgx.Logger, logLevel pgx.LogLevel) (SqlDB, error) {
	// set default behaviour to be unquoted identifiers (instead of ANSI SQL quote marks)
	quote.DefaultQuoter = quote.NoQuoter

	configCopy := config
	configCopy.Password = "***"
	lgr.Log(pgx.LogLevelInfo, "DB connection", map[string]interface{}{"dbconfig": configCopy})

	config.Logger = lgr
	config.LogLevel = logLevel

	pool, err := createConnectionPool(lgr, config, 10, time.Second)
	if err != nil {
		return nil, errors.Wrap(err, "Unable to connect to the database.\n")
	}

	// ping the connection using an empty statement
	_, err = pool.ExecEx(context.Background(), ";", nil)
	if err != nil {
		return nil, errors.Wrap(err, "Unable to communicate with to the database.\n")
	}

	return &shim{ex: pool, lgr: &toggleLogger{lgr: lgr, enabled: 1}, isTx: false}, nil
}

//-------------------------------------------------------------------------------------------------

const (
	// https://www.postgresql.org/docs/current/static/errcodes-appendix.html
	// This error code is seen sometimes when the database is starting.
	psqlCannotConnectNow = "57P03"
)

func createConnectionPool(lgr pgx.Logger, config pgx.ConnConfig, maxConnections int, connAcquireTimeout time.Duration) (*pgx.ConnPool, error) {
	backOff := backoff.NewExponentialBackOff()

	maxElapsedTime, _ := time.ParseDuration(os.Getenv("DB_CONNECT_TIMEOUT"))
	backOff.MaxElapsedTime = maxElapsedTime

	pgxConnPoolConfig := pgx.ConnPoolConfig{
		ConnConfig:     config,
		MaxConnections: maxConnections,
		AfterConnect:   nil,
		AcquireTimeout: connAcquireTimeout,
	}

	var pool *pgx.ConnPool
	var err error

	dbConnectDelay, _ := time.ParseDuration(os.Getenv("DB_CONNECT_DELAY"))
	if dbConnectDelay > 0 {
		lgr.Log(pgx.LogLevelInfo, "Waiting to connect to Postgres.", nil)
		time.Sleep(dbConnectDelay)
	}

	lgr.Log(pgx.LogLevelInfo, "Connecting to Postgres.", nil)

	// Construct a connection pool, retry until a connection pool can be established
	err = backoff.RetryNotify(
		func() error {
			pool, err = pgx.NewConnPool(pgxConnPoolConfig)
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
