package where_test

import (
	"fmt"
	"github.com/benmoss/matchers"
	. "github.com/onsi/gomega"
	"github.com/rickb777/sqlapi/where"
	"testing"
)

func TestBuildWhereClause_happyCases(t *testing.T) {
	g := NewGomegaWithT(t)

	nameEqFred := where.Eq("name", "Fred")
	nameEqJohn := where.Eq("name", "John")
	ageLt10 := where.Lt("age", 10)
	ageGt5 := where.Gt("age", 5)

	cases := []struct {
		wh        where.Expression
		expSQL    string
		expString string
		args      []interface{}
	}{
		{where.NoOp(), "", "", nil},

		{
			where.Condition{Column: "name", Predicate: " not nil", Args: nil},
			`WHERE "name" not nil`,
			`"name" not nil`,
			nil,
		},

		{
			where.Condition{Column: "p.name", Predicate: " not nil", Args: nil},
			`WHERE "p"."name" not nil`,
			`"p"."name" not nil`,
			nil,
		},

		{
			where.Null("name"),
			`WHERE "name" IS NULL`,
			`"name" IS NULL`,
			nil,
		},

		{
			where.NotNull("name"),
			`WHERE "name" IS NOT NULL`,
			`"name" IS NOT NULL`,
			nil,
		},

		{
			where.Condition{Column: "name", Predicate: " <>?", Args: []interface{}{"Boo"}},
			`WHERE "name" <>?`,
			`"name" <>'Boo'`,
			[]interface{}{"Boo"},
		},

		{
			nameEqFred,
			`WHERE "name"=?`,
			`"name"='Fred'`,
			[]interface{}{"Fred"},
		},

		{
			where.Like("name", "F%"),
			`WHERE "name" LIKE ?`,
			`"name" LIKE 'F%'`,
			[]interface{}{"F%"},
		},

		{
			where.NoOp().And(nameEqFred),
			`WHERE ("name"=?)`,
			`("name"='Fred')`,
			[]interface{}{"Fred"},
		},

		{
			nameEqFred.And(where.Gt("age", 10)),
			`WHERE ("name"=?) AND ("age">?)`,
			`("name"='Fred') AND ("age">10)`,
			[]interface{}{"Fred", 10},
		},

		{
			nameEqFred.Or(where.Gt("age", 10)),
			`WHERE ("name"=?) OR ("age">?)`,
			`("name"='Fred') OR ("age">10)`,
			[]interface{}{"Fred", 10},
		},

		{
			nameEqFred.And(ageGt5).And(where.Gt("weight", 15)),
			`WHERE ("name"=?) AND ("age">?) AND ("weight">?)`,
			`("name"='Fred') AND ("age">5) AND ("weight">15)`,
			[]interface{}{"Fred", 5, 15},
		},

		{
			nameEqFred.Or(ageGt5).Or(where.Gt("weight", 15)),
			`WHERE ("name"=?) OR ("age">?) OR ("weight">?)`,
			`("name"='Fred') OR ("age">5) OR ("weight">15)`,
			[]interface{}{"Fred", 5, 15},
		},

		{
			where.Between("age", 12, 18).Or(where.Gt("weight", 45)),
			`WHERE ("age" BETWEEN ? AND ?) OR ("weight">?)`,
			`("age" BETWEEN 12 AND 18) OR ("weight">45)`,
			[]interface{}{12, 18, 45},
		},

		{
			where.GtEq("age", 10),
			`WHERE "age">=?`,
			`"age">=10`,
			[]interface{}{10},
		},

		{
			where.LtEq("age", 10),
			`WHERE "age"<=?`,
			`"age"<=10`,
			[]interface{}{10},
		},

		{
			where.NotEq("age", 10),
			`WHERE "age"<>?`,
			`"age"<>10`,
			[]interface{}{10},
		},

		{
			where.In("age", 10, 12, 14),
			`WHERE "age" IN (?,?,?)`,
			`"age" IN (10,12,14)`,
			[]interface{}{10, 12, 14},
		},

		{
			where.In("age", []int{10, 12, 14}),
			`WHERE "age" IN (?,?,?)`,
			`"age" IN (10,12,14)`,
			[]interface{}{10, 12, 14},
		},

		{ // 'In' without any vararg parameters
			where.In("age"),
			``,
			``,
			nil,
		},

		{ // 'In' with mixed value and nil vararg parameters
			where.In("age", 1, nil, 2, nil),
			`WHERE ("age" IN (?,?)) OR ("age" IS NULL)`,
			`("age" IN (1,2)) OR ("age" IS NULL)`,
			[]interface{}{1, 2},
		},

		{ // 'In' with only a nil vararg parameter
			where.In("age", nil),
			`WHERE ("age" IS NULL)`,
			`("age" IS NULL)`,
			nil,
		},

		{
			where.Not(nameEqFred),
			`WHERE NOT ("name"=?)`,
			`NOT ("name"='Fred')`,
			[]interface{}{"Fred"},
		},

		{
			where.Not(nameEqFred.And(ageLt10)),
			`WHERE NOT (("name"=?) AND ("age"<?))`,
			`NOT (("name"='Fred') AND ("age"<10))`,
			[]interface{}{"Fred", 10},
		},

		{
			where.Not(nameEqFred.Or(ageLt10)),
			`WHERE NOT (("name"=?) OR ("age"<?))`,
			`NOT (("name"='Fred') OR ("age"<10))`,
			[]interface{}{"Fred", 10},
		},

		{
			where.Not(nameEqFred).And(ageLt10),
			`WHERE (NOT ("name"=?)) AND ("age"<?)`,
			`(NOT ("name"='Fred')) AND ("age"<10)`,
			[]interface{}{"Fred", 10},
		},

		{
			where.Not(nameEqFred).Or(ageLt10),
			`WHERE (NOT ("name"=?)) OR ("age"<?)`,
			`(NOT ("name"='Fred')) OR ("age"<10)`,
			[]interface{}{"Fred", 10},
		},

		{
			where.And(nameEqFred, ageLt10),
			`WHERE ("name"=?) AND ("age"<?)`,
			`("name"='Fred') AND ("age"<10)`,
			[]interface{}{"Fred", 10},
		},

		{
			where.And(nameEqFred).And(where.And(ageLt10)),
			`WHERE ("name"=?) AND ("age"<?)`,
			`("name"='Fred') AND ("age"<10)`,
			[]interface{}{"Fred", 10},
		},

		{
			where.Or(nameEqFred, ageLt10),
			`WHERE ("name"=?) OR ("age"<?)`,
			`("name"='Fred') OR ("age"<10)`,
			[]interface{}{"Fred", 10},
		},

		{
			where.And(nameEqFred.Or(nameEqJohn), ageLt10),
			`WHERE (("name"=?) OR ("name"=?)) AND ("age"<?)`,
			`(("name"='Fred') OR ("name"='John')) AND ("age"<10)`,
			[]interface{}{"Fred", "John", 10},
		},

		{
			where.Or(nameEqFred, ageLt10.And(ageGt5)),
			`WHERE ("name"=?) OR (("age"<?) AND ("age">?))`,
			`("name"='Fred') OR (("age"<10) AND ("age">5))`,
			[]interface{}{"Fred", 10, 5},
		},

		{
			where.Or(nameEqFred, nameEqJohn).And(ageGt5),
			`WHERE (("name"=?) OR ("name"=?)) AND ("age">?)`,
			`(("name"='Fred') OR ("name"='John')) AND ("age">5)`,
			[]interface{}{"Fred", "John", 5},
		},

		{
			where.Or(nameEqFred, nameEqJohn, where.And(ageGt5)),
			`WHERE ("name"=?) OR ("name"=?) OR (("age">?))`,
			`("name"='Fred') OR ("name"='John') OR (("age">5))`,
			[]interface{}{"Fred", "John", 5},
		},

		{
			where.Or().Or(where.NoOp()).And(where.NoOp()),
			"",
			"",
			nil,
		},

		{
			where.And(where.Or(where.NoOp())),
			"",
			"",
			nil,
		},
	}

	for i, c := range cases {
		info := fmt.Sprintf("%d: %s", i, c.expSQL)

		sql, args := where.Where(c.wh)

		g.Expect(sql).To(Equal(c.expSQL), info)
		g.Expect(args).To(matchers.DeepEqual(c.args), info)

		s := c.wh.String()

		g.Expect(s).To(Equal(c.expString), info)
	}
}

func TestQueryConstraint(t *testing.T) {
	g := NewGomegaWithT(t)

	cases := []struct {
		qc          where.QueryConstraint
		expPostgres string
	}{
		{nil, ""},
		{where.Literal("order by foo"), "order by foo"},
		{where.OrderBy("foo").Asc(), `ORDER BY "foo"`},
		{where.OrderBy("foo").Desc(), `ORDER BY "foo" DESC`},
		{where.OrderBy("foo").OrderBy("bar"), `ORDER BY "bar"`},
		{where.Limit(0), ""},
		{where.Limit(10), "LIMIT 10"},
		{where.Offset(20), "OFFSET 20"},
		{where.Limit(5).OrderBy("foo", "bar"), `ORDER BY "foo","bar" LIMIT 5`},
		{where.OrderBy("foo").Desc().Limit(10).Offset(20), `ORDER BY "foo" DESC LIMIT 10 OFFSET 20`},
	}

	for _, c := range cases {
		var sql string

		if c.qc != nil {
			sql = where.BuildQueryConstraint(c.qc)
		}

		g.Expect(sql).To(Equal(c.expPostgres))
	}
}
