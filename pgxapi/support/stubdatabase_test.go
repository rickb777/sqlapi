package support

import (
	"context"
	"database/sql"
	"fmt"
	"github.com/jackc/pgx"
	"github.com/rickb777/sqlapi/dialect"
	"github.com/rickb777/sqlapi/pgxapi"
	"github.com/rickb777/sqlapi/util"
	"regexp"
)

type StubDatabase struct {
	execer        stubExecer
	loggedQueries []string
}

// Type conformance checks
var _ pgxapi.Database = &StubDatabase{}

func (*StubDatabase) DB() pgxapi.Execer {
	panic("implement DB")
}

func (*StubDatabase) Dialect() dialect.Dialect {
	panic("implement Dialect")
}

func (*StubDatabase) Logger() pgxapi.Logger {
	panic("implement Logger")
}

func (*StubDatabase) Wrapper() interface{} {
	panic("implement Wrapper")
}

func (*StubDatabase) PingContext(ctx context.Context) error {
	panic("implement PingContext")
}

func (*StubDatabase) Ping() error {
	panic("implement Ping")
}

func (*StubDatabase) Stats() sql.DBStats {
	panic("implement Stats")
}

func (*StubDatabase) ListTables(re *regexp.Regexp) (util.StringList, error) {
	panic("implement ListTables")
}

//-------------------------------------------------------------------------------------------------

type stubExecer struct {
	stubResult int64
}

// Type conformance checks
var _ pgxapi.Execer = &stubExecer{}

func (e stubExecer) QueryContext(ctx context.Context, query string, args ...interface{}) (pgxapi.SqlRows, error) {
	fmt.Printf("QueryContext: "+query+" %v", args...)
	return nil, nil
}

func (e stubExecer) QueryExRaw(ctx context.Context, sql string, options *pgx.QueryExOptions, args ...interface{}) (pgxapi.SqlRows, error) {
	panic("implement me")
}

func (e stubExecer) QueryRowContext(ctx context.Context, query string, args ...interface{}) pgxapi.SqlRow {
	fmt.Printf("QueryRowContext: "+query+" %v", args...)
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
	fmt.Printf("ExecContext: "+query+" %v", args...)
	return e.stubResult, nil
}

func (e stubExecer) PrepareContext(ctx context.Context, name, query string) (*pgx.PreparedStatement, error) {
	fmt.Printf("PrepareContext: " + query)
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
