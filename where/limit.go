package where

import (
	"bytes"
	"strconv"
)

// QueryConstraint is a value that is appended to a SELECT statement.
type QueryConstraint interface {
	Build() string
}

func BuildQueryConstraint(qc QueryConstraint) string {
	if qc == nil {
		return ""
	}
	return qc.Build()
}

//-------------------------------------------------------------------------------------------------

type literal string

// Literal returns the literal string supplied, converting it to a QueryConstraint.
func Literal(sqlPart string) QueryConstraint {
	return literal(sqlPart)
}

func (qc literal) Build() string {
	return string(qc)
}

//-------------------------------------------------------------------------------------------------

type queryConstraint struct {
	orderBy       []string
	desc          bool
	limit, offset int
}

var _ QueryConstraint = &queryConstraint{}

// OrderBy lists the column(s) by which the database will be asked to sort its results.
// The columns passed in here will be quoted according to the needs of the current dialect.
func OrderBy(column ...string) *queryConstraint {
	return &queryConstraint{orderBy: column}
}

// Limit sets the upper limit on the number of records to be returned.
func Limit(n int) *queryConstraint {
	return &queryConstraint{limit: n}
}

// Offset sets the offset into the result set; previous items will be discarded.
func Offset(n int) *queryConstraint {
	return &queryConstraint{offset: n}
}

// OrderBy lists the column(s) by which the database will be asked to sort its results.
// The columns passed in here will be quoted according to the needs of the current dialect.
func (qc *queryConstraint) OrderBy(column ...string) *queryConstraint {
	qc.orderBy = column
	return qc
}

// Asc sets the sort order to be ascending.
func (qc *queryConstraint) Asc() *queryConstraint {
	qc.desc = false
	return qc
}

// Desc sets the sort order to be descending.
func (qc *queryConstraint) Desc() *queryConstraint {
	qc.desc = true
	return qc
}

// Limit sets the upper limit on the number of records to be returned.
func (qc *queryConstraint) Limit(n int) *queryConstraint {
	qc.limit = n
	return qc
}

// Offset sets the offset into the result set; previous items will be discarded.
func (qc *queryConstraint) Offset(n int) *queryConstraint {
	qc.offset = n
	return qc
}

func (qc *queryConstraint) Build() string {
	b := &bytes.Buffer{}
	spacer := ""
	if len(qc.orderBy) > 0 {
		b.WriteString("ORDER BY ")
		separater := `"`
		for _, col := range qc.orderBy {
			b.WriteString(separater)
			b.WriteString(col)
			separater = `","`
		}
		b.WriteString(`"`)
		if qc.desc {
			b.WriteString(" DESC")
		}
		spacer = " "
	}
	if qc.limit > 0 {
		b.WriteString(spacer)
		b.WriteString("LIMIT ")
		b.WriteString(strconv.Itoa(qc.limit))
		spacer = " "
	}
	if qc.offset > 0 {
		b.WriteString(spacer)
		b.WriteString("OFFSET ")
		b.WriteString(strconv.Itoa(qc.offset))
	}
	return b.String()
}
