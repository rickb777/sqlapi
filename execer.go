package sqlapi

import (
	"context"
	"database/sql"
	"io"
	"time"

	"github.com/rickb777/sqlapi/dialect"
)

// Logger provides the specialised logging operations within this API.
type Logger interface {
	Log(format string, v ...interface{})
	LogT(msg string, startTime *time.Time, data ...interface{})
	LogQuery(query string, args ...interface{})
	LogIfError(err error) error
	LogError(err error) error
	TraceLogging(on bool)
	SetOutput(out io.Writer)
}

// Getter provides the core methods for reading information from databases.
type Getter interface {
	// QueryContext executes a query that returns rows, typically a SELECT.
	// The arguments are for any placeholder parameters in the query.
	// Placeholders in the SQL are automatically replaced with numbered placeholders.
	QueryContext(ctx context.Context, sql string, arguments ...interface{}) (SqlRows, error)

	// QueryRowContext executes a query that is expected to return at most one row.
	// QueryRowContext always returns a non-nil value. Errors are deferred until
	// Row's Scan method is called.
	// If the query selects no rows, the *Row's Scan will return ErrNoRows.
	// Otherwise, the *Row's Scan scans the first selected row and discards
	// the rest.
	//
	// Placeholders in the SQL are automatically replaced with numbered placeholders.
	QueryRowContext(ctx context.Context, query string, arguments ...interface{}) SqlRow
}

// Execer is a precis of *sql.DB and *sql.Tx (see database/sql).
type Execer interface {
	Getter

	// ExecContext executes a query without returning any rows.
	// The arguments are for any placeholder parameters in the query.
	ExecContext(ctx context.Context, sql string, arguments ...interface{}) (int64, error)

	// InsertContext executes a query and returns the insertion ID.
	// The primary key column, pk, is used for some dialects, notably PostgreSQL.
	// The arguments are for any placeholder parameters in the query.
	InsertContext(ctx context.Context, pk, query string, arguments ...interface{}) (int64, error)

	// PrepareContext creates a prepared statement for later queries or executions.
	// Multiple queries or executions may be run concurrently from the
	// returned statement.
	// The caller must call the statement's Close method
	// when the statement is no longer needed.
	//
	// The provided context is used for the preparation of the statement, not for the
	// execution of the statement.
	PrepareContext(ctx context.Context, name, sql string) (SqlStmt, error)

	IsTx() bool

	// Logger gets the trace logger. Note that you can use this to rotate the output writer
	// via its SetOutput method. Also, it can even disable it completely (via ioutil.Discard).
	Logger() Logger

	// Dialect gets the database dialect.
	Dialect() dialect.Dialect
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
	Transact(ctx context.Context, txOptions *sql.TxOptions, fn func(SqlTx) error) error

	// PingContext tests connectivity to the database server.
	PingContext(ctx context.Context) error

	// Stats gets statistics from the database server.
	Stats() DBStats

	// SingleConn takes exclusive use of a connection for use by the supplied function.
	// The connection will be automatically released after the function has terminated.
	SingleConn(ctx context.Context, fn func(ex Execer) error) error

	// Close closes the database connection.
	Close() error
}

// SqlTx is a precis of *sql.Tx
type SqlTx interface {
	Execer
	Commit() error
	Rollback() error
}

// SqlStmt is a precis of *sql.Stmt
type SqlStmt interface {
	// ExecContext executes a query without returning any rows.
	// The arguments are for any placeholder parameters in the query.
	ExecContext(ctx context.Context, arguments ...interface{}) (sql.Result, error)

	// QueryContext executes a query that returns rows, typically a SELECT.
	// The arguments are for any placeholder parameters in the query.
	QueryContext(ctx context.Context, arguments ...interface{}) (*sql.Rows, error)

	// QueryRowContext executes a query that is expected to return at most one row.
	// QueryRowContext always returns a non-nil value. Errors are deferred until
	// Row's Scan method is called.
	// If the query selects no rows, the *Row's Scan will return ErrNoRows.
	// Otherwise, the *Row's Scan scans the first selected row and discards
	// the rest.
	QueryRowContext(ctx context.Context, arguments ...interface{}) *sql.Row

	Close() error
}

// Type conformance assertions
var _ SqlStmt = new(sql.Stmt)
