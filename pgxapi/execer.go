package pgxapi

import (
	"context"
	"io"
	"time"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/tracelog"
	"github.com/rickb777/sqlapi/driver"
)

// Logger provides the specialised logging operations within this API.
type Logger interface {
	tracelog.Logger
	LogT(ctx context.Context, level tracelog.LogLevel, msg string, startTime *time.Time, data ...interface{})
	LogQuery(ctx context.Context, query string, args ...interface{})
	LogIfError(ctx context.Context, err error) error
	LogError(ctx context.Context, err error) error
	TraceLogging(on bool)
	SetOutput(out io.Writer) // no-op
}

// Getter provides the core methods for reading information from databases.
type Getter interface {
	// Query executes a query that returns rows, typically a SELECT.
	// The arguments are for any placeholder parameters in the query.
	// Placeholders in the SQL are automatically replaced with numbered placeholders.
	Query(ctx context.Context, sql string, arguments ...interface{}) (SqlRows, error)

	// QueryRowContext executes a query that is expected to return at most one row.
	// QueryRowContext always returns a non-nil value. Errors are deferred until
	// Row's Scan method is called.
	// If the query selects no rows, the *Row's Scan will return ErrNoRows.
	// Otherwise, the *Row's Scan scans the first selected row and discards
	// the rest.
	//
	// Placeholders in the SQL are automatically replaced with numbered placeholders.
	QueryRow(ctx context.Context, query string, arguments ...interface{}) SqlRow
}

// Execer is a precis of *pgx.ConnPool and *pgx.Tx.
type Execer interface {
	Getter

	// Exec executes a query without returning any rows.
	// The arguments are for any placeholder parameters in the query.
	Exec(ctx context.Context, sql string, arguments ...interface{}) (int64, error)

	// Insert executes a query and returns the insertion ID.
	// The primary key column, pk, is used for some dialects, notably PostgreSQL.
	// The arguments are for any placeholder parameters in the query.
	Insert(ctx context.Context, pk, query string, arguments ...interface{}) (int64, error)

	IsTx() bool

	// Logger gets the trace logger. Note that you can use this to rotate the output writer
	// via its SetOutput method. Also, it can even disable it completely (via ioutil.Discard).
	Logger() Logger

	// Dialect gets the database dialect.
	Dialect() driver.Dialect
}

//-------------------------------------------------------------------------------------------------

// SqlDB is able to make queries and begin transactions.
type SqlDB interface {
	Execer

	// Transact handles a transaction according to some function. If the function completes
	// without error, the transaction will be committed automatically. If there is an error
	// or a panic, the transaction will be rolled back automatically.
	//
	// The function fn should avoid using the original SqlDB; this is easily achieved by
	// using a named function instead of an anonymous closure.
	Transact(ctx context.Context, txOptions *pgx.TxOptions, fn func(SqlTx) error) error

	// PingContext tests connectivity to the database server.
	Ping(ctx context.Context) error

	// Stats gets statistics from the database server.
	Stats() DBStats

	// SingleConn takes exclusive use of a connection for use by the supplied function.
	// The connection will be automatically released after the function has terminated.
	SingleConn(ctx context.Context, fn func(ex Execer) error) error

	// Close closes the database connection.
	Close() error

	// With returns a modified SqlDB with a user-supplied item.
	With(wrapped interface{}) SqlDB

	// UserItem gets a user-supplied item associated with this DB.
	UserItem() interface{}
}

// SqlTx is a precis of *pgx.Tx
type SqlTx interface {
	Execer
	Commit(ctx context.Context) error
	Rollback(ctx context.Context) error
}
