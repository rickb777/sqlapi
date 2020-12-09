package test

import (
	"context"

	"github.com/jackc/pgx/v4"
	"github.com/rickb777/sqlapi"
	"github.com/rickb777/sqlapi/driver"
)

// StubExecer provides a non-functioning Execer for testing purposes.
type StubExecer struct {
	N    int64
	Row  sqlapi.SqlRow
	Rows sqlapi.SqlRows
	Err  error
	Lgr  sqlapi.Logger
	Di   driver.Dialect
	User interface{}
}

var _ sqlapi.Execer = &StubExecer{}

// n.b. logging is absent here (it happens in the support functions)

func (e StubExecer) Query(ctx context.Context, query string, args ...interface{}) (sqlapi.SqlRows, error) {
	//e.Lgr.Log("%s %v", query, args)
	return e.Rows, e.Err
}

func (e StubExecer) QueryRow(ctx context.Context, query string, args ...interface{}) sqlapi.SqlRow {
	//e.Lgr.Log("%s %v", query, args)
	return e.Row
}

func (e StubExecer) Exec(ctx context.Context, query string, args ...interface{}) (int64, error) {
	//e.Lgr.Log("%s %v", query, args)
	return e.N, e.Err
}

func (e StubExecer) Insert(ctx context.Context, pk, query string, args ...interface{}) (int64, error) {
	//e.Lgr.Log("%s [%s] %v", query, pk, args)
	return e.N, e.Err
}

func (StubExecer) IsTx() bool {
	return false
}

func (e StubExecer) Logger() sqlapi.Logger {
	return e.Lgr
}

func (e StubExecer) Dialect() driver.Dialect {
	return e.Di
}

//-------------------------------------------------------------------------------------------------

func (e StubExecer) Transact(_ context.Context, txOptions *pgx.TxOptions, fn func(sqlapi.SqlTx) error) error {
	return fn(e)
}

func (e StubExecer) Ping(_ context.Context) error {
	return e.Err
}

func (e StubExecer) Stats() sqlapi.DBStats {
	return sqlapi.DBStats{}
}

func (e StubExecer) SingleConn(_ context.Context, fn func(ex sqlapi.Execer) error) error {
	return fn(e)
}

func (e StubExecer) Close() error {
	return e.Err
}

func (e StubExecer) With(userItem interface{}) sqlapi.SqlDB {
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
