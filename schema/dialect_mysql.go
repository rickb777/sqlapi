package schema

import (
	"fmt"
	"github.com/rickb777/sqlapi/types"
	"io"
)

type mysql struct{}

var Mysql Dialect = mysql{}

func (d mysql) Index() int {
	return MysqlIndex
}

func (d mysql) String() string {
	return "Mysql"
}

func (d mysql) Alias() string {
	return "MySQL"
}

// see https://dev.mysql.com/doc/refman/5.7/en/data-types.html

func (dialect mysql) FieldAsColumn(field *Field) string {
	switch field.Encode {
	case ENCJSON:
		return "json"
	case ENCTEXT:
		return varchar(field.Tags.Size)
	}

	column := "mediumblob"
	dflt := field.Tags.Default

	switch field.Type.Base {
	case types.Int, types.Int64:
		column = "bigint"
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
		column = "tinyint(1)"
	case types.String:
		column = varchar(field.Tags.Size)
		dflt = fmt.Sprintf("'%s'", field.Tags.Default)
	}

	column = fieldTags(field, column, dflt)

	if field.Tags.Auto {
		column += " auto_increment"
	}

	return column
}

func varchar(size int) string {
	// Assigns an arbitrary size if none is provided.
	// 255 is chosen because the max. index key length is 767 and UTF8 might use up to three bytes per character.
	if size == 0 {
		size = 255
	}
	return fmt.Sprintf("varchar(%d)", size)
}

// see https://dev.mysql.com/doc/refman/5.7/en/integer-types.html

func (dialect mysql) TableDDL(table *TableDescription) string {
	return baseTableDDL(table, dialect, " \"\\n\"+\n", `"`)
}

func (dialect mysql) FieldDDL(w io.Writer, field *Field, comma string) string {
	return backTickFieldDDL(w, field, comma, dialect)
}

func (dialect mysql) InsertHasReturningPhrase() bool {
	return false
}

func (dialect mysql) UpdateDML(table *TableDescription) string {
	return baseUpdateDML(table, backTickQuoted, baseParamIsQuery)
}

func (dialect mysql) TruncateDDL(tableName string, force bool) []string {
	truncate := fmt.Sprintf("TRUNCATE %s", tableName)
	if !force {
		return []string{truncate}
	}

	return []string{
		"SET FOREIGN_KEY_CHECKS=0",
		truncate,
		"SET FOREIGN_KEY_CHECKS=1",
	}
}

func (dialect mysql) SplitAndQuote(csv string) string {
	return baseSplitAndQuote(csv, "`", "`,`", "`")
}

func (dialect mysql) Quote(identifier string) string {
	return backTickQuoted(identifier)
}

func (dialect mysql) QuoteW(w io.Writer, identifier string) {
	backTickQuotedW(w, identifier)
}

func (dialect mysql) QuoteWithPlaceholder(w io.Writer, identifier string, idx int) {
	backTickQuotedW(w, identifier)
	io.WriteString(w, "=?")
}

func (dialect mysql) Placeholder(name string, j int) string {
	return "?"
}

func (dialect mysql) Placeholders(n int) string {
	return baseQueryPlaceholders(n)
}

// ReplacePlaceholders converts a string containing '?' placeholders to
// the form used by MySQL and SQLite - i.e. unchanged.
func (dialect mysql) ReplacePlaceholders(sql string, _ []interface{}) string {
	return sql
}

func (dialect mysql) CreateTableSettings() string {
	return " ENGINE=InnoDB DEFAULT CHARSET=utf8"
}
