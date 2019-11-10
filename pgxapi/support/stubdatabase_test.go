package support

import (
	"context"
	"fmt"
	"github.com/jackc/pgx"
	"github.com/rickb777/collection"
	"github.com/rickb777/sqlapi/dialect"
	"github.com/rickb777/sqlapi/pgxapi"
	"regexp"
)

type StubDatabase struct {
	execer stubExecer
}

// Type conformance checks
var _ pgxapi.Database = &StubDatabase{}

func (*StubDatabase) DB() pgxapi.SqlDB {
	panic("implement DB")
}

func (*StubDatabase) Dialect() dialect.Dialect {
	panic("implement Dialect")
}

func (d *StubDatabase) Logger() pgxapi.Logger {
	return d.execer.Logger()
}

func (*StubDatabase) Wrapper() interface{} {
	panic("implement Wrapper")
}

func (*StubDatabase) ListTables(re *regexp.Regexp) (collection.StringList, error) {
	panic("implement ListTables")
}

//-------------------------------------------------------------------------------------------------

type stubLogger struct {
	logged []string
}

func (r *stubLogger) Log(level pgx.LogLevel, msg string, data map[string]interface{}) {
	r.logged = append(r.logged, fmt.Sprintf("%s %v", msg, data))
}

//-------------------------------------------------------------------------------------------------

type stubExecer struct {
	stubResult int64
	rows       pgxapi.SqlRows
	pgxLog     *stubLogger
}

// Type conformance checks
var _ pgxapi.Execer = &stubExecer{}

func (e stubExecer) QueryContext(ctx context.Context, query string, args ...interface{}) (pgxapi.SqlRows, error) {
	e.pgxLog.Log(pgx.LogLevelInfo, fmt.Sprintf("%s %v", query, args), nil)
	return e.rows, nil
}

func (e stubExecer) QueryExRaw(ctx context.Context, sql string, options *pgx.QueryExOptions, args ...interface{}) (pgxapi.SqlRows, error) {
	e.pgxLog.Log(pgx.LogLevelInfo, fmt.Sprintf("%s %v", sql, args), nil)
	return e.rows, nil
}

func (e stubExecer) QueryRowContext(ctx context.Context, query string, args ...interface{}) pgxapi.SqlRow {
	return nil
}

func (e stubExecer) QueryRowExRaw(ctx context.Context, query string, options *pgx.QueryExOptions, args ...interface{}) pgxapi.SqlRow {
	panic("implement me")
}

func (e stubExecer) BeginBatch() *pgx.Batch {
	panic("implement me")
}

func (e stubExecer) InsertContext(ctx context.Context, query string, args ...interface{}) (int64, error) {
	panic("implement me")
}

func (e stubExecer) ExecContext(ctx context.Context, query string, args ...interface{}) (int64, error) {
	e.pgxLog.Log(pgx.LogLevelInfo, fmt.Sprintf("%s %v", query, args), nil)
	return e.stubResult, nil
}

func (e stubExecer) PrepareContext(ctx context.Context, name, query string) (*pgx.PreparedStatement, error) {
	return nil, nil
}

func (e stubExecer) IsTx() bool {
	panic("implement me")
}

func (e stubExecer) Logger() pgxapi.Logger {
	return pgxapi.NewLogger(e.pgxLog)
}
