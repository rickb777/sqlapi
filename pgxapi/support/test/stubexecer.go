package test

import (
	"context"
	"fmt"
	"github.com/jackc/pgx/v4"
	"github.com/rickb777/sqlapi/dialect"
	"github.com/rickb777/sqlapi/pgxapi"
	"github.com/rickb777/where/quote"
)

// StubExecer provides a non-functioning Execer for testing purposes.
type StubExecer struct {
	N    int64
	Row  pgxapi.SqlRow
	Rows pgxapi.SqlRows
	Err  error
	Lgr  pgxapi.Logger
	Q    quote.Quoter
	User interface{}
}

var _ pgxapi.Execer = &StubExecer{}

// n.b. logging is included here because this emulates the behaviour of pgx

func (e StubExecer) Query(ctx context.Context, query string, args ...interface{}) (pgxapi.SqlRows, error) {
	e.Lgr.Log(ctx, pgx.LogLevelInfo, query, argMap(args...))
	return e.Rows, e.Err
}

func (e StubExecer) QueryRow(ctx context.Context, query string, args ...interface{}) pgxapi.SqlRow {
	e.Lgr.Log(ctx, pgx.LogLevelInfo, query, argMap(args...))
	return e.Row
}

func (e StubExecer) Exec(ctx context.Context, query string, args ...interface{}) (int64, error) {
	e.Lgr.Log(ctx, pgx.LogLevelInfo, query, argMap(args...))
	return e.N, e.Err
}

func (e StubExecer) Insert(ctx context.Context, pk, query string, args ...interface{}) (int64, error) {
	e.Lgr.Log(ctx, pgx.LogLevelInfo, query, argMap(args...))
	return e.N, e.Err
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

func (e StubExecer) Ping(_ context.Context) error {
	return e.Err
}

func (e StubExecer) Stats() pgxapi.DBStats {
	return pgxapi.DBStats{}
}

func (e StubExecer) SingleConn(_ context.Context, fn func(ex pgxapi.Execer) error) error {
	return fn(e)
}

func (e StubExecer) Close() error {
	return nil
}

func (e StubExecer) With(userItem interface{}) pgxapi.SqlDB {
	e.User = userItem
	return e
}

func (e StubExecer) UserItem() interface{} {
	return e.User
}

//-------------------------------------------------------------------------------------------------

func (e StubExecer) Commit(_ context.Context) error {
	return e.Err
}

func (e StubExecer) Rollback(_ context.Context) error {
	return e.Err
}

func argMap(args ...interface{}) map[string]interface{} {
	m := make(map[string]interface{})
	for i, a := range args {
		m[fmt.Sprintf("$%d", i+1)] = a
	}
	return m
}
