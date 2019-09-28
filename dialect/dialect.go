package dialect

import (
	"strings"

	"github.com/rickb777/sqlapi/schema"
	"github.com/rickb777/where/quote"
)

// Dialect is an abstraction of a type of database.
type Dialect interface {
	// Index returns a consistent ID for this dialect, regardless of other settings.
	Index() int
	// String returns the name of this dialect.
	String() string
	// Alias is an alternative name for this dialect.
	Alias() string
	// Quoter is the tool used for quoting identifiers.
	Quoter() quote.Quoter
	// WithQuoter returns a modified Dialect with a given quoter.
	WithQuoter(q quote.Quoter) Dialect

	FieldAsColumn(field *schema.Field) string
	TruncateDDL(tableName string, force bool) []string
	CreateTableSettings() string
	InsertHasReturningPhrase() bool
	ShowTables() string

	ReplacePlaceholders(sql string, args []interface{}) string
	Placeholders(n int) string
	HasNumberedPlaceholders() bool
	HasLastInsertId() bool
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

type dialect string
