package support

import (
	"database/sql"
	"github.com/benmoss/matchers"
	. "github.com/onsi/gomega"
	"github.com/rickb777/sqlapi/dialect"
	"github.com/rickb777/sqlapi/pgxapi"
	"github.com/rickb777/sqlapi/require"
	"github.com/rickb777/where"
	"github.com/rickb777/where/quote"
	"testing"
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
		stdLog := &stubLogger{}
		tbl := StubTable{
			name: pgxapi.TableName{
				Prefix: "p.",
				Name:   "table",
			},
			dialect: c.dialect,
			execer:  StubExecer{},
			logger:  pgxapi.NewLogger(stdLog),
		}
		wh := where.Eq("room", 101).And(where.Eq("fun", true))

		q, a := sliceSql(tbl, "foo", wh, where.OrderBy("xyz"))

		g.Expect(q).To(Equal(c.expected))
		g.Expect(a).To(matchers.DeepEqual([]interface{}{101, true}))
	}
}

func TestQuery_happy(t *testing.T) {
	g := NewGomegaWithT(t)

	stdLog := &stubLogger{}
	logger := pgxapi.NewLogger(stdLog)
	tbl := StubTable{
		name: pgxapi.TableName{
			Prefix: "p.",
			Name:   "table",
		},
		dialect: dialect.Postgres,
		execer:  StubExecer{Lgr: logger},
		logger:  logger,
	}

	_, err := Query(tbl, "SELECT foo FROM p.table WHERE x=?", 123)

	g.Expect(err).NotTo(HaveOccurred())
	g.Expect(stdLog.logged).To(ConsistOf(`SELECT foo FROM p.table WHERE x=$1 map[0:123]`))
}

func TestExec_happy(t *testing.T) {
	g := NewGomegaWithT(t)

	stdLog := &stubLogger{}
	logger := pgxapi.NewLogger(stdLog)
	e := StubExecer{StubResult: 2, Lgr: logger}
	tbl := StubTable{
		name: pgxapi.TableName{
			Prefix: "p.",
			Name:   "table",
		},
		dialect: dialect.Postgres,
		execer:  e,
		logger:  logger,
	}

	_, err := Exec(tbl, require.Exactly(2), "DELETE FROM p.table WHERE x=?", 123)

	g.Expect(err).NotTo(HaveOccurred())
	g.Expect(stdLog.logged).To(ConsistOf(`DELETE FROM p.table WHERE x=$1 map[0:123]`))
}

func TestUpdateFields(t *testing.T) {
	g := NewGomegaWithT(t)

	stdLog := &stubLogger{}
	logger := pgxapi.NewLogger(stdLog)
	e := StubExecer{StubResult: 2, Lgr: logger}
	tbl := StubTable{
		name: pgxapi.TableName{
			Prefix: "p.",
			Name:   "table",
		},
		dialect: dialect.Postgres,
		execer:  e,
		logger:  logger,
	}

	_, err := UpdateFields(tbl, require.Exactly(2), where.Eq("foo", "bar"), sql.Named("c1", 1), sql.Named("c2", 2))

	g.Expect(err).NotTo(HaveOccurred())
	g.Expect(stdLog.logged).To(ConsistOf(`UPDATE "p"."table" SET "c1"=$1, "c2"=$2 WHERE "foo"=$3 map[0:1 1:2 2:bar]`))
}

func TestDeleteByColumn(t *testing.T) {
	g := NewGomegaWithT(t)

	stdLog := &stubLogger{}
	logger := pgxapi.NewLogger(stdLog)
	e := StubExecer{StubResult: 2, Lgr: logger}
	tbl := StubTable{
		name: pgxapi.TableName{
			Prefix: "p.",
			Name:   "table",
		},
		dialect: dialect.Postgres,
		execer:  e,
		logger:  logger,
	}

	_, err := DeleteByColumn(tbl, require.Exactly(2), "foo", 1, 2, 3, 4)

	g.Expect(err).NotTo(HaveOccurred())
	g.Expect(stdLog.logged).To(ConsistOf(`DELETE FROM "p"."table" WHERE "foo" IN ($1,$2,$3,$4) map[0:1 1:2 2:3 3:4]`))
}

func TestGetIntIntIndex_happy(t *testing.T) {
	g := NewGomegaWithT(t)

	stdLog := &stubLogger{}
	logger := pgxapi.NewLogger(stdLog)
	e := StubExecer{Rows: &StubRows{
		Rows: []StubRow{{int64(2), int64(16)}, {int64(3), int64(81)}},
	}, Lgr: logger}
	tbl := StubTable{
		name: pgxapi.TableName{
			Prefix: "p.",
			Name:   "table",
		},
		dialect: dialect.Postgres,
		execer:  e,
		logger:  logger,
	}

	m, err := GetIntIntIndex(tbl, quote.AnsiQuoter, "aa", "bb", where.Eq("foo", "bar"))

	g.Expect(err).NotTo(HaveOccurred())
	g.Expect(m).To(Equal(map[int64]int64{2: 16, 3: 81}))
	g.Expect(stdLog.logged).To(ConsistOf(`SELECT "aa", "bb" FROM "p"."table" WHERE "foo"=? map[0:bar]`))
}

func TestGetStringIntIndex_happy(t *testing.T) {
	g := NewGomegaWithT(t)

	stdLog := &stubLogger{}
	logger := pgxapi.NewLogger(stdLog)
	e := StubExecer{Rows: &StubRows{
		Rows: []StubRow{{"two", int64(16)}, {"three", int64(81)}},
	}, Lgr: logger}
	tbl := StubTable{
		name: pgxapi.TableName{
			Prefix: "p.",
			Name:   "table",
		},
		dialect: dialect.Postgres,
		execer:  e,
		logger:  logger,
	}

	m, err := GetStringIntIndex(tbl, quote.AnsiQuoter, "aa", "bb", where.Eq("foo", "bar"))

	g.Expect(err).NotTo(HaveOccurred())
	g.Expect(m).To(Equal(map[string]int64{"two": 16, "three": 81}))
	g.Expect(stdLog.logged).To(ConsistOf(`SELECT "aa", "bb" FROM "p"."table" WHERE "foo"=? map[0:bar]`))
}

func TestGetIntStringIndex_happy(t *testing.T) {
	g := NewGomegaWithT(t)

	stdLog := &stubLogger{}
	logger := pgxapi.NewLogger(stdLog)
	e := StubExecer{Rows: &StubRows{
		Rows: []StubRow{{int64(2), "16"}, {int64(3), "81"}},
	}, Lgr: logger}
	tbl := StubTable{
		name: pgxapi.TableName{
			Prefix: "p.",
			Name:   "table",
		},
		dialect: dialect.Postgres,
		execer:  e,
		logger:  logger,
	}

	m, err := GetIntStringIndex(tbl, quote.AnsiQuoter, "aa", "bb", where.Eq("foo", "bar"))

	g.Expect(err).NotTo(HaveOccurred())
	g.Expect(m).To(Equal(map[int64]string{2: "16", 3: "81"}))
	g.Expect(stdLog.logged).To(ConsistOf(`SELECT "aa", "bb" FROM "p"."table" WHERE "foo"=? map[0:bar]`))
}
