package dialect

import (
	"bytes"
	"github.com/rickb777/sqlapi/schema"
	"strings"
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
	Quoter() Quoter
	// WithQuoter returns a modified Dialect with a given quoter.
	WithQuoter(q Quoter) Dialect

	FieldAsColumn(field *schema.Field) string
	TruncateDDL(tableName string, force bool) []string
	CreateTableSettings() string
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

//-------------------------------------------------------------------------------------------------

// Quoter wraps identifiers in quote marks. Compound identifers (i.e. those with an alias prefix)
// are handled according to SQL grammar.
type Quoter interface {
	Quote(identifier string) string
	QuoteN(identifiers []string) []string
	QuoteW(w StringWriter, identifier string)
}

const (
	ansiQuoter  = quoter(`"`)
	mySqlQuoter = quoter("`")
	noQuoter    = quoter("")
)

var (
	// AnsiQuoter wraps identifiers in double quotes.
	AnsiQuoter Quoter = ansiQuoter

	// MySqlQuoter wraps identifiers in back-ticks.
	MySqlQuoter Quoter = mySqlQuoter

	// NoQuoter leaves identifiers unquoted.
	NoQuoter Quoter = noQuoter

	// DefaultQuoter is used by default. Change this to affect the default setting for every
	// SQL construction function.
	DefaultQuoter = AnsiQuoter
)

// NewQuoter gets a quoter using arbitrary quote marks.
func NewQuoter(mark string) Quoter {
	return quoter(mark)
}

// quoter wraps identifiers in quote marks. Compound identifers (i.e. those with an alias prefix)
// are handled according to SQL grammar.
type quoter string

// Quote renders an identifier within quote marks. If the identifier consists of both a
// prefix and a name, each part is quoted separately. For better performance, use QuoteW
// instead of Quote wherever possible.
func (q quoter) Quote(identifier string) string {
	if len(q) == 0 {
		return identifier
	}

	w := bytes.NewBuffer(make([]byte, 0, len(identifier)+4))
	q.QuoteW(w, identifier)
	return w.String()
}

// QuoteN quotes a list of identifiers using Quote.
func (q quoter) QuoteN(identifiers []string) []string {
	if len(q) == 0 {
		return identifiers
	}

	var r []string
	for _, id := range identifiers {
		r = append(r, q.Quote(id))
	}
	return r
}

// QuoteW renders an identifier within quote marks. If the identifier consists of both a
// prefix and a name, each part is quoted separately.
func (q quoter) QuoteW(w StringWriter, identifier string) {
	if len(q) == 0 {
		w.WriteString(identifier)
	} else {
		elements := strings.Split(identifier, ".")
		quoteW(w, string(q), string(q)+"."+string(q), string(q), elements...)
	}
}

func quoteW(w StringWriter, before, sep, after string, elements ...string) {
	if len(elements) > 0 {
		w.WriteString(before)
		for i, e := range elements {
			if i > 0 {
				w.WriteString(sep)
			}
			w.WriteString(e)
		}
		w.WriteString(after)
	}
}
