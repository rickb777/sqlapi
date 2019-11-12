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
	ShowTables() string

	// ReplacePlaceholders alters a query string by replacing the '?' placeholders with the appropriate
	// placeholders needed by this dialect. For MySQL and SQlite3, the string is returned unchanged.
	ReplacePlaceholders(sql string, args []interface{}) string
	// Placeholders returns a comma-separated list of n placeholders.
	Placeholders(n int) string
	// HasNumberedPlaceholders returns true for dialects such as PostgreSQL that use numbered placeholders.
	HasNumberedPlaceholders() bool
	// HasLastInsertId returns true for dialects such as MySQL that return a last-insert ID after each
	// INSERT. This allows the corresponding feature of the database/sql API to work.
	// It is the inverse of InsertHasReturningPhrase.
	HasLastInsertId() bool
	// InsertHasReturningPhrase returns true for dialects such as Postgres that use a RETURNING phrase to
	// obtain the last-insert ID after each INSERT.
	// It is the inverse of HasLastInsertId.
	InsertHasReturningPhrase() bool
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
