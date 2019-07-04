package support

import (
	"github.com/jackc/pgx"
	"github.com/rickb777/sqlapi/pgxapi"
)

type StubRow []interface{}

type StubRows struct {
	I          int
	Rows       []StubRow
	Error      error
	Fields     []pgx.FieldDescription
	ValueSlice []interface{}
}

var _ pgxapi.SqlRows = new(StubRows)

func (r *StubRows) Next() bool {
	return r.I < len(r.Rows)
}

func (r *StubRows) Scan(dest ...interface{}) error {
	for i, v := range r.Rows[r.I] {
		dest[i] = v
	}
	r.I++
	return nil
}

func (r *StubRows) FieldDescriptions() []pgx.FieldDescription {
	return r.Fields
}

func (r *StubRows) Values() ([]interface{}, error) {
	return r.ValueSlice, nil
}

func (r *StubRows) Close() {}

func (r *StubRows) Err() error {
	return r.Error
}
