package driver

import (
	"fmt"

	"github.com/rickb777/sqlapi/schema"
	"github.com/rickb777/sqlapi/types"
	"github.com/rickb777/where/v2/dialect"
	"github.com/rickb777/where/v2/quote"
)

func Mysql(q ...quote.Quoter) Dialect {
	if len(q) > 0 {
		dialect.MySqlQuoter = q[0]
	}
	return Dialect{d: dialect.Mysql}
}

// see https://dev.mysql.com/doc/refman/5.7/en/data-types.html

func mysqlFieldAsColumn(field *schema.Field) string {
	tags := field.GetTags()
	indexed := len(tags.Index) > 0 || len(tags.Unique) > 0

	switch field.Encode {
	case schema.ENCJSON:
		return "json"
	case schema.ENCTEXT:
		return varchar(tags.Size, indexed)
	}

	column := "mediumblob"
	dflt := tags.Default

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
		column = "boolean"
	case types.String:
		column = varchar(tags.Size, indexed)
		dflt = fmt.Sprintf("'%s'", tags.Default)
	}

	column = fieldTags(field.Type.IsPtr, tags, column, dflt)

	if tags.Auto {
		column += " auto_increment"
	}

	return column
}

func varchar(size int, indexed bool) string {
	if size == 0 { // unspecified
		if indexed {
			// largest variable size that has only one length byte
			return "varchar(255)"
		}
		return "text"
	}
	if size >= 2<<24 {
		return "longtext"
	}
	if size >= 2<<16 {
		return "mediumtext"
	}
	return fmt.Sprintf("varchar(%d)", size)
}

// see https://dev.mysql.com/doc/refman/5.7/en/integer-types.html

func mysqlTruncateDDL(tableName string, force bool) []string {
	truncate := fmt.Sprintf("TRUNCATE %s", dialect.MySqlQuoter.Quote(tableName))
	if !force {
		return []string{truncate}
	}

	return []string{
		"SET FOREIGN_KEY_CHECKS=0",
		truncate,
		"SET FOREIGN_KEY_CHECKS=1",
	}
}

//-------------------------------------------------------------------------------------------------

const (
	mysqlShowTables              = `SHOW TABLES`
	mysqlHasNumberedPlaceholders = false
	mysqlHasLastInsertId         = true
	mysqlCreateTableSettings     = " ENGINE=InnoDB DEFAULT CHARSET=utf8"
)

func mysqlPlaceholders(n int) string {
	return simpleQueryPlaceholders(n)
}

// ReplacePlaceholders converts a string containing '?' placeholders to
// the form used by MySQL and SQLite - i.e. unchanged.
func mysqlReplacePlaceholders(sql string, _ []interface{}) string {
	return sql
}
