package support

import (
	"context"

	"github.com/rickb777/sqlapi"
	"github.com/rickb777/sqlapi/dialect"
)

type StubExecer struct {
	StubResult int64
	Rows       sqlapi.SqlRows
	Lgr        sqlapi.Logger
}

func (e StubExecer) QueryContext(ctx context.Context, query string, args ...interface{}) (sqlapi.SqlRows, error) {
	return e.Rows, nil
}

func (e StubExecer) ExecContext(ctx context.Context, query string, args ...interface{}) (int64, error) {
	return e.StubResult, nil
}

func (e StubExecer) InsertContext(ctx context.Context, pk, query string, args ...interface{}) (int64, error) {
	return e.StubResult, nil
}

func (StubExecer) PrepareContext(ctx context.Context, name, query string) (sqlapi.SqlStmt, error) {
	return nil, nil
}

func (StubExecer) QueryRowContext(ctx context.Context, query string, args ...interface{}) sqlapi.SqlRow {
	return nil
}

func (StubExecer) IsTx() bool {
	panic("implement me")
}

func (e StubExecer) Logger() sqlapi.Logger {
	return e.Lgr
}

func (e StubExecer) Dialect() dialect.Dialect {
	panic("implement me")
}
