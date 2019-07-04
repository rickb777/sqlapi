package pgxapi

import (
	"context"
	"database/sql"
	"github.com/jackc/pgx"
	"time"
)

type Getter interface {
	// QueryContext executes a query that returns rows, typically a SELECT.
	// The arguments are for any placeholder parameters in the query.
	// Placeholders in the SQL are automatically replaced with numbered placeholders.
	QueryContext(ctx context.Context, sql string, arguments ...interface{}) (SqlRows, error)

	// QueryExRaw directly accesses the pgx.QueryEx connection pool method, exposing
	// the query options. The placeholders will not be automatically replaced.
	QueryExRaw(ctx context.Context, sql string, options *pgx.QueryExOptions, arguments ...interface{}) (SqlRows, error)

	// QueryRowContext executes a query that is expected to return at most one row.
	// QueryRowContext always returns a non-nil value. Errors are deferred until
	// Row's Scan method is called.
	// If the query selects no rows, the *Row's Scan will return ErrNoRows.
	// Otherwise, the *Row's Scan scans the first selected row and discards
	// the rest.
	//
	// Placeholders in the SQL are automatically replaced with numbered placeholders.
	QueryRowContext(ctx context.Context, query string, arguments ...interface{}) SqlRow

	// QueryRowExRaw directly accesses the pgx.QueryRowEx connection pool method, exposing
	// the query options. The placeholders will not be automatically replaced.
	QueryRowExRaw(ctx context.Context, query string, options *pgx.QueryExOptions, arguments ...interface{}) SqlRow
}

type Batcher interface {
	// BeginBatch exposes the pgx batch operations.
	BeginBatch() *pgx.Batch
}

type Logger interface {
	pgx.Logger
	LogT(level pgx.LogLevel, msg string, startTime *time.Time, data ...interface{})
	LogQuery(query string, args ...interface{})
	LogIfError(err error) error
	LogError(err error) error
	TraceLogging(on bool)
}

// Execer is a precis of *pgx.ConnPool and *pgx.Tx.
type Execer interface {
	Getter
	Batcher

	// ExecContext executes a query without returning any rows.
	// The arguments are for any placeholder parameters in the query.
	ExecContext(ctx context.Context, sql string, arguments ...interface{}) (int64, error)

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
	PrepareContext(ctx context.Context, name, sql string) (*pgx.PreparedStatement, error)

	IsTx() bool
	Logger() Logger
}

//-------------------------------------------------------------------------------------------------

// SqlDB is able to make queries and begin transactions.
type SqlDB interface {
	Execer

	// Transact handles a transaction according to some function. If the function completes
	// without error, the transaction will be committed automatically. If there is an error
	// or a panic, the transaction will be rolled back automatically.
	Transact(ctx context.Context, txOptions *pgx.TxOptions, fn func(SqlTx) error) error

	// PingContext tests connectivity to the database server.
	PingContext(ctx context.Context) error

	// Stats gets statistics from the database server.
	Stats() sql.DBStats

	// Close closes the database connection.
	Close()
}

// SqlTx is a precis of *pgx.Tx except that the commit and rollback methods are unexported.
type SqlTx interface {
	Execer
	commit() error
	rollback() error
}

// SqlStmt is a precis of *sql.Stmt
type SqlStmt interface {
	// ExecContext executes a query without returning any rows.
	// The arguments are for any placeholder parameters in the query.
	ExecContext(ctx context.Context, arguments ...interface{}) (sql.Result, error)

	// QueryContext executes a query that returns rows, typically a SELECT.
	// The arguments are for any placeholder parameters in the query.
	QueryContext(ctx context.Context, arguments ...interface{}) (*pgx.Rows, error)

	// QueryRowContext executes a query that is expected to return at most one row.
	// QueryRowContext always returns a non-nil value. Errors are deferred until
	// Row's Scan method is called.
	// If the query selects no rows, the *Row's Scan will return ErrNoRows.
	// Otherwise, the *Row's Scan scans the first selected row and discards
	// the rest.
	QueryRowContext(ctx context.Context, arguments ...interface{}) *pgx.Row

	Close() error
}

// Type conformance assertions
//var _ SqlStmt = new(pgx.PreparedStatement)
