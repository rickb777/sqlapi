package dialect

import (
	"bytes"
	"github.com/rickb777/sqlapi/schema"
	"strings"
)

// Dialect is an abstraction of a type of database.
type Dialect interface {
	Index() int
	String() string
	Alias() string
	Quoter() Quoter
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
	QuoteW(w StringWriter, identifier string)
	SplitAndQuote(csv, before, separator, after string) string
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

// SplitAndQuote splits a comma-separated list and renders it with each word quoted.
// The result is prefixed by 'before' and suffixed by 'after'.
func (q quoter) SplitAndQuote(csv, before, separator, after string) string {
	identifiers := strings.Split(csv, ",")
	w := bytes.NewBuffer(make([]byte, 0, len(identifiers)*16))
	quoteW(w, before+string(q), string(q)+separator+string(q), string(q)+after, identifiers...)
	return w.String()
}

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
