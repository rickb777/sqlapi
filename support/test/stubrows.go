package test

import (
	"database/sql"
	"reflect"

	"github.com/rickb777/sqlapi"
)

// StubRow provides a non-functioning pgxapi.SqlRow for testing purposes.
type StubRow []interface{}

// StubRow provides a non-functioning pgxapi.SqlRows for testing purposes.
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

func setValue(dest, v interface{}) {
	vv := reflect.ValueOf(v)
	reflect.ValueOf(dest).Elem().Set(vv)
}

func (r StubRow) Scan(dest ...interface{}) error {
	for i, v := range r {
		setValue(dest[i], v)
	}
	return nil
}

func (r *StubRows) Scan(dest ...interface{}) error {
	for i, v := range r.Rows[r.I] {
		setValue(dest[i], v)
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
