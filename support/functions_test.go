package support

import (
	"database/sql"
	"github.com/benmoss/matchers"
	. "github.com/onsi/gomega"
	"github.com/rickb777/sqlapi/dialect"
	"github.com/rickb777/sqlapi/where"
	"testing"
)

func TestUpdateFieldsSQL(t *testing.T) {
	g := NewGomegaWithT(t)

	cases := []struct {
		dialect  dialect.Dialect
		expected string
	}{
		{
			dialect:  dialect.Mysql,
			expected: `UPDATE "foo" SET "col1"=?, "col2"=? WHERE ("room"=?) AND ("fun"=?)`,
		},
		{
			dialect:  dialect.Postgres,
			expected: `UPDATE "foo" SET "col1"=?, "col2"=? WHERE ("room"=?) AND ("fun"=?)`,
		},
	}

	for _, c := range cases {
		f1 := sql.Named("col1", 111)
		f2 := sql.Named("col2", 222)
		wh := where.Eq("room", 101).And(where.Eq("fun", true))

		q, a := updateFieldsSQL("foo", c.dialect, wh, f1, f2)

		g.Expect(q).To(Equal(c.expected), c.dialect.String())
		g.Expect(a).To(matchers.DeepEqual([]interface{}{111, 222, 101, true}), c.dialect.String())
	}
}
