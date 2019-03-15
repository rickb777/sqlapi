package support

import (
	"database/sql"
	"github.com/benmoss/matchers"
	. "github.com/onsi/gomega"
	"github.com/rickb777/sqlapi"
	"github.com/rickb777/sqlapi/dialect"
	"github.com/rickb777/sqlapi/require"
	"github.com/rickb777/sqlapi/where"
	"testing"
)

func TestUpdateFieldsSQL(t *testing.T) {
	g := NewGomegaWithT(t)

	cases := []struct {
		quoter   dialect.Quoter
		expected string
	}{
		{
			quoter:   dialect.MySqlQuoter,
			expected: "UPDATE `foo` SET `col1`=?, `col2`=? WHERE (`room`=?) AND (`fun`=?)",
		},
		{
			quoter:   dialect.AnsiQuoter,
			expected: `UPDATE "foo" SET "col1"=?, "col2"=? WHERE ("room"=?) AND ("fun"=?)`,
		},
	}

	for _, c := range cases {
		f1 := sql.Named("col1", 111)
		f2 := sql.Named("col2", 222)
		wh := where.Eq("room", 101).And(where.Eq("fun", true))

		q, a := updateFieldsSQL("foo", c.quoter, wh, f1, f2)

		g.Expect(q).To(Equal(c.expected))
		g.Expect(a).To(matchers.DeepEqual([]interface{}{111, 222, 101, true}))
	}
}

func TestSliceSql(t *testing.T) {
	g := NewGomegaWithT(t)

	cases := []struct {
		dialect  dialect.Dialect
		expected string
	}{
		{
			dialect:  dialect.Mysql.WithQuoter(dialect.NoQuoter),
			expected: "SELECT foo FROM p.table WHERE (room=?) AND (fun=?) ORDER BY xyz",
		},
		{
			dialect:  dialect.Mysql,
			expected: "SELECT `foo` FROM `p`.`table` WHERE (`room`=?) AND (`fun`=?) ORDER BY `xyz`",
		},
		{
			dialect:  dialect.Mysql.WithQuoter(dialect.AnsiQuoter),
			expected: `SELECT "foo" FROM "p"."table" WHERE ("room"=?) AND ("fun"=?) ORDER BY "xyz"`,
		},
		{
			dialect:  dialect.Postgres,
			expected: `SELECT "foo" FROM "p"."table" WHERE ("room"=?) AND ("fun"=?) ORDER BY "xyz"`,
		},
	}

	for _, c := range cases {
		d := &StubDatabase{}
		tbl := StubTable{
			name: sqlapi.TableName{
				Prefix: "p.",
				Name:   "table",
			},
			dialect:  c.dialect,
			database: d,
		}
		wh := where.Eq("room", 101).And(where.Eq("fun", true))

		q, a := sliceSql(tbl, "foo", wh, where.OrderBy("xyz"))

		g.Expect(q).To(Equal(c.expected))
		g.Expect(a).To(matchers.DeepEqual([]interface{}{101, true}))
	}
}

func TestQuery_happy(t *testing.T) {
	g := NewGomegaWithT(t)

	d := &StubDatabase{}
	tbl := StubTable{
		name: sqlapi.TableName{
			Prefix: "p.",
			Name:   "table",
		},
		dialect:  dialect.Postgres,
		database: d,
	}

	_, err := Query(tbl, "SELECT foo FROM p.table WHERE x=?", 123)

	g.Expect(err).NotTo(HaveOccurred())
	g.Expect(d.loggedQueries).To(Equal([]string{"SELECT foo FROM p.table WHERE x=$1", "[123]"}))
}

func TestExec_happy(t *testing.T) {
	g := NewGomegaWithT(t)

	r := stubResult{ra: 2}
	e := stubExecer{stubResult: r}
	d := &StubDatabase{execer: e}
	tbl := StubTable{
		name: sqlapi.TableName{
			Prefix: "p.",
			Name:   "table",
		},
		dialect:  dialect.Postgres,
		database: d,
	}

	_, err := Exec(tbl, require.Exactly(2), "DELETE FROM p.table WHERE x=?", 123)

	g.Expect(err).NotTo(HaveOccurred())
	g.Expect(d.loggedQueries).To(Equal([]string{"DELETE FROM p.table WHERE x=$1", "[123]"}))
}
