package sqlapi

import (
	"context"
	"database/sql"
	"fmt"
	"log"
	"net"
	"os"
	"strconv"
	"strings"
	"time"

	"github.com/cenkalti/backoff/v4"
	"github.com/jackc/pgx/v4"
	"github.com/rickb777/sqlapi/driver"
	"github.com/rickb777/where/quote"
)

// MustConnectEnv is as per ConnectEnv but with a fatal termination on error.
func MustConnectEnv(ctx context.Context, lgr pgx.Logger, logLevel pgx.LogLevel, tries int) SqlDB {
	db, err := ConnectEnv(ctx, lgr, logLevel, tries)
	if err != nil {
		log.Fatalf("%v\n", err)
	}

	return db
}

const sqliteInMemory = "file::memory:?mode=memory&cache=shared"

// ConnectEnv connects to the database server using environment variables:
// DB_URL, DB_DRIVER and DB_QUOTE.
// Also available are DB_MAX_CONNECTIONS, DB_CONNECT_DELAY and DB_CONNECT_TIMEOUT.
// Use DB_QUOTE to set "ansi", "mysql" or "none" as the policy for quoting identifiers (the default
// is none).
func ConnectEnv(ctx context.Context, lgr pgx.Logger, logLevel pgx.LogLevel, tries int) (SqlDB, error) {
	dbUrl := os.Getenv("DB_URL")
	dbDriver := os.Getenv("DB_DRIVER")
	if dbDriver == "" {
		dbDriver = "sqlite3"
		if dbUrl == "" {
			dbUrl = sqliteInMemory
		}
	}

	dbUser := os.Getenv("DB_USER")
	dbName := os.Getenv("DB_NAME")
	if dbUrl == "" {
		dbPassword := os.Getenv("DB_PASSWORD")
		if dbUser != "" && dbPassword != "" && dbName != "" {
			dbUrl = fmt.Sprintf("%s:%s@/%s", dbUser, dbPassword, dbName)
		}
	}

	di := driver.PickDialect(dbDriver)
	if di == nil {
		di = driver.Sqlite()
	}

	quoter := quote.PickQuoter(os.Getenv("DB_QUOTE"))
	if quoter != nil {
		di = di.WithQuoter(quoter)
	}

	lgr.Log(ctx, logLevel, "Connecting to DB", map[string]interface{}{
		"url":      dbUrl,
		"driver":   dbDriver,
		"dialect":  di,
		"quote":    quoter,
		"user":     dbUser,
		"database": dbName,
	})
	logger := NewLogger(lgr)

	return Connect(ctx, dbDriver, dbUrl, di, logger, tries)
}

// MustConnect is as per Connect but with a fatal termination on error.
func MustConnect(ctx context.Context, driverName, dataSourceName string, di driver.Dialect, lgr Logger, tries int) SqlDB {
	db, err := Connect(ctx, driverName, dataSourceName, di, lgr, tries)
	if err != nil {
		log.Fatalf("%v\n", err)
	}

	return db
}

// Connect opens a database connection and pings the server.
// If the connection fails, it is retried using an exponential backoff.
// the maximum number of (re-)tries can be specified; if this is zero, there is no limit.
func Connect(ctx context.Context, driver, dsn string, di driver.Dialect, lgr Logger, tries int) (SqlDB, error) {
	if dsn == "" {
		return nil, fmt.Errorf("DB connect to %s failed: DSN is blank", driver)
	}

	if (driver == "postgres" || driver == "pgx") && dsn != "" && !strings.HasPrefix("postgres://", dsn) {
		dsn = "postgres://" + dsn
	}

	backOff := backoff.NewExponentialBackOff()
	backOff.MaxElapsedTime = osGetEnvDuration("DB_CONNECT_TIMEOUT", 0)

	dbConnectDelay := osGetEnvDuration("DB_CONNECT_DELAY", 0)
	if dbConnectDelay > 0 {
		lgr.Log(ctx, pgx.LogLevelInfo, "Waiting to connect to "+di.String(), nil)
		time.Sleep(dbConnectDelay)
	}

	var db *sql.DB
	var err error

	info := map[string]interface{}{
		"url":     dsn,
		"driver":  driver,
		"dialect": di,
		"tries":   tries,
	}

	// Construct a connection pool, retry until a connection pool can be established
	err = backoff.RetryNotify(
		func() error {
			tries--
			lgr.Log(ctx, pgx.LogLevelInfo, "Connecting to Docker DB", info)
			db, err = sql.Open(driver, dsn)
			if err != nil {
				if tries == 0 {
					return backoff.Permanent(err) // no more tries
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
		return nil, fmt.Errorf("%w - unable to connect to the database.", err)
	}

	defer func() {
		if err != nil {
			db.Close()
		}
	}()

	// ping the connection using an empty statement
	_, err = db.ExecContext(ctx, ";")
	if err != nil {
		return nil, fmt.Errorf("%w - unable to communicate with to the database.", err)
	}

	lgr.Log(ctx, pgx.LogLevelInfo, "Connected successfully to "+driver, nil)

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
