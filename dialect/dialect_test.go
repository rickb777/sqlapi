package dialect

import (
	"bytes"
	. "github.com/onsi/gomega"
	"github.com/rickb777/sqlapi/schema"
	"testing"
)

// assertion of conformance
var _ StringWriter = &bytes.Buffer{}

func TestQuote(t *testing.T) {
	g := NewGomegaWithT(t)

	cases := []struct {
		q        Quoter
		expected string
	}{
		{AnsiQuoter, `"x"."Aaaa"`},
		{MySqlQuoter, "`x`.`Aaaa`"},
		{NoQuoter, `x.Aaaa`},
	}
	for _, c := range cases {
		s1 := c.q.Quote("x.Aaaa")
		g.Expect(s1).Should(Equal(c.expected))

		b2 := &bytes.Buffer{}
		c.q.QuoteW(b2, "x.Aaaa")
		g.Expect(b2.String()).Should(Equal(c.expected))
	}
}

func TestSplitAndQuote(t *testing.T) {
	g := NewGomegaWithT(t)

	cases := []struct {
		q        Quoter
		expected string
	}{
		{AnsiQuoter, `["aa", "bb", "cc"]`},
		{MySqlQuoter, "[`aa`, `bb`, `cc`]"},
		{NoQuoter, `[aa, bb, cc]`},
	}
	for _, c := range cases {
		s1 := c.q.SplitAndQuote("aa,bb,cc", "[", ", ", "]")
		g.Expect(s1).Should(Equal(c.expected))
	}
}

func TestPlaceholders(t *testing.T) {
	g := NewGomegaWithT(t)

	cases := []struct {
		di       Dialect
		n        int
		expected string
	}{
		{Mysql, 0, ""},
		{Mysql, 1, "?"},
		{Mysql, 3, "?,?,?"},
		{Mysql, 11, "?,?,?,?,?,?,?,?,?,?,?"},

		{Postgres, 0, ""},
		{Postgres, 1, "$1"},
		{Postgres, 3, "$1,$2,$3"},
		{Postgres, 11, "$1,$2,$3,$4,$5,$6,$7,$8,$9,$10,$11"},
	}
	for _, c := range cases {
		s := c.di.Placeholders(c.n)
		g.Expect(s).Should(Equal(c.expected))
	}
}

func TestReplacePlaceholders(t *testing.T) {
	g := NewGomegaWithT(t)

	s := Mysql.ReplacePlaceholders("?,?,?,?,?,?,?,?,?,?,?", nil)
	g.Expect(s).Should(Equal("?,?,?,?,?,?,?,?,?,?,?"))

	s = Postgres.ReplacePlaceholders("?,?,?,?,?,?,?,?,?,?,?", nil)
	g.Expect(s).Should(Equal("$1,$2,$3,$4,$5,$6,$7,$8,$9,$10,$11"))
}

func TestPickDialect(t *testing.T) {
	g := NewGomegaWithT(t)

	cases := []struct {
		di   Dialect
		name string
	}{
		{Mysql, "MySQL"},
		{Postgres, "Postgres"},
		{Postgres, "PostgreSQL"},
		{Sqlite, "SQLite"},
		{Sqlite, "sqlite3"},
	}
	for _, c := range cases {
		s := PickDialect(c.name)
		g.Expect(s).Should(Equal(c.di))
	}
}

func TestFieldAsColumn(t *testing.T) {
	g := NewGomegaWithT(t)

	cases := []struct {
		di       Dialect
		field    *schema.Field
		expected string
	}{
		{Mysql, id, "\t\"id\"\tbigint not null primary key auto_increment"},
		{Mysql, name, "\t\"username\"\tvarchar(2048) not null"},
		{Mysql, active, "\t\"active\"\tboolean not null"},
		{Mysql, age, "\t\"age\"\tint unsigned default null"},
		{Mysql, bmi, "\t\"bmi\"\tfloat default null"},
		{Mysql, labels, "\t\"labels\"\tjson"},

		{Postgres, id, "\t\"id\"\tbigserial not null primary key"},
		{Postgres, name, "\t\"username\"\ttext not null"},
		{Postgres, active, "\t\"active\"\tboolean not null"},
		{Postgres, age, "\t\"age\"\tbigint default null"},
		{Postgres, bmi, "\t\"bmi\"\treal default null"},
		{Postgres, labels, "\t\"labels\"\tjson"},

		{Sqlite, id, "\t\"id\"\tinteger not null primary key autoincrement"},
		{Sqlite, name, "\t\"username\"\ttext not null"},
		{Sqlite, active, "\t\"active\"\tboolean not null"},
		{Sqlite, age, "\t\"age\"\tint unsigned default null"},
		{Sqlite, bmi, "\t\"bmi\"\tfloat default null"},
		{Sqlite, labels, "\t\"labels\"\ttext"},
	}
	for _, c := range cases {
		b := &bytes.Buffer{}
		c.di.FieldAsColumn(b, c.field)
		g.Expect(b.String()).Should(Equal(c.expected), c.di.String())
	}
}
