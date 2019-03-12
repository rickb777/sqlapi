package dialect

import (
	"bytes"
	"github.com/rickb777/sqlapi/schema"
	"strings"
	"text/tabwriter"
)

// Table returns a SQL statement to create the table.
func baseTableDDL(table *schema.TableDescription, dialect Dialect, initial, final string) string {

	// use a large default buffer size of so that
	// the tabbing doesn't get prematurely flushed
	// resulting in un-even lines.
	var byt = make([]byte, 0, 100*len(table.Fields))
	var buf = bytes.NewBuffer(byt)

	// use a tab writer to evenly space the column
	// names and column types.
	var tab = tabwriter.NewWriter(buf, 0, 8, 1, ' ', 0)
	w := Adapt(tab)
	comma := initial
	for _, field := range table.Fields {
		comma = dialect.FieldDDL(w, field, comma)
	}
	w.WriteString(final)

	// flush the tab writer to write to the buffer
	tab.Flush()

	return buf.String()
}

func baseFieldDDL(w StringWriter, field *schema.Field, comma string, dialect Dialect) string {
	w.WriteString(comma)
	w.WriteString("\t\"")
	w.WriteString(string(field.SqlName))
	w.WriteString("\"\t")
	w.WriteString(dialect.FieldAsColumn(field))
	return ",\n" // for next iteration
}

func baseUpdateDML(table *schema.TableDescription) string {
	w := &bytes.Buffer{}
	w.WriteString(`"`)
	table.Fields.NoAuto().SqlNames().MkString3W(w, `"`, `"=?,"`, `"=? `)
	baseWhereClauseW(w, schema.FieldList{table.Primary}, 0)
	w.WriteString(`"`)
	return w.String()
}

// Quote renders an identifier within double quotes. If the identifier consists of both a
// prefix and a name, each part is quoted separately.
func Quote(identifier string) string {
	w := bytes.NewBuffer(make([]byte, 0, len(identifier)+4))
	QuoteW(w, identifier)
	return w.String()
}

// QuoteW renders an identifier within double quotes. If the identifier consists of both a
// prefix and a name, each part is quoted separately.
func QuoteW(w StringWriter, identifier string) {
	elements := strings.Split(identifier, ".")
	doubleQuoteW(w, ".", elements...)
}

func DoubleQuotedList(csv string) string {
	identifiers := strings.Split(csv, ",")
	w := bytes.NewBuffer(make([]byte, 0, len(identifiers)*16))
	doubleQuoteW(w, ",", identifiers...)
	return w.String()
}

func doubleQuoteW(w StringWriter, sep string, elements ...string) {
	if len(elements) > 0 {
		w.WriteString(`"`)
		for i, e := range elements {
			if i > 0 {
				w.WriteString(`"`)
				w.WriteString(sep)
				w.WriteString(`"`)
			}
			w.WriteString(e)
		}
		w.WriteString(`"`)
	}
}

//-------------------------------------------------------------------------------------------------

const placeholders = "?,?,?,?,?,?,?,?,?,?"

func baseQueryPlaceholders(n int) string {
	if n == 0 {
		return ""
	} else if n <= 10 {
		m := (n * 2) - 1
		return placeholders[:m]
	}
	return strings.Repeat("?,", n-1) + "?"
}

// helper function to generate the Where clause
// section of a SQL statement
func baseWhereClauseW(w StringWriter, fields schema.FieldList, pos int) {
	j := pos

	for i, field := range fields {
		switch {
		case i == 0:
			w.WriteString("WHERE")
		default:
			w.WriteString("\nAND")
		}

		w.WriteString(` "`)
		w.WriteString(field.SqlName)
		w.WriteString(`"=?`)

		j++
	}
}
