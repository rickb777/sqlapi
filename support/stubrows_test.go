package support

import (
	"database/sql"
	"reflect"

	"github.com/rickb777/sqlapi"
)

type StubRow []interface{}

type StubRows struct {
	I     int
	Rows  []StubRow
	Error error
	Cols  []string
	Types []*sql.ColumnType
}

var _ sqlapi.SqlRows = new(StubRows)

func (r *StubRows) Next() bool {
	return r.I < len(r.Rows)
}

func (r *StubRows) Scan(dest ...interface{}) error {
	for i, v := range r.Rows[r.I] {
		vv := reflect.ValueOf(v)
		reflect.ValueOf(dest[i]).Elem().Set(vv)
	}
	r.I++
	return nil
}

func (r *StubRows) Columns() ([]string, error) {
	return r.Cols, nil
}

func (r *StubRows) ColumnTypes() ([]*sql.ColumnType, error) {
	return r.Types, nil
}

func (r *StubRows) Close() error {
	return nil
}

func (r *StubRows) Err() error {
	return r.Error
}
