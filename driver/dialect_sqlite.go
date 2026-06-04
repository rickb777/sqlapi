package driver

import (
	"fmt"

	"github.com/rickb777/sqlapi/schema"
	"github.com/rickb777/sqlapi/types"
	"github.com/rickb777/where/v2/dialect"
	"github.com/rickb777/where/v2/quote"
)

func Sqlite(q ...quote.Quoter) Dialect {
	if len(q) > 0 {
		dialect.SqliteQuoter = q[0]
	}

	return Dialect{d: dialect.Sqlite}
}

// For integers, the value is a signed integer, stored in 1, 2, 3, 4, 6, or 8 bytes depending on the magnitude of the value
// For reals, the value is a floating point value, stored as an 8-byte IEEE floating point number.

func sqliteFieldAsColumn(field *schema.Field) string {
	tags := field.GetTags()
	if tags.Auto {
		// In sqlite, "autoincrement" is less efficient than built-in "rowid"
		// and the datatype must be "integer" (https://sqlite.org/autoinc.html).
		return "integer not null primary key autoincrement"
	}

	switch field.Encode {
	case schema.ENCJSON:
		return "text"
	case schema.ENCTEXT:
		return "text"
	}

	column := "blob"
	dflt := tags.Default

	switch field.Type.Base {
	case types.Int, types.Int64:
		column = "bigint"
		dflt = tags.Default
	case types.Int8:
		column = "tinyint"
	case types.Int16:
		column = "smallint"
	case types.Int32:
		column = "int"
	case types.Uint, types.Uint64:
		column = "bigint unsigned"
	case types.Uint8:
		column = "tinyint unsigned"
	case types.Uint16:
		column = "smallint unsigned"
	case types.Uint32:
		column = "int unsigned"
	case types.Float32:
		column = "float"
	case types.Float64:
		column = "double"
	case types.Bool:
		column = "boolean"
	case types.String:
		column = "text"
		dflt = fmt.Sprintf("'%s'", tags.Default)
	}

	return fieldTags(field.Type.IsPtr, tags, column, dflt)
}

func fieldTags(fieldIsPtr bool, tags types.Tag, column, dflt string) string {
	if fieldIsPtr {
		column += " default null"
	} else {
		column += " not null"

		if tags.Default != "" {
			column += " default " + dflt
		}

	}

	if tags.Primary {
		column += " primary key"
	}

	return column
}

func sqliteTruncateDDL(tableName string, force bool) []string {
	truncate := fmt.Sprintf("DELETE FROM %s", dialect.SqliteQuoter.Quote(tableName))
	return []string{truncate}
}

//-------------------------------------------------------------------------------------------------

const (
	sqliteShowTables              = `SELECT name FROM sqlite_master WHERE type = 'table'`
	sqliteHasNumberedPlaceholders = false
	sqliteHasLastInsertId         = true
	sqliteCreateTableSettings     = ""
)

func sqlitePlaceholders(n int) string {
	return simpleQueryPlaceholders(n)
}

// ReplacePlaceholders converts a string containing '?' placeholders to
// the form used by MySQL and SQLite - i.e. unchanged.
func sqliteReplacePlaceholders(sql string, _ []interface{}) string {
	return sql
}
