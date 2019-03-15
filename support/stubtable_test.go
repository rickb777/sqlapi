package support

import (
	"context"
	"github.com/rickb777/sqlapi"
	"github.com/rickb777/sqlapi/dialect"
	"log"
)

type StubTable struct {
	name     sqlapi.TableName
	dialect  dialect.Dialect
	database *StubDatabase
}

// Type conformance checks
var _ sqlapi.Table = &StubTable{}

func (tbl StubTable) Database() sqlapi.Database {
	return tbl.database
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

func (tbl StubTable) DB() sqlapi.SqlDB {
	return nil
}

func (tbl StubTable) Execer() sqlapi.Execer {
	return tbl.database.execer
}

func (tbl StubTable) Tx() sqlapi.SqlTx {
	return nil
}

func (tbl StubTable) IsTx() bool {
	return false
}

func (tbl StubTable) Query(query string, args ...interface{}) (sqlapi.SqlRows, error) {
	return nil, nil
}
