package support

import (
	"context"
	"database/sql"
	"github.com/rickb777/sqlapi"
	"github.com/rickb777/sqlapi/dialect"
	"log"
)

type StubTable struct {
	name    sqlapi.TableName
	dialect dialect.Dialect
}

// Type conformance checks
var _ sqlapi.Table = &StubTable{}

func (tbl StubTable) Database() sqlapi.Database {
	return nil
}

func (tbl StubTable) Logger() *log.Logger {
	return nil
}

func (tbl StubTable) Ctx() context.Context {
	return context.Background()
}

func (tbl StubTable) Dialect() dialect.Dialect {
	return tbl.dialect
}

func (tbl StubTable) Name() sqlapi.TableName {
	return tbl.name
}

func (tbl StubTable) DB() *sql.DB {
	return nil
}

func (tbl StubTable) Execer() sqlapi.Execer {
	return nil
}

func (tbl StubTable) Tx() *sql.Tx {
	return nil
}

func (tbl StubTable) IsTx() bool {
	return false
}

func (tbl StubTable) Query(query string, args ...interface{}) (sqlapi.SqlRows, error) {
	return nil, nil
}

//-------------------------------------------------------------------------------------------------
