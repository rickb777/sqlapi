package pgxapi

import (
	"github.com/jackc/pgproto3/v2"
	"github.com/rickb777/collection"
)

// SqlRow is a precis of *sql.Row.
type SqlRow interface {
	Scan(dest ...interface{}) error
}

// SqlRows is a precis of pgx.Rows.
type SqlRows interface {
	SqlRow

	// Next prepares the next row for reading. It returns true if there is another
	// row and false if no more rows are available. It automatically closes rows
	// when all rows are read.
	Next() bool

	FieldDescriptions() []pgproto3.FieldDescription

	// Values returns an array of the row values
	Values() ([]interface{}, error)

	// Close closes the rows, making the connection ready for use again. It is safe
	// to call Close after rows is already closed.
	Close()

	Err() error
}

//-------------------------------------------------------------------------------------------------

// Rows provides a tool for scanning result *sql.Rows of arbitrary or varying length.
// The internal *sql.Rows field is exported and is usable as per normal via its Next
// and Scan methods, or the Next and ScanToMap methods can be used instead.
type Rows struct {
	fields []pgproto3.FieldDescription
	values []interface{}
	Rows   SqlRows
}

// RowData holds a single row result from the database.
type RowData struct {
	Fields []pgproto3.FieldDescription
	Data   collection.StringAnyMap
}

// WrapRows wraps a *sql.Rows result so that its data can be scanned into a series of
// maps, one for each row.
func WrapRows(rows SqlRows) (*Rows, error) {
	cols := rows.FieldDescriptions()

	vv, err := rows.Values()
	if err != nil {
		return nil, err
	}

	return &Rows{cols, vv, rows}, nil
}

// Next prepares the next result row for reading with the Scan method. It
// returns true on success, or false if there is no next result row or an error
// happened while preparing it. Err should be consulted to distinguish between
// the two cases.
//
// Every call to Scan, even the first one, must be preceded by a call to Next.
func (rams *Rows) Next() bool {
	return rams.Rows.Next()
}

// Scan copies the columns in the current row into the values pointed
// at by dest. The number of values in dest must be the same as the
// number of columns in the wrapped Rows.
func (rams *Rows) Scan(dest ...interface{}) error {
	return rams.Rows.Scan(dest...)
}

// ScanToMap copies all the column data of the current row into a map. The
// map is keyed by column name.
//
// The result describes a single row from the database, consisting of the
// column names, types and data.
func (rams *Rows) ScanToMap() (RowData, error) {
	values := make([]interface{}, len(rams.fields))
	valuePointers := make([]interface{}, len(rams.fields))
	for i := range values {
		valuePointers[i] = &values[i]
	}

	if err := rams.Rows.Scan(valuePointers...); err != nil {
		return RowData{}, err
	}

	m := make(collection.StringAnyMap)
	for i, field := range rams.fields {
		v := valuePointers[i].(*interface{})
		m[string(field.Name)] = *v
	}

	return RowData{rams.fields, m}, nil
}

// Close closes the Rows, preventing further enumeration. If Next is called
// and returns false and there are no further result sets,
// the Rows are closed automatically and it will suffice to check the
// result of Err. Close is idempotent and does not affect the result of Err.
func (rams *Rows) Close() {
	rams.Rows.Close()
}

// Err returns the error, if any, that was encountered during iteration.
// Err may be called after an explicit or implicit Close.
func (rams *Rows) Err() error {
	return rams.Rows.Err()
}

func (rams *Rows) FieldDescriptions() []pgproto3.FieldDescription {
	return rams.fields
}
