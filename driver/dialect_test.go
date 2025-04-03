package driver

import (
	"bytes"
	"testing"

	"github.com/rickb777/expect"
	"github.com/rickb777/sqlapi/schema"
	"github.com/rickb777/where/quote"
)

// assertion of conformance
var _ StringWriter = &bytes.Buffer{}

func TestQuote(t *testing.T) {
	cases := []struct {
		q        quote.Quoter
		expected string
	}{
		{quote.AnsiQuoter, `"x"."Aaaa"`},
		{quote.MySqlQuoter, "`x`.`Aaaa`"},
		{quote.NoQuoter, `x.Aaaa`},
	}
	for _, c := range cases {
		s1 := c.q.Quote("x.Aaaa")
		expect.String(s1).ToBe(t, c.expected)

		b2 := &bytes.Buffer{}
		c.q.QuoteW(b2, "x.Aaaa")
		expect.String(b2.String()).ToBe(t, c.expected)
	}
}

func TestSplitAndQuote(t *testing.T) {
	cases := []struct {
		q        quote.Quoter
		expected []string
	}{
		{quote.AnsiQuoter, []string{`"aa"`, `"bb"`, `"cc"`}},
		{quote.MySqlQuoter, []string{"`aa`", "`bb`", "`cc`"}},
		{quote.NoQuoter, []string{`aa`, `bb`, `cc`}},
	}
	for _, c := range cases {
		s1 := c.q.QuoteN([]string{"aa", "bb", "cc"})
		expect.Slice(s1).ToBe(t, c.expected...)
	}
}

func TestPlaceholders(t *testing.T) {
	cases := []struct {
		di       Dialect
		n        int
		expected string
	}{
		{Mysql(), 0, ""},
		{Mysql(), 1, "?"},
		{Mysql(), 3, "?,?,?"},
		{Mysql(), 11, "?,?,?,?,?,?,?,?,?,?,?"},

		{Postgres(), 0, ""},
		{Postgres(), 1, "$1"},
		{Postgres(), 3, "$1,$2,$3"},
		{Postgres(), 11, "$1,$2,$3,$4,$5,$6,$7,$8,$9,$10,$11"},
	}
	for _, c := range cases {
		s := c.di.Placeholders(c.n)
		expect.String(s).ToBe(t, c.expected)
	}
}

func TestReplacePlaceholders(t *testing.T) {
	s := Mysql().ReplacePlaceholders("?,?,?,?,?,?,?,?,?,?,?", nil)
	expect.String(s).ToBe(t, "?,?,?,?,?,?,?,?,?,?,?")

	s = Postgres().ReplacePlaceholders("?,?,?,?,?,?,?,?,?,?,?", nil)
	expect.String(s).ToBe(t, "$1,$2,$3,$4,$5,$6,$7,$8,$9,$10,$11")
}

func TestPickDialect(t *testing.T) {
	cases := []struct {
		di   Dialect
		name string
	}{
		{Mysql(), "MySQL"},
		{Postgres(), "Postgres"},
		{Postgres(), "PostgreSQL"},
		{Sqlite(), "SQLite"},
		{Sqlite(), "sqlite3"},
	}
	for _, c := range cases {
		s := PickDialect(c.name)
		expect.Any(s).ToBe(t, c.di)
	}
}

func TestFieldAsColumn(t *testing.T) {
	cases := []struct {
		di       Dialect
		field    *schema.Field
		expected string
	}{
		{Mysql(), id, "bigint not null primary key auto_increment"},
		{Mysql(), name, "varchar(2048) not null"},
		{Mysql(), active, "boolean not null"},
		{Mysql(), age, "int unsigned default null"},
		{Mysql(), bmi, "float default null"},
		{Mysql(), labels, "json"},

		{Postgres(), id, "bigserial not null primary key"},
		{Postgres(), name, "text not null"},
		{Postgres(), active, "boolean not null"},
		{Postgres(), age, "bigint default null"},
		{Postgres(), bmi, "real default null"},
		{Postgres(), labels, "json"},

		{Sqlite(), id, "integer not null primary key autoincrement"},
		{Sqlite(), name, "text not null"},
		{Sqlite(), active, "boolean not null"},
		{Sqlite(), age, "int unsigned default null"},
		{Sqlite(), bmi, "float default null"},
		{Sqlite(), labels, "text"},
	}
	for _, c := range cases {
		s := c.di.FieldAsColumn(c.field)
		expect.String(s).I(c.di.Name()).ToBe(t, c.expected)
	}
}
