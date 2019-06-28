package pgxapi

import (
	"context"

	"github.com/jackc/pgx"
	"github.com/pkg/errors"
	"github.com/rickb777/sqlapi/dialect"
)

type Getter interface {
	QueryEx(ctx context.Context, sql string, options *pgx.QueryExOptions, args ...interface{}) (*pgx.Rows, error)
	QueryExRaw(ctx context.Context, sql string, options *pgx.QueryExOptions, args ...interface{}) (*pgx.Rows, error)
	QueryRowEx(ctx context.Context, query string, options *pgx.QueryExOptions, args ...interface{}) *pgx.Row
	QueryRowExRaw(ctx context.Context, query string, options *pgx.QueryExOptions, args ...interface{}) *pgx.Row
}

type Batcher interface {
	// BeginBatch exposes the pgx batch operations.
	BeginBatch() *pgx.Batch
}

type Execer interface {
	Getter
	Batcher
	InsertEx(ctx context.Context, query string, options *pgx.QueryExOptions, args ...interface{}) (int64, error)
	ExecEx(ctx context.Context, sql string, options *pgx.QueryExOptions, arguments ...interface{}) (pgx.CommandTag, error)
	PrepareEx(ctx context.Context, name, sql string, opts *pgx.PrepareExOptions) (*pgx.PreparedStatement, error)
}

//-------------------------------------------------------------------------------------------------

func WrapTX(tx *pgx.Tx) Execer {
	return &txShim{tx: tx}
}

type txShim struct {
	tx *pgx.Tx
}

// QueryExReplacePlaceholders executes any query returning data rows,
// having first replaced all '?' placeholders with numbered ones.
func (ts *txShim) QueryEx(ctx context.Context, query string, options *pgx.QueryExOptions, args ...interface{}) (*pgx.Rows, error) {
	qr := dialect.Postgres.ReplacePlaceholders(query, nil)
	return ts.QueryExRaw(ctx, qr, options, args...)
}

// QueryEx executes any query returning data rows.
func (ts *txShim) QueryExRaw(ctx context.Context, query string, options *pgx.QueryExOptions, args ...interface{}) (*pgx.Rows, error) {
	rows, err := ts.tx.QueryEx(ctx, query, options, args...)
	if err != nil {
		return nil, errors.Wrapf(err, "%s %v", query, args)
	}
	return rows, nil
}

// QueryRowEx executes any query returning one data row,
// having first replaced all '?' placeholders with numbered ones.
func (ts *txShim) QueryRowEx(ctx context.Context, query string, options *pgx.QueryExOptions, args ...interface{}) *pgx.Row {
	qr := dialect.Postgres.ReplacePlaceholders(query, nil)
	return ts.QueryRowExRaw(ctx, qr, options, args...)
}

// QueryRowEx executes any query returning one data row.
func (ts *txShim) QueryRowExRaw(ctx context.Context, query string, options *pgx.QueryExOptions, args ...interface{}) *pgx.Row {
	return ts.tx.QueryRowEx(ctx, query, options, args...)
}

// InsertEx executes an insert query returning an ID.
func (ts *txShim) InsertEx(ctx context.Context, query string, options *pgx.QueryExOptions, args ...interface{}) (int64, error) {
	row := ts.tx.QueryRowEx(ctx, query, nil, args...)
	var id int64
	err := row.Scan(&id)
	if err != nil && err != pgx.ErrNoRows {
		return 0, errors.Wrapf(err, "%s %v", query, args)
	}
	return id, errors.Wrapf(err, "%s %v", query, args)
}

// ExecEx executes any query returning nothing.
func (ts *txShim) ExecEx(ctx context.Context, query string, options *pgx.QueryExOptions, args ...interface{}) (pgx.CommandTag, error) {
	tag, err := ts.tx.ExecEx(ctx, query, options, args...)
	return tag, errors.Wrapf(err, "%s %v", query, args)
}

// PrepareEx prepares a statement for repeated use.
func (ts *txShim) PrepareEx(ctx context.Context, name, sql string, opts *pgx.PrepareExOptions) (*pgx.PreparedStatement, error) {
	ps, err := ts.tx.PrepareEx(ctx, name, sql, opts)
	return ps, errors.Wrapf(err, "%s %s", name, sql)
}

// BeginBatch begins a batch operation.
func (ts *txShim) BeginBatch() *pgx.Batch {
	return ts.tx.BeginBatch()
}
