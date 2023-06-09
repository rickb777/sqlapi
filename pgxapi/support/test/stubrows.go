package test

import (
	"reflect"

	"github.com/jackc/pgx/v5/pgconn"
	"github.com/rickb777/sqlapi/pgxapi"
)

// StubRow provides a non-functioning pgxapi.SqlRow for testing purposes.
type StubRow []interface{}

// StubRow provides a non-functioning pgxapi.SqlRows for testing purposes.
type StubRows struct {
	I          int
	Rows       []StubRow
	Error      error
	Fields     []pgconn.FieldDescription
	ValueSlice []interface{}
}

var _ pgxapi.SqlRows = new(StubRows)

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

func (r *StubRows) FieldDescriptions() []pgconn.FieldDescription {
	return r.Fields
}

func (r *StubRows) Values() ([]interface{}, error) {
	return r.ValueSlice, nil
}

func (r *StubRows) Close() {}

func (r *StubRows) Err() error {
	return r.Error
}
