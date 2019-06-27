package dialect

import (
	"fmt"
	"github.com/rickb777/sqlapi/schema"
	"github.com/rickb777/sqlapi/types"
	"github.com/rickb777/where/quote"
)

type sqlite struct {
	q quote.Quoter
}

var Sqlite Dialect = sqlite{q: quote.AnsiQuoter}

func (d sqlite) Index() int {
	return SqliteIndex
}

func (d sqlite) String() string {
	return "Sqlite"
}

func (d sqlite) Alias() string {
	return "SQLite3"
}

func (d sqlite) Quoter() quote.Quoter {
	return d.q
}

func (d sqlite) WithQuoter(q quote.Quoter) Dialect {
	d.q = q
	return d
}

// For integers, the value is a signed integer, stored in 1, 2, 3, 4, 6, or 8 bytes depending on the magnitude of the value
// For reals, the value is a floating point value, stored as an 8-byte IEEE floating point number.

func (dialect sqlite) FieldAsColumn(field *schema.Field) string {
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

func (dialect sqlite) InsertHasReturningPhrase() bool {
	return false
}

func (dialect sqlite) TruncateDDL(tableName string, force bool) []string {
	truncate := fmt.Sprintf("DELETE FROM %s", dialect.Quoter().Quote(tableName))
	return []string{truncate}
}

func (dialect sqlite) ShowTables() string {
	return `SELECT name FROM sqlite_master WHERE type = 'table'`
}

//-------------------------------------------------------------------------------------------------

func (dialect sqlite) HasNumberedPlaceholders() bool {
	return false
}

func (dialect sqlite) Placeholders(n int) string {
	return simpleQueryPlaceholders(n)
}

// ReplacePlaceholders converts a string containing '?' placeholders to
// the form used by MySQL and SQLite - i.e. unchanged.
func (dialect sqlite) ReplacePlaceholders(sql string, _ []interface{}) string {
	return sql
}

func (dialect sqlite) CreateTableSettings() string {
	return ""
}
