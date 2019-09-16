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
	pgxLog pgxapi.Logger
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
	return d.pgxLog
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
}

// Type conformance checks
var _ pgxapi.Execer = &stubExecer{}

func (e stubExecer) QueryContext(ctx context.Context, query string, args ...interface{}) (pgxapi.SqlRows, error) {
	return nil, nil
}

func (e stubExecer) QueryExRaw(ctx context.Context, sql string, options *pgx.QueryExOptions, args ...interface{}) (pgxapi.SqlRows, error) {
	panic("implement me")
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
	return e.stubResult, nil
}

func (e stubExecer) PrepareContext(ctx context.Context, name, query string) (*pgx.PreparedStatement, error) {
	return nil, nil
}

func (e stubExecer) IsTx() bool {
	panic("implement me")
}

func (e stubExecer) Logger() pgxapi.Logger {
	panic("implement me")
}

//-------------------------------------------------------------------------------------------------

type stubResult struct {
	li, ra int64
	err    error
}

func (r stubResult) LastInsertId() (int64, error) {
	return r.li, r.err
}

func (r stubResult) RowsAffected() (int64, error) {
	return r.ra, r.err
}
