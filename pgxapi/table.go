package pgxapi

import (
	"github.com/rickb777/sqlapi/dialect"
)

// Table provides the generic features of each generated table handler.
type Table interface {
	// Name gets the table name. without prefix
	Name() TableName

	// Database gets the shared database information.
	Database() Database

	// Execer gets the wrapped database or transaction handle.
	Execer() Execer

	// DB gets the wrapped database handle, provided this is not within a transaction.
	// Panics if it is in the wrong state - use IsTx() if necessary.
	DB() SqlDB

	// Tx gets the wrapped transaction handle, provided this is within a transaction.
	// Panics if it is in the wrong state - use IsTx() if necessary.
	Tx() SqlTx

	// IsTx tests whether this is within a transaction.
	IsTx() bool

	// Constraints returns the table's constraints.
	// (not included here because of package inter-dependencies)
	//Constraints() constraint.Constraints

	// Dialect gets the database dialect.
	Dialect() dialect.Dialect

	// Logger gets the trace logger.
	Logger() Logger

	//---------------------------------------------------------------------------------------------
	// The following type-specific methods are also provided (but are not part of this interface).

	// WithPrefix sets the table name prefix for subsequent queries.
	// The result is a modified copy of the table; the original is unchanged.
	// WithPrefix(pfx string) SomeTypeTable

	// WithLogger sets the logger for subsequent queries. An alias for SetLogger.
	// The result is a modified copy of the table; the original is unchanged.
	//WithLogger(logger *log.Logger) SomeTypeTable {

	// Begin starts a transaction. The default isolation level is dependent on the driver.
	// The result is a modified copy of the table; the original is unchanged.
	//BeginTx(opts *sql.TxOptions) (SomeTypeTable, error)

	// Using returns a modified Table using the transaction supplied. This is needed
	// when making multiple queries across several tables within a single transaction.
	// The result is a modified copy of the table; the original is unchanged.
	//Using(tx *sql.Tx) SomeTypeTable

	// WithConstraint returns a modified Table with added data consistency constraints.
	//WithConstraint(cc ...sqlgen2.Constraint) SomeTypeTable {
	//---------------------------------------------------------------------------------------------

	// Query is the low-level request method for this table. The query is logged using whatever logger is
	// configured. If an error arises, this too is logged.
	//
	// If you need a context other than the background, use WithContext before calling Query.
	//
	// The args are for any placeholder parameters in the query.
	//
	// The caller must call rows.Close() on the result.
	//Query(query string, args ...interface{}) ([]Values, error)
}

// TableCreator is a table with create/delete/truncate methods.
type TableCreator interface {
	Table

	// CreateTable creates the database table.
	CreateTable(ifNotExists bool) (int64, error)

	// DropTable drops the database table.
	DropTable(ifExists bool) (int64, error)

	// Truncate empties the table
	Truncate(force bool) (err error)
}

// TableWithIndexes is a table creator with create/delete methods for the indexes.
type TableWithIndexes interface {
	TableCreator

	// CreateIndexes creates the indexes for the database table.
	CreateIndexes(ifNotExist bool) (err error)

	// DropIndexes executes a query that drops all the indexes on the database table.
	DropIndexes(ifExist bool) (err error)

	// CreateTableWithIndexes creates the database table and its indexes.
	CreateTableWithIndexes(ifNotExist bool) (err error)
}
