package support

import (
	"context"
	"fmt"
	"github.com/mercury-holidays/sqlapi"
	"github.com/mercury-holidays/sqlapi/dialect"
	"github.com/mercury-holidays/sqlapi/util"
	"io"
	"regexp"
)

type StubDatabase struct {
	execer stubExecer
	stdLog *stubLogger
}

func (*StubDatabase) DB() sqlapi.SqlDB {
	panic("implement me")
}

func (*StubDatabase) Dialect() dialect.Dialect {
	panic("implement me")
}

func (d *StubDatabase) Logger() sqlapi.Logger {
	return sqlapi.NewLogger(d.stdLog)
}

func (*StubDatabase) Wrapper() interface{} {
	panic("implement me")
}

func (*StubDatabase) ListTables(re *regexp.Regexp) (util.StringList, error) {
	panic("implement me")
}

//-------------------------------------------------------------------------------------------------

type stubLogger struct {
	logged []string
}

func (r *stubLogger) Printf(format string, v ...interface{}) {
	r.logged = append(r.logged, fmt.Sprintf(format, v...))
}

func (r *stubLogger) SetOutput(w io.Writer) {}

//-------------------------------------------------------------------------------------------------

type stubExecer struct {
	stubResult int64
	rows       sqlapi.SqlRows
}

func (e stubExecer) ExecContext(ctx context.Context, query string, args ...interface{}) (int64, error) {
	return e.stubResult, nil
}

func (e stubExecer) InsertContext(ctx context.Context, query string, args ...interface{}) (int64, error) {
	return e.stubResult, nil
}

func (stubExecer) PrepareContext(ctx context.Context, name, query string) (sqlapi.SqlStmt, error) {
	return nil, nil
}

func (se stubExecer) QueryContext(ctx context.Context, query string, args ...interface{}) (sqlapi.SqlRows, error) {
	return se.rows, nil
}

func (stubExecer) QueryRowContext(ctx context.Context, query string, args ...interface{}) sqlapi.SqlRow {
	return nil
}

func (stubExecer) IsTx() bool {
	panic("implement me")
}
