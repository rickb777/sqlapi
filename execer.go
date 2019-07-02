package sqlapi

import (
	"context"
	"database/sql"
	"time"
)

type Logger interface {
	Log(format string, v ...interface{})
	LogT(msg string, startTime *time.Time, data ...interface{})
	LogQuery(query string, args ...interface{})
	LogIfError(err error) error
	LogError(err error) error
	TraceLogging(on bool)
}

// Execer is a precis of *sql.DB and *sql.Tx.
// See database/sql.
type Execer interface {
	// ExecContext executes a query without returning any rows.
	// The arguments are for any placeholder parameters in the query.
	ExecContext(ctx context.Context, query string, arguments ...interface{}) (int64, error)

	// InsertContext executes a query and returns the insertion ID.
	// The arguments are for any placeholder parameters in the query.
	InsertContext(ctx context.Context, query string, arguments ...interface{}) (int64, error)

	// PrepareContext creates a prepared statement for later queries or executions.
	// Multiple queries or executions may be run concurrently from the
	// returned statement.
	// The caller must call the statement's Close method
	// when the statement is no longer needed.
	//
	// The provided context is used for the preparation of the statement, not for the
	// execution of the statement.
	PrepareContext(ctx context.Context, name, query string) (SqlStmt, error)

	// QueryContext executes a query that returns rows, typically a SELECT.
	// The arguments are for any placeholder parameters in the query.
	QueryContext(ctx context.Context, query string, arguments ...interface{}) (SqlRows, error)

	// QueryRowContext executes a query that is expected to return at most one row.
	// QueryRowContext always returns a non-nil value. Errors are deferred until
	// Row's Scan method is called.
	// If the query selects no rows, the *Row's Scan will return ErrNoRows.
	// Otherwise, the *Row's Scan scans the first selected row and discards
	// the rest.
	QueryRowContext(ctx context.Context, query string, arguments ...interface{}) SqlRow

	IsTx() bool
}

// SqlDB is able to make queries and begin transactions.
type SqlDB interface {
	Execer
	BeginTx(ctx context.Context, opts *sql.TxOptions) (SqlTx, error)
	Transact(ctx context.Context, txOptions *sql.TxOptions, fn func(SqlTx) error) error
	PingContext(ctx context.Context) error
	Stats() sql.DBStats
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
var _ SqlStmt = &sql.Stmt{}
