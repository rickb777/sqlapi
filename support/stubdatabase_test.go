package support

import (
	"context"
	"database/sql"
	"fmt"
	"github.com/rickb777/sqlapi"
	"github.com/rickb777/sqlapi/dialect"
	"github.com/rickb777/sqlapi/util"
	"log"
	"regexp"
)

type StubDatabase struct {
	execer        stubExecer
	loggedQueries []string
}

func (*StubDatabase) DB() sqlapi.Execer {
	panic("implement me")
}

func (*StubDatabase) BeginTx(ctx context.Context, opts *sql.TxOptions) (sqlapi.SqlTx, error) {
	panic("implement me")
}

func (*StubDatabase) Begin() (sqlapi.SqlTx, error) {
	panic("implement me")
}

func (*StubDatabase) Dialect() dialect.Dialect {
	panic("implement me")
}

func (*StubDatabase) Logger() *log.Logger {
	panic("implement me")
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

func (*StubDatabase) TraceLogging(on bool) {
	panic("implement me")
}

func (d *StubDatabase) LogQuery(query string, args ...interface{}) {
	d.loggedQueries = append(d.loggedQueries, query)
	d.loggedQueries = append(d.loggedQueries, fmt.Sprintf("%+v", args))
}

func (d *StubDatabase) LogIfError(err error) error {
	if err != nil {
		d.loggedQueries = append(d.loggedQueries, err.Error())
	}
	return err
}

func (*StubDatabase) LogError(err error) error {
	panic("implement me")
}

func (*StubDatabase) Exec(query string, args ...interface{}) (sql.Result, error) {
	panic("implement me")
}

func (d *StubDatabase) ExecContext(ctx context.Context, query string, args ...interface{}) (sql.Result, error) {
	return d.execer.ExecContext(ctx, query, args...)
}

func (d *StubDatabase) Prepare(query string) (sqlapi.SqlStmt, error) {
	return d.PrepareContext(context.Background(), query)
}

func (d *StubDatabase) PrepareContext(ctx context.Context, query string) (sqlapi.SqlStmt, error) {
	return d.execer.PrepareContext(ctx, query)
}

func (d *StubDatabase) Query(query string, args ...interface{}) (sqlapi.SqlRows, error) {
	return d.QueryContext(context.Background(), query, args...)
}

func (d *StubDatabase) QueryContext(ctx context.Context, query string, args ...interface{}) (sqlapi.SqlRows, error) {
	return d.execer.QueryContext(ctx, query, args...)
}

func (*StubDatabase) QueryRow(query string, args ...interface{}) sqlapi.SqlRow {
	panic("implement me")
}

func (*StubDatabase) QueryRowContext(ctx context.Context, query string, args ...interface{}) sqlapi.SqlRow {
	panic("implement me")
}

func (*StubDatabase) ListTables(re *regexp.Regexp) (util.StringList, error) {
	panic("implement me")
}

//-------------------------------------------------------------------------------------------------

type stubExecer struct {
	stubResult stubResult
}

func (e stubExecer) ExecContext(ctx context.Context, query string, args ...interface{}) (sql.Result, error) {
	fmt.Printf("ExecContext: "+query+" %v", args...)
	return e.stubResult, nil
}

func (stubExecer) PrepareContext(ctx context.Context, query string) (*sql.Stmt, error) {
	fmt.Printf("PrepareContext: " + query)
	return nil, nil
}

func (stubExecer) QueryContext(ctx context.Context, query string, args ...interface{}) (*sql.Rows, error) {
	fmt.Printf("QueryContext: "+query+" %v", args...)
	return nil, nil
}

func (stubExecer) QueryRowContext(ctx context.Context, query string, args ...interface{}) *sql.Row {
	fmt.Printf("QueryRowContext: "+query+" %v", args...)
	return nil
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
