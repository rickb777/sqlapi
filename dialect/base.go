package dialect

import (
	"bytes"
	"github.com/rickb777/sqlapi/schema"
	"strings"
)

func baseUpdateDML(table *schema.TableDescription) string {
	w := &bytes.Buffer{}
	w.WriteString(`"`)
	table.Fields.NoAuto().SqlNames().MkString3W(w, `"`, `"=?,"`, `"=? `)
	baseWhereClauseW(w, schema.FieldList{table.Primary}, 0)
	w.WriteString(`"`)
	return w.String()
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
