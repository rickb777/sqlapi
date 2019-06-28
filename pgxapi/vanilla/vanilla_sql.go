// THIS FILE WAS AUTO-GENERATED. DO NOT MODIFY.
// sqlapi v0.16.0; sqlgen v0.43.0

package vanilla

import (
	"context"
	"database/sql"
	"github.com/jackc/pgx"
	"github.com/rickb777/sqlapi/dialect"
	"github.com/rickb777/sqlapi/pgxapi"
	"github.com/rickb777/sqlapi/pgxapi/constraint"
	"github.com/rickb777/sqlapi/pgxapi/support"
	"github.com/rickb777/sqlapi/require"
)

// RecordTable holds a given table name with the database reference, providing access methods below.
// The Prefix field is often blank but can be used to hold a table name prefix (e.g. ending in '_'). Or it can
// specify the name of the schema, in which case it should have a trailing '.'.
type RecordTable struct {
	name        pgxapi.TableName
	database    pgxapi.Database
	db          pgxapi.Execer
	constraints constraint.Constraints
	ctx         context.Context
	pk          string
}

// Type conformance checks
var _ pgxapi.Table = &RecordTable{}
var _ pgxapi.Table = &RecordTable{}

// NewRecordTable returns a new table instance.
// If a blank table name is supplied, the default name "records" will be used instead.
// The request context is initialised with the background.
func NewRecordTable(name string, d pgxapi.Database) RecordTable {
	if name == "" {
		name = "records"
	}
	var constraints constraint.Constraints
	return RecordTable{
		name:        pgxapi.TableName{Name: name},
		database:    d,
		db:          d.DB(),
		constraints: constraints,
		ctx:         context.Background(),
		pk:          "id",
	}
}

// CopyTableAsRecordTable copies a table instance, retaining the name etc but
// providing methods appropriate for 'Record'. It doesn't copy the constraints of the original table.
//
// It serves to provide methods appropriate for 'Record'. This is most useful when this is used to represent a
// join result. In such cases, there won't be any need for DDL methods, nor Exec, Insert, Update or Delete.
func CopyTableAsRecordTable(origin pgxapi.Table) RecordTable {
	return RecordTable{
		name:        origin.Name(),
		database:    origin.Database(),
		db:          origin.DB(),
		constraints: nil,
		ctx:         context.Background(),
		pk:          "id",
	}
}

// SetPkColumn sets the name of the primary key column. It defaults to "id".
// The result is a modified copy of the table; the original is unchanged.
func (tbl RecordTable) SetPkColumn(pk string) RecordTable {
	tbl.pk = pk
	return tbl
}

// WithPrefix sets the table name prefix for subsequent queries.
// The result is a modified copy of the table; the original is unchanged.
func (tbl RecordTable) WithPrefix(pfx string) RecordTable {
	tbl.name.Prefix = pfx
	return tbl
}

// WithContext sets the context for subsequent queries via this table.
// The result is a modified copy of the table; the original is unchanged.
//
// The shared context in the *Database is not altered by this method. So it
// is possible to use different contexts for different (groups of) queries.
func (tbl RecordTable) WithContext(ctx context.Context) RecordTable {
	tbl.ctx = ctx
	return tbl
}

// Database gets the shared database information.
func (tbl RecordTable) Database() pgxapi.Database {
	return tbl.database
}

// Logger gets the trace logger.
func (tbl RecordTable) Logger() pgx.Logger {
	return tbl.database.Logger()
}

// WithConstraint returns a modified Table with added data consistency constraints.
func (tbl RecordTable) WithConstraint(cc ...constraint.Constraint) RecordTable {
	tbl.constraints = append(tbl.constraints, cc...)
	return tbl
}

// Constraints returns the table's constraints.
func (tbl RecordTable) Constraints() constraint.Constraints {
	return tbl.constraints
}

// Ctx gets the current request context.
func (tbl RecordTable) Ctx() context.Context {
	return tbl.ctx
}

// Dialect gets the database dialect.
func (tbl RecordTable) Dialect() dialect.Dialect {
	return tbl.database.Dialect()
}

// Name gets the table name.
func (tbl RecordTable) Name() pgxapi.TableName {
	return tbl.name
}

// PkColumn gets the column name used as a primary key.
func (tbl RecordTable) PkColumn() string {
	return tbl.pk
}

// DB gets the wrapped database handle, provided this is not within a transaction.
// Panics if it is in the wrong state - use IsTx() if necessary.
func (tbl RecordTable) DB() pgxapi.SqlDB {
	return tbl.db.(pgxapi.SqlDB)
}

// Execer gets the wrapped database or transaction handle.
func (tbl RecordTable) Execer() pgxapi.Execer {
	return tbl.db
}

// Tx gets the wrapped transaction handle, provided this is within a transaction.
// Panics if it is in the wrong state - use IsTx() if necessary.
func (tbl RecordTable) Tx() pgxapi.SqlTx {
	return tbl.db.(pgxapi.SqlTx)
}

// IsTx tests whether this is within a transaction.
func (tbl RecordTable) IsTx() bool {
	return pgxapi.IsTx(tbl.db)
}

// BeginTx starts a transaction using the table's context.
// This context is used until the transaction is committed or rolled back.
//
// If this context is cancelled, the sql package will roll back the transaction.
// In this case, Tx.Commit will then return an error.
//
// The provided TxOptions is optional and may be nil if defaults should be used.
// If a non-default isolation level is used that the driver doesn't support,
// an error will be returned.
//
// Panics if the Execer is not SqlDB.
func (tbl RecordTable) BeginTx(opts *pgx.TxOptions) (RecordTable, error) {
	var err error
	tbl.db, err = tbl.db.(pgxapi.SqlDB).BeginTx(tbl.ctx, opts)
	return tbl, tbl.logIfError(err)
}

// Using returns a modified Table using the transaction supplied. This is needed
// when making multiple queries across several tables within a single transaction.
// The result is a modified copy of the table; the original is unchanged.
func (tbl RecordTable) Using(tx pgxapi.SqlTx) RecordTable {
	tbl.db = tx
	return tbl
}

//func (tbl RecordTable) logQuery(query string, args ...interface{}) {
//	tbl.database.LogQuery(query, args...)
//}

func (tbl RecordTable) logError(err error) error {
	return tbl.database.LogError(err)
}

func (tbl RecordTable) logIfError(err error) error {
	return tbl.database.LogIfError(err)
}

//--------------------------------------------------------------------------------

// NumRecordColumns is the total number of columns in Record.
const NumRecordColumns = 1

// NumRecordDataColumns is the number of columns in Record not including the auto-increment key.
const NumRecordDataColumns = 0

// RecordColumnNames is the list of columns in Record.
const RecordColumnNames = "id"

// RecordDataColumnNames is the list of data columns in Record.
const RecordDataColumnNames = ""

//--------------------------------------------------------------------------------

// Query is the low-level request method for this table. The query is logged using whatever logger is
// configured. If an error arises, this too is logged.
//
// If you need a context other than the background, use WithContext before calling Query.
//
// The args are for any placeholder parameters in the query.
//
// The caller must call rows.Close() on the result.
//
// Wrap the result in *pgxapi.Rows if you need to access its data as a map.
func (tbl RecordTable) Query(query string, args ...interface{}) (pgxapi.SqlRows, error) {
	return support.Query(tbl, query, args...)
}

//--------------------------------------------------------------------------------

// QueryOneNullString is a low-level access method for one string. This can be used for function queries and
// such like. If the query selected many rows, only the first is returned; the rest are discarded.
// If not found, the result will be invalid.
//
// Note that this applies ReplaceTableName to the query string.
//
// The args are for any placeholder parameters in the query.
func (tbl RecordTable) QueryOneNullString(req require.Requirement, query string, args ...interface{}) (result sql.NullString, err error) {
	err = support.QueryOneNullThing(tbl, req, &result, query, args...)
	return result, err
}

// QueryOneNullInt64 is a low-level access method for one int64. This can be used for 'COUNT(1)' queries and
// such like. If the query selected many rows, only the first is returned; the rest are discarded.
// If not found, the result will be invalid.
//
// Note that this applies ReplaceTableName to the query string.
//
// The args are for any placeholder parameters in the query.
func (tbl RecordTable) QueryOneNullInt64(req require.Requirement, query string, args ...interface{}) (result sql.NullInt64, err error) {
	err = support.QueryOneNullThing(tbl, req, &result, query, args...)
	return result, err
}

// QueryOneNullFloat64 is a low-level access method for one float64. This can be used for 'AVG(...)' queries and
// such like. If the query selected many rows, only the first is returned; the rest are discarded.
// If not found, the result will be invalid.
//
// Note that this applies ReplaceTableName to the query string.
//
// The args are for any placeholder parameters in the query.
func (tbl RecordTable) QueryOneNullFloat64(req require.Requirement, query string, args ...interface{}) (result sql.NullFloat64, err error) {
	err = support.QueryOneNullThing(tbl, req, &result, query, args...)
	return result, err
}