package pgxapi

import (
	"context"
	"database/sql"
	"github.com/jackc/pgx"
	"time"
)

type Logger = pgx.Logger

type Getter interface {
	QueryEx(ctx context.Context, sql string, options *pgx.QueryExOptions, args ...interface{}) (SqlRows, error)
	QueryExRaw(ctx context.Context, sql string, options *pgx.QueryExOptions, args ...interface{}) (SqlRows, error)
	QueryRowEx(ctx context.Context, query string, options *pgx.QueryExOptions, args ...interface{}) SqlRow
	QueryRowExRaw(ctx context.Context, query string, options *pgx.QueryExOptions, args ...interface{}) SqlRow
}

type Batcher interface {
	// BeginBatch exposes the pgx batch operations.
	BeginBatch() *pgx.Batch
}

type Lgr interface {
	pgx.Logger
	LogT(level pgx.LogLevel, msg string, startTime *time.Time, data ...interface{})
	TraceLogging(on bool)
}

type Execer interface {
	Getter
	Batcher
	Lgr

	InsertEx(ctx context.Context, query string, options *pgx.QueryExOptions, args ...interface{}) (int64, error)
	ExecEx(ctx context.Context, sql string, options *pgx.QueryExOptions, arguments ...interface{}) (pgx.CommandTag, error)
	PrepareEx(ctx context.Context, name, sql string, opts *pgx.PrepareExOptions) (*pgx.PreparedStatement, error)
	IsTx() bool
}

func IsTx(ex Execer) bool {
	return ex.IsTx()
}

//-------------------------------------------------------------------------------------------------

// SqlDB is able to make queries and begin transactions.
type SqlDB interface {
	Execer
	BeginTx(ctx context.Context, opts *pgx.TxOptions) (SqlTx, error)
	Transact(ctx context.Context, txOptions *pgx.TxOptions, fn func(Execer) error) error
	PingContext(ctx context.Context) error
	Stats() sql.DBStats
	Close()
}

// SqlTx is a precis of *pgx.Tx
type SqlTx interface {
	Execer
	Commit() error
	Rollback() error
}

// SqlStmt is a precis of *sql.Stmt
type SqlStmt interface {
	// ExecContext executes a query without returning any rows.
	// The args are for any placeholder parameters in the query.
	ExecContext(ctx context.Context, args ...interface{}) (sql.Result, error)

	// QueryContext executes a query that returns rows, typically a SELECT.
	// The args are for any placeholder parameters in the query.
	QueryContext(ctx context.Context, args ...interface{}) (*pgx.Rows, error)

	// QueryRowContext executes a query that is expected to return at most one row.
	// QueryRowContext always returns a non-nil value. Errors are deferred until
	// Row's Scan method is called.
	// If the query selects no rows, the *Row's Scan will return ErrNoRows.
	// Otherwise, the *Row's Scan scans the first selected row and discards
	// the rest.
	QueryRowContext(ctx context.Context, args ...interface{}) *pgx.Row

	Close() error
}
