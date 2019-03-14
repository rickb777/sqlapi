package dialect

import (
	"fmt"
	"github.com/rickb777/sqlapi/schema"
	"github.com/rickb777/sqlapi/types"
)

type mysql quoter

var Mysql Dialect = mysql(mySqlQuoter)

func (d mysql) Index() int {
	return MysqlIndex
}

func (d mysql) String() string {
	return "Mysql"
}

func (d mysql) Alias() string {
	return "MySQL"
}

func (d mysql) Quoter() Quoter {
	return quoter(d)
}

func (d mysql) WithQuoter(q Quoter) Dialect {
	return mysql(q.(quoter))
}

// see https://dev.mysql.com/doc/refman/5.7/en/data-types.html

func (dialect mysql) FieldAsColumn(field *schema.Field) string {
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

func (dialect mysql) InsertHasReturningPhrase() bool {
	return false
}

func (dialect mysql) TruncateDDL(tableName string, force bool) []string {
	truncate := fmt.Sprintf("TRUNCATE %s", dialect.Quoter().Quote(tableName))
	if !force {
		return []string{truncate}
	}

	return []string{
		"SET FOREIGN_KEY_CHECKS=0",
		truncate,
		"SET FOREIGN_KEY_CHECKS=1",
	}
}

func (dialect mysql) ShowTables() string {
	return `SHOW TABLES`
}

//-------------------------------------------------------------------------------------------------

func (dialect mysql) HasNumberedPlaceholders() bool {
	return false
}

func (dialect mysql) Placeholders(n int) string {
	return simpleQueryPlaceholders(n)
}

// ReplacePlaceholders converts a string containing '?' placeholders to
// the form used by MySQL and SQLite - i.e. unchanged.
func (dialect mysql) ReplacePlaceholders(sql string, _ []interface{}) string {
	return sql
}

func (dialect mysql) CreateTableSettings() string {
	return " ENGINE=InnoDB DEFAULT CHARSET=utf8"
}
