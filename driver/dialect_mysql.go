package driver

import (
	"fmt"

	"github.com/rickb777/sqlapi/schema"
	"github.com/rickb777/sqlapi/types"
	"github.com/rickb777/where/dialect"
	"github.com/rickb777/where/quote"
)

type mysql struct {
	d dialect.DialectConfig
}

func of(dflt dialect.DialectConfig, d ...dialect.DialectConfig) dialect.DialectConfig {
	if len(d) > 0 {
		return d[0]
	}
	return dflt
}

func Mysql(d ...dialect.DialectConfig) Dialect {
	return mysql{d: of(dialect.MysqlConfig, d...)}
}

func (d mysql) Index() dialect.Dialect {
	return dialect.Mysql
}

func (d mysql) String() string {
	if d.d.Quoter != nil {
		return fmt.Sprintf("Mysql/%s", d.d.Quoter)
	}
	return "Mysql"
}

func (d mysql) Name() string {
	return "Mysql"
}

func (d mysql) Alias() string {
	return "MySQL"
}

func (d mysql) Config() dialect.DialectConfig {
	return d.d
}

func (d mysql) Quoter() quote.Quoter {
	return d.d.Quoter
}

func (d mysql) WithQuoter(q quote.Quoter) Dialect {
	d.d.Quoter = q
	return d
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

func (dialect mysql) HasLastInsertId() bool {
	return true
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
