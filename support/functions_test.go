package support

import (
	"database/sql"
	"testing"

	"github.com/rickb777/sqlapi/support/test"

	"github.com/benmoss/matchers"
	. "github.com/onsi/gomega"
	"github.com/rickb777/sqlapi"
	"github.com/rickb777/sqlapi/dialect"
	"github.com/rickb777/sqlapi/require"
	"github.com/rickb777/where"
	"github.com/rickb777/where/quote"
)

func TestUpdateFieldsSQL(t *testing.T) {
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

func TestSliceSql(t *testing.T) {
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
		stdLog := &test.StubLogger{}
		lgr := sqlapi.NewLogger(stdLog)
		ex := &test.StubExecer{Di: c.dialect, Lgr: lgr}
		tbl := sqlapi.CoreTable{
			Nm: sqlapi.TableName{
				Prefix: "p.",
				Name:   "table",
			},
			Ex: ex,
		}
		wh := where.Eq("room", 101).And(where.Eq("fun", true))

		q, a := sliceSql(tbl, "foo", wh, where.OrderBy("xyz"))

		g.Expect(q).To(Equal(c.expected))
		g.Expect(a).To(matchers.DeepEqual([]interface{}{101, true}))
	}
}

func TestQuery_happy(t *testing.T) {
	g := NewGomegaWithT(t)

	stdLog := &test.StubLogger{Testing: t}
	lgr := sqlapi.NewLogger(stdLog)
	ex := &test.StubExecer{Di: dialect.Postgres, Lgr: lgr}
	tbl := sqlapi.CoreTable{
		Nm: sqlapi.TableName{
			Prefix: "p.",
			Name:   "table",
		},
		Ex: ex,
	}

	_, err := Query(tbl, "SELECT foo FROM p.table WHERE x=?", 123)

	g.Expect(err).NotTo(HaveOccurred())
	g.Expect(stdLog.Logged).To(ConsistOf(`info  SELECT foo FROM p.table WHERE x=$1 [$1=123]`))
}

func TestExec_happy(t *testing.T) {
	g := NewGomegaWithT(t)

	stdLog := &test.StubLogger{}
	lgr := sqlapi.NewLogger(stdLog)
	ex := &test.StubExecer{Di: dialect.Postgres, N: 2, Lgr: lgr}
	tbl := sqlapi.CoreTable{
		Nm: sqlapi.TableName{
			Prefix: "p.",
			Name:   "table",
		},
		Ex: ex,
	}

	_, err := Exec(tbl, require.Exactly(2), "DELETE FROM p.table WHERE x=?", 123)

	g.Expect(err).NotTo(HaveOccurred())
	g.Expect(stdLog.Logged).To(ConsistOf(`info  DELETE FROM p.table WHERE x=$1 [$1=123]`))
}

func TestUpdateFields(t *testing.T) {
	g := NewGomegaWithT(t)

	stdLog := &test.StubLogger{}
	lgr := sqlapi.NewLogger(stdLog)
	ex := &test.StubExecer{Di: dialect.Postgres, N: 2, Lgr: lgr}
	tbl := sqlapi.CoreTable{
		Nm: sqlapi.TableName{
			Prefix: "p.",
			Name:   "table",
		},
		Ex: ex,
	}

	_, err := UpdateFields(tbl, require.Exactly(2), where.Eq("foo", "bar"), sql.Named("c1", 1), sql.Named("c2", 2))

	g.Expect(err).NotTo(HaveOccurred())
	g.Expect(stdLog.Logged).To(ConsistOf(`info  UPDATE "p"."table" SET "c1"=$1, "c2"=$2 WHERE "foo"=$3 [$1=1, $2=2, $3=bar]`))
}

func TestDeleteByColumn(t *testing.T) {
	g := NewGomegaWithT(t)

	stdLog := &test.StubLogger{}
	lgr := sqlapi.NewLogger(stdLog)
	ex := &test.StubExecer{Di: dialect.Postgres, N: 2, Lgr: lgr}
	tbl := sqlapi.CoreTable{
		Nm: sqlapi.TableName{
			Prefix: "p.",
			Name:   "table",
		},
		Ex: ex,
	}

	_, err := DeleteByColumn(tbl, require.Exactly(2), "foo", 1, 2, 3, 4)

	g.Expect(err).NotTo(HaveOccurred())
	g.Expect(stdLog.Logged).To(ConsistOf(`info  DELETE FROM "p"."table" WHERE "foo" IN ($1,$2,$3,$4) [$1=1, $2=2, $3=3, $4=4]`))
}

func TestGetIntIntIndex_happy(t *testing.T) {
	g := NewGomegaWithT(t)

	stdLog := &test.StubLogger{}
	lgr := sqlapi.NewLogger(stdLog)
	ex := &test.StubExecer{Rows: &test.StubRows{
		Rows: []test.StubRow{{int64(2), int64(16)}, {int64(3), int64(81)}},
	}, Di: dialect.Postgres, Lgr: lgr}
	tbl := sqlapi.CoreTable{
		Nm: sqlapi.TableName{
			Prefix: "p.",
			Name:   "table",
		},
		Ex: ex,
	}

	m, err := GetIntIntIndex(tbl, quote.AnsiQuoter, "aa", "bb", where.Eq("foo", "bar"))

	g.Expect(err).NotTo(HaveOccurred())
	g.Expect(m).To(Equal(map[int64]int64{2: 16, 3: 81}))
	g.Expect(stdLog.Logged).To(ConsistOf(`info  SELECT "aa", "bb" FROM "p"."table" WHERE "foo"=$1 [$1=bar]`))
}

func TestGetStringIntIndex_happy(t *testing.T) {
	g := NewGomegaWithT(t)

	stdLog := &test.StubLogger{}
	lgr := sqlapi.NewLogger(stdLog)
	ex := &test.StubExecer{Rows: &test.StubRows{
		Rows: []test.StubRow{{"two", int64(16)}, {"three", int64(81)}},
	}, Di: dialect.Postgres, Lgr: lgr}
	tbl := sqlapi.CoreTable{
		Nm: sqlapi.TableName{
			Prefix: "p.",
			Name:   "table",
		},
		Ex: ex,
	}

	m, err := GetStringIntIndex(tbl, quote.AnsiQuoter, "aa", "bb", where.Eq("foo", "bar"))

	g.Expect(err).NotTo(HaveOccurred())
	g.Expect(m).To(Equal(map[string]int64{"two": 16, "three": 81}))
	g.Expect(stdLog.Logged).To(ConsistOf(`info  SELECT "aa", "bb" FROM "p"."table" WHERE "foo"=$1 [$1=bar]`))
}

func TestGetIntStringIndex_happy(t *testing.T) {
	g := NewGomegaWithT(t)

	stdLog := &test.StubLogger{}
	lgr := sqlapi.NewLogger(stdLog)
	ex := &test.StubExecer{Rows: &test.StubRows{
		Rows: []test.StubRow{{int64(2), "16"}, {int64(3), "81"}},
	}, Di: dialect.Postgres, Lgr: lgr}
	tbl := sqlapi.CoreTable{
		Nm: sqlapi.TableName{
			Prefix: "p.",
			Name:   "table",
		},
		Ex: ex,
	}

	m, err := GetIntStringIndex(tbl, quote.AnsiQuoter, "aa", "bb", where.Eq("foo", "bar"))

	g.Expect(err).NotTo(HaveOccurred())
	g.Expect(m).To(Equal(map[int64]string{2: "16", 3: "81"}))
	g.Expect(stdLog.Logged).To(ConsistOf(`info  SELECT "aa", "bb" FROM "p"."table" WHERE "foo"=$1 [$1=bar]`))
}
