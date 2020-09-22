package support

import (
	"context"
	"github.com/jackc/pgx"
	"github.com/rickb777/sqlapi/dialect"
	"github.com/rickb777/sqlapi/pgxapi"
	"strconv"
)

type StubExecer struct {
	StubResult int64
	Rows       pgxapi.SqlRows
	Lgr        pgxapi.Logger
}

func (e StubExecer) QueryContext(ctx context.Context, query string, args ...interface{}) (pgxapi.SqlRows, error) {
	e.Lgr.Log(pgx.LogLevelInfo, query, argMap(args...))
	return e.Rows, nil
}

func (e StubExecer) QueryExRaw(ctx context.Context, query string, options *pgx.QueryExOptions, args ...interface{}) (pgxapi.SqlRows, error) {
	e.Lgr.Log(pgx.LogLevelInfo, query, argMap(args...))
	return nil, nil
}

func (e StubExecer) QueryRowContext(ctx context.Context, query string, arguments ...interface{}) pgxapi.SqlRow {
	panic("implement me")
}

func (e StubExecer) QueryRowExRaw(ctx context.Context, query string, options *pgx.QueryExOptions, arguments ...interface{}) pgxapi.SqlRow {
	panic("implement me")
}

func (e StubExecer) BeginBatch() *pgx.Batch {
	panic("implement me")
}

func (e StubExecer) PrepareContext(ctx context.Context, name, sql string) (*pgx.PreparedStatement, error) {
	return nil, nil
}

func (e StubExecer) ExecContext(ctx context.Context, query string, args ...interface{}) (int64, error) {
	e.Lgr.Log(pgx.LogLevelInfo, query, argMap(args...))
	return e.StubResult, nil
}

func (e StubExecer) InsertContext(ctx context.Context, pk, query string, args ...interface{}) (int64, error) {
	e.Lgr.Log(pgx.LogLevelInfo, query, argMap(args...))
	return e.StubResult, nil
}

func (StubExecer) IsTx() bool {
	panic("implement me")
}

func (e StubExecer) Logger() pgxapi.Logger {
	return e.Lgr
}

func (e StubExecer) Dialect() dialect.Dialect {
	return dialect.Postgres
}

func argMap(args ...interface{}) map[string]interface{} {
	m := make(map[string]interface{})
	for i, a := range args {
		m[strconv.Itoa(i)] = a
	}
	return m
}
