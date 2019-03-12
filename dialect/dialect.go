package dialect

import (
	"github.com/rickb777/sqlapi/schema"
	"strings"
)

// Dialect is an abstraction of a type of database.
type Dialect interface {
	Index() int
	String() string
	Alias() string

	TableDDL(*schema.TableDescription) string
	FieldDDL(w StringWriter, field *schema.Field, comma string) string
	TruncateDDL(tableName string, force bool) []string
	CreateTableSettings() string
	FieldAsColumn(*schema.Field) string
	InsertHasReturningPhrase() bool
	ShowTables() string

	ReplacePlaceholders(sql string, args []interface{}) string
	Placeholders(n int) string
	HasNumberedPlaceholders() bool
}

//-------------------------------------------------------------------------------------------------

const (
	SqliteIndex = iota
	MysqlIndex
	PostgresIndex
	PgxIndex
)

//-------------------------------------------------------------------------------------------------

// AllDialects lists all currently-supported dialects.
var AllDialects = []Dialect{Sqlite, Mysql, Postgres, Pgx}

// PickDialect finds a dialect that matches by name, ignoring letter case.
// It returns nil if not found.
func PickDialect(name string) Dialect {
	for _, d := range AllDialects {
		if strings.EqualFold(name, d.String()) || strings.EqualFold(name, d.Alias()) {
			return d
		}
	}
	return nil
}
