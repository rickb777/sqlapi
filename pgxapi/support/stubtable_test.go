package support

import (
	"context"
	"github.com/mercury-holidays/sqlapi/dialect"
	"github.com/mercury-holidays/sqlapi/pgxapi"
)

type StubTable struct {
	name     pgxapi.TableName
	dialect  dialect.Dialect
	database *StubDatabase
}

// Type conformance checks
var _ pgxapi.Table = &StubTable{}

func (tbl StubTable) Database() pgxapi.Database {
	return tbl.database
}

func (tbl StubTable) Logger() pgxapi.Logger {
	return nil
}

func (tbl StubTable) Ctx() context.Context {
	return context.Background()
}

func (tbl StubTable) Dialect() dialect.Dialect {
	return tbl.dialect
}

func (tbl StubTable) Name() pgxapi.TableName {
	return tbl.name
}

func (tbl StubTable) DB() pgxapi.SqlDB {
	return nil
}

func (tbl StubTable) Execer() pgxapi.Execer {
	return tbl.database.execer
}

func (tbl StubTable) Tx() pgxapi.SqlTx {
	return nil
}

func (tbl StubTable) IsTx() bool {
	return false
}

func (tbl StubTable) Query(query string, args ...interface{}) (pgxapi.SqlRows, error) {
	return nil, nil
}
