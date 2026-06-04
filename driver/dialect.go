package driver

import (
	"fmt"

	"github.com/rickb777/sqlapi/schema"
	"github.com/rickb777/where/v2/dialect"
	"github.com/rickb777/where/v2/quote"
)

// Dialect is an abstraction of a type of database.
type Dialect struct {
	d dialect.Dialect
}

// Index returns a consistent ID for this dialect, regardless of other settings.
func (d Dialect) Index() dialect.Dialect {
	return dialect.Mysql
}

// String returns the name (and quoter if present) of this dialect.
func (d Dialect) String() string {
	return fmt.Sprintf("%s/%s", d.d.String(), d.d.Quoter())
}

// Name returns the name of this dialect.
func (d Dialect) Name() string {
	return d.d.String()
}

// Alias is an alternative name for this dialect.
func (d Dialect) Alias() string {
	switch d.d {
	case dialect.Mysql:
		return "MySQL"
	case dialect.Postgres:
		return "PostgreSQL"
	case dialect.Sqlite:
		return "SQLite3"
		//case dialect.SqlServer:
		//	return MSSqlQuoter
	}
	panic(d.d)
}

// Format is the dialect format option.
func (d Dialect) Format() dialect.FormatOption {
	return d.d.Placeholder()
}

// Quoter is the tool used for quoting identifiers.
func (d Dialect) Quoter() quote.Quoter {
	return d.Quoter()
}
func (d Dialect) FieldAsColumn(field *schema.Field) string {
	switch d.d {
	case dialect.Mysql:
		return mysqlFieldAsColumn(field)
	case dialect.Postgres:
		return postgresFieldAsColumn(field)
	case dialect.Sqlite:
		return ""
	}
	panic(d.d)
}

func (d Dialect) InsertHasReturningPhrase() bool {
	switch d.d {
	case dialect.Mysql:
		return false
	case dialect.Postgres:
		return true
	case dialect.Sqlite:
		return false
	}
	panic(d.d)
}

func (d Dialect) TruncateDDL(tableName string, force bool) []string {
	switch d.d {
	case dialect.Mysql:
		return mysqlTruncateDDL(tableName, force)
	case dialect.Postgres:
		return postgresTruncateDDL(tableName, force)
	case dialect.Sqlite:
		return sqliteTruncateDDL(tableName, force)
	}
	panic(d.d)
}

func (d Dialect) ShowTables() string {
	switch d.d {
	case dialect.Mysql:
		return mysqlShowTables
	case dialect.Postgres:
		return showTableNamePostgres
	case dialect.Sqlite:
		return sqliteShowTables
		//case dialect.SqlServer:
		//	return MSSqlQuoter
	}
	panic(d.d)
}

//CreateTableSettings() string
//
//// ReplacePlaceholders alters a query string by replacing the '?' placeholders with the appropriate
//// placeholders needed by this dialect. For MySQL and SQlite3, the string is returned unchanged.
//ReplacePlaceholders(sql string, args []interface{}) string
//// Placeholders returns a comma-separated list of n placeholders.
//Placeholders(n int) string
//// HasNumberedPlaceholders returns true for dialects such as PostgreSQL that use numbered placeholders.
//HasNumberedPlaceholders() bool
//// HasLastInsertId returns true for dialects such as MySQL that return a last-insert ID after each
//// INSERT. This allows the corresponding feature of the database/sql API to work.
//// It is the inverse of InsertHasReturningPhrase.
//HasLastInsertId() bool
//// InsertHasReturningPhrase returns true for dialects such as Postgres that use a RETURNING phrase to
//// obtain the last-insert ID after each INSERT.
//// It is the inverse of HasLastInsertId.
//InsertHasReturningPhrase() bool
