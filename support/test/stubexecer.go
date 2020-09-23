package test

import (
	"context"
	"database/sql"

	"github.com/rickb777/sqlapi"
	"github.com/rickb777/sqlapi/dialect"
)

// StubExecer provides a non-functioning Execer for testing purposes.
type StubExecer struct {
	N    int64
	Row  sqlapi.SqlRow
	Rows sqlapi.SqlRows
	Err  error
	Lgr  sqlapi.Logger
	Di   dialect.Dialect
}

var _ sqlapi.Execer = &StubExecer{}

// n.b. logging is absent here (it happens in the support functions)

func (e StubExecer) QueryContext(ctx context.Context, query string, args ...interface{}) (sqlapi.SqlRows, error) {
	//e.Lgr.Log("%s %v", query, args)
	return e.Rows, e.Err
}

func (e StubExecer) QueryRowContext(ctx context.Context, query string, args ...interface{}) sqlapi.SqlRow {
	//e.Lgr.Log("%s %v", query, args)
	return e.Row
}

func (e StubExecer) ExecContext(ctx context.Context, query string, args ...interface{}) (int64, error) {
	//e.Lgr.Log("%s %v", query, args)
	return e.N, e.Err
}

func (e StubExecer) InsertContext(ctx context.Context, pk, query string, args ...interface{}) (int64, error) {
	//e.Lgr.Log("%s [%s] %v", query, pk, args)
	return e.N, e.Err
}

func (e StubExecer) PrepareContext(ctx context.Context, name, query string) (sqlapi.SqlStmt, error) {
	return nil, e.Err
}

func (StubExecer) IsTx() bool {
	return false
}

func (e StubExecer) Logger() sqlapi.Logger {
	return e.Lgr
}

func (e StubExecer) Dialect() dialect.Dialect {
	return e.Di
}

//-------------------------------------------------------------------------------------------------

func (e StubExecer) Transact(_ context.Context, txOptions *sql.TxOptions, fn func(sqlapi.SqlTx) error) error {
	return fn(e)
}

func (e StubExecer) PingContext(_ context.Context) error {
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

//-------------------------------------------------------------------------------------------------

func (e StubExecer) Commit() error {
	return e.Err
}

func (e StubExecer) Rollback() error {
	return e.Err
}
