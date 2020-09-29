package sqlapi

import (
	"context"
	"database/sql"
	"log"
	"net"
	"os"
	"strconv"
	"time"

	"github.com/rickb777/where/quote"

	"github.com/cenkalti/backoff/v3"
	"github.com/jackc/pgx/v4"
	"github.com/rickb777/sqlapi/dialect"
)

// MustConnectEnv is as per ConnectEnv but with a fatal termination on error.
func MustConnectEnv(ctx context.Context, lgr pgx.Logger, logLevel pgx.LogLevel) SqlDB {
	db, err := ConnectEnv(ctx, lgr, logLevel)
	if err != nil {
		log.Fatalf("%v\n", err)
	}

	return db
}

// ConnectEnv connects to the database server using environment variables:
// DB_URL, DB_DRIVER, DB_DIALECT and DB_QUOTE.
// Also available are DB_MAX_CONNECTIONS, DB_CONNECT_DELAY and DB_CONNECT_TIMEOUT.
// Use DB_QUOTE to set "ansi", "mysql" or "none" as the policy for quoting identifiers (the default
// is none).
func ConnectEnv(ctx context.Context, lgr pgx.Logger, logLevel pgx.LogLevel) (SqlDB, error) {
	dbUrl := os.Getenv("DB_URL")
	driver := os.Getenv("DB_DRIVER")
	if driver == "" {
		driver = "sqlite3"
	}

	di := dialect.PickDialect(os.Getenv("DB_DIALECT"))
	if di == nil {
		di = dialect.Sqlite
	}

	quoter := quote.PickQuoter(os.Getenv("DB_QUOTE"))
	if quoter != nil {
		di = di.WithQuoter(quoter)
	}

	logger := NewLogger(lgr)
	return Connect(ctx, driver, dbUrl, di, logger)
}

// MustConnect is as per Connect but with a fatal termination on error.
func MustConnect(ctx context.Context, driverName, dataSourceName string, di dialect.Dialect, lgr Logger) SqlDB {
	db, err := Connect(ctx, driverName, dataSourceName, di, lgr)
	if err != nil {
		log.Fatalf("%v\n", err)
	}

	return db
}

func Connect(ctx context.Context, driverName, dataSourceName string, di dialect.Dialect, lgr Logger) (SqlDB, error) {
	backOff := backoff.NewExponentialBackOff()
	backOff.MaxElapsedTime = osGetEnvDuration("DB_CONNECT_TIMEOUT", 0)

	dbConnectDelay := osGetEnvDuration("DB_CONNECT_DELAY", 0)
	if dbConnectDelay > 0 {
		lgr.Log(ctx, pgx.LogLevelInfo, "Waiting to connect to "+di.String(), nil)
		time.Sleep(dbConnectDelay)
	}

	var db *sql.DB
	var err error

	// Construct a connection pool, retry until a connection pool can be established
	err = backoff.RetryNotify(
		func() error {
			lgr.Log(ctx, pgx.LogLevelInfo, "Connecting to "+di.String(), nil)
			db, err = sql.Open(driverName, dataSourceName)
			if err != nil {
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
			db.Close()
		}
	}()

	lgr.Log(ctx, pgx.LogLevelInfo, "Connected successfully to "+di.String(), nil)

	return WrapDB(db, di, lgr), nil
}

func notify(ctx context.Context, lgr Logger, err error, next time.Duration) {
	lgr.Log(ctx, pgx.LogLevelWarn, "Failed to open DB connection",
		map[string]interface{}{
			"error":    err,
			"retry_in": next.Truncate(time.Millisecond),
		})
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
