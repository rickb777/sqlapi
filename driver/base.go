package driver

import "strings"

const placeholders = "?,?,?,?,?,?,?,?,?,?"

func simpleQueryPlaceholders(n int) string {
	if n == 0 {
		return ""
	} else if n <= 10 {
		m := (n * 2) - 1
		return placeholders[:m]
	}
	return strings.Repeat("?,", n-1) + "?"
}

func baseFieldAsColumn(w StringWriter, name, field string) {
	w.WriteString("\t\"")
	w.WriteString(name)
	w.WriteString("\":\t\"")
	w.WriteString(field)
	w.WriteString("\"")
}
