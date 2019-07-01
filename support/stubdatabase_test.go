package support

import (
	"context"
	"database/sql"
	"fmt"
	"github.com/rickb777/sqlapi"
	"github.com/rickb777/sqlapi/dialect"
	"github.com/rickb777/sqlapi/util"
	"io"
	"regexp"
)

type StubDatabase struct {
	execer stubExecer
	stdLog *stubLogger
}

func (*StubDatabase) DB() sqlapi.Execer {
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

func (*StubDatabase) PingContext(ctx context.Context) error {
	panic("implement me")
}

func (*StubDatabase) Ping() error {
	panic("implement me")
}

func (*StubDatabase) Stats() sql.DBStats {
	panic("implement me")
}

func (*StubDatabase) ListTables(re *regexp.Regexp) (util.StringList, error) {
	panic("implement me")
}

//-------------------------------------------------------------------------------------------------

type stubExecer struct {
	stubResult int64
}

func (e stubExecer) ExecContext(ctx context.Context, query string, args ...interface{}) (int64, error) {
	fmt.Printf("ExecContext: "+query+" %v", args...)
	return e.stubResult, nil
}

func (e stubExecer) InsertContext(ctx context.Context, query string, args ...interface{}) (int64, error) {
	fmt.Printf("InsertContext: "+query+" %v", args...)
	return e.stubResult, nil
}

func (stubExecer) PrepareContext(ctx context.Context, name, query string) (sqlapi.SqlStmt, error) {
	fmt.Printf("PrepareContext: " + query)
	return nil, nil
}

func (stubExecer) QueryContext(ctx context.Context, query string, args ...interface{}) (sqlapi.SqlRows, error) {
	fmt.Printf("QueryContext: "+query+" %v", args...)
	return nil, nil
}

func (stubExecer) QueryRowContext(ctx context.Context, query string, args ...interface{}) sqlapi.SqlRow {
	fmt.Printf("QueryRowContext: "+query+" %v", args...)
	return nil
}

func (stubExecer) IsTx() bool {
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
