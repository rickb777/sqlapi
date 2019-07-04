package support

import (
	"database/sql"
	"github.com/benmoss/matchers"
	. "github.com/onsi/gomega"
	"github.com/rickb777/sqlapi"
	"github.com/rickb777/sqlapi/dialect"
	"github.com/rickb777/sqlapi/require"
	"github.com/rickb777/where"
	"github.com/rickb777/where/quote"
	"testing"
)

func xTestUpdateFieldsSQL(t *testing.T) {
	g := NewGomegaWithT(t)

	cases := []struct {
		quoter   quote.Quoter
		expected string
	}{
		{
			quoter:   quote.MySqlQuoter,
			expected: "UPDATE `foo` SET `col1`=?, `col2`=? WHERE (`room`=?) AND (`fun`=?)",
		},
		{
			quoter:   quote.AnsiQuoter,
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

func xTestSliceSql(t *testing.T) {
	g := NewGomegaWithT(t)

	cases := []struct {
		dialect  dialect.Dialect
		expected string
	}{
		{
			dialect:  dialect.Mysql.WithQuoter(quote.NoQuoter),
			expected: "SELECT foo FROM p.table WHERE (room=?) AND (fun=?) ORDER BY xyz",
		},
		{
			dialect:  dialect.Mysql,
			expected: "SELECT `foo` FROM `p`.`table` WHERE (`room`=?) AND (`fun`=?) ORDER BY `xyz`",
		},
		{
			dialect:  dialect.Mysql.WithQuoter(quote.AnsiQuoter),
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

	d := &StubDatabase{
		execer: stubExecer{},
		stdLog: &stubLogger{},
	}
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
	g.Expect(d.stdLog.logged).To(ConsistOf(`SELECT foo FROM p.table WHERE x=$1 [123]` + "\n"))
}

func TestExec_happy(t *testing.T) {
	g := NewGomegaWithT(t)

	e := stubExecer{stubResult: 2}
	d := &StubDatabase{execer: e, stdLog: &stubLogger{}}
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
	g.Expect(d.stdLog.logged).To(ConsistOf(`DELETE FROM p.table WHERE x=$1 [123]` + "\n"))
}

func TestGetIntIntIndex_happy(t *testing.T) {
	g := NewGomegaWithT(t)

	e := stubExecer{rows: &StubRows{
		Rows: []StubRow{{int64(2), int64(16)}, {int64(3), int64(81)}},
	}}
	d := &StubDatabase{execer: e, stdLog: &stubLogger{}}
	tbl := StubTable{
		name: sqlapi.TableName{
			Prefix: "p.",
			Name:   "table",
		},
		dialect:  dialect.Postgres,
		database: d,
	}

	m, err := GetIntIntIndex(tbl, quote.AnsiQuoter, "aa", "bb", where.Eq("foo", "bar"))

	g.Expect(err).NotTo(HaveOccurred())
	g.Expect(m).To(Equal(map[int64]int64{2: 16, 3: 81}))
	g.Expect(d.stdLog.logged).To(ConsistOf(`SELECT "aa", "bb" FROM "p"."table" WHERE "foo"=$1 [bar]` + "\n"))
}

func TestGetStringIntIndex_happy(t *testing.T) {
	g := NewGomegaWithT(t)

	e := stubExecer{rows: &StubRows{
		Rows: []StubRow{{"two", int64(16)}, {"three", int64(81)}},
	}}
	d := &StubDatabase{execer: e, stdLog: &stubLogger{}}
	tbl := StubTable{
		name: sqlapi.TableName{
			Prefix: "p.",
			Name:   "table",
		},
		dialect:  dialect.Postgres,
		database: d,
	}

	m, err := GetStringIntIndex(tbl, quote.AnsiQuoter, "aa", "bb", where.Eq("foo", "bar"))

	g.Expect(err).NotTo(HaveOccurred())
	g.Expect(m).To(Equal(map[string]int64{"two": 16, "three": 81}))
	g.Expect(d.stdLog.logged).To(ConsistOf(`SELECT "aa", "bb" FROM "p"."table" WHERE "foo"=$1 [bar]` + "\n"))
}

func TestGetIntStringIndex_happy(t *testing.T) {
	g := NewGomegaWithT(t)

	e := stubExecer{rows: &StubRows{
		Rows: []StubRow{{int64(2), "16"}, {int64(3), "81"}},
	}}
	d := &StubDatabase{execer: e, stdLog: &stubLogger{}}
	tbl := StubTable{
		name: sqlapi.TableName{
			Prefix: "p.",
			Name:   "table",
		},
		dialect:  dialect.Postgres,
		database: d,
	}

	m, err := GetIntStringIndex(tbl, quote.AnsiQuoter, "aa", "bb", where.Eq("foo", "bar"))

	g.Expect(err).NotTo(HaveOccurred())
	g.Expect(m).To(Equal(map[int64]string{2: "16", 3: "81"}))
	g.Expect(d.stdLog.logged).To(ConsistOf(`SELECT "aa", "bb" FROM "p"."table" WHERE "foo"=$1 [bar]` + "\n"))
}
