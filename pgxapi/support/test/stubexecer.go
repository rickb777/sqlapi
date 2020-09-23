package test

import (
	"context"
	"github.com/jackc/pgx"
	"github.com/rickb777/sqlapi/dialect"
	"github.com/rickb777/sqlapi/pgxapi"
	"github.com/rickb777/where/quote"
	"strconv"
)

// StubExecer provides a non-functioning Execer for testing purposes.
type StubExecer struct {
	N    int64
	Row  pgxapi.SqlRow
	Rows pgxapi.SqlRows
	Err  error
	Lgr  pgxapi.Logger
	Q    quote.Quoter
}

var _ pgxapi.Execer = &StubExecer{}

// n.b. logging is included here because this emulates the behaviour of pgx

func (e StubExecer) QueryContext(ctx context.Context, query string, args ...interface{}) (pgxapi.SqlRows, error) {
	e.Lgr.Log(pgx.LogLevelInfo, query, argMap(args...))
	return e.Rows, e.Err
}

func (e StubExecer) QueryExRaw(ctx context.Context, query string, options *pgx.QueryExOptions, args ...interface{}) (pgxapi.SqlRows, error) {
	e.Lgr.Log(pgx.LogLevelInfo, query, argMap(args...))
	return nil, e.Err
}

func (e StubExecer) QueryRowContext(ctx context.Context, query string, args ...interface{}) pgxapi.SqlRow {
	e.Lgr.Log(pgx.LogLevelInfo, query, argMap(args...))
	return e.Row
}

func (e StubExecer) QueryRowExRaw(ctx context.Context, query string, options *pgx.QueryExOptions, args ...interface{}) pgxapi.SqlRow {
	e.Lgr.Log(pgx.LogLevelInfo, query, argMap(args...))
	return e.Row
}

func (e StubExecer) ExecContext(ctx context.Context, query string, args ...interface{}) (int64, error) {
	e.Lgr.Log(pgx.LogLevelInfo, query, argMap(args...))
	return e.N, e.Err
}

func (e StubExecer) InsertContext(ctx context.Context, pk, query string, args ...interface{}) (int64, error) {
	e.Lgr.Log(pgx.LogLevelInfo, query, argMap(args...))
	return e.N, e.Err
}

// PrepareContext is not implemented
func (e StubExecer) PrepareContext(ctx context.Context, name, query string) (*pgx.PreparedStatement, error) {
	return nil, e.Err
}

func (StubExecer) IsTx() bool {
	return false
}

func (e StubExecer) Logger() pgxapi.Logger {
	return e.Lgr
}

// BeginBatch is not implemented
func (e StubExecer) BeginBatch() *pgx.Batch {
	panic("implement me")
}

func (e StubExecer) Dialect() dialect.Dialect {
	if e.Q == nil {
		return postgresNoQuotes
	}
	return dialect.Postgres.WithQuoter(e.Q)
}

var postgresNoQuotes = dialect.Postgres.WithQuoter(quote.NoQuoter)

//-------------------------------------------------------------------------------------------------

func (e StubExecer) Transact(_ context.Context, txOptions *pgx.TxOptions, fn func(pgxapi.SqlTx) error) error {
	return fn(e)
}

func (e StubExecer) PingContext(_ context.Context) error {
	return e.Err
}

func (e StubExecer) Stats() pgxapi.DBStats {
	return pgxapi.DBStats{}
}

func (e StubExecer) SingleConn(_ context.Context, fn func(ex pgxapi.Execer) error) error {
	return fn(e)
}

func (e StubExecer) Close() {}

//-------------------------------------------------------------------------------------------------

func (e StubExecer) Commit() error {
	return e.Err
}

func (e StubExecer) Rollback() error {
	return e.Err
}

func argMap(args ...interface{}) map[string]interface{} {
	m := make(map[string]interface{})
	for i, a := range args {
		m[strconv.Itoa(i)] = a
	}
	return m
}
