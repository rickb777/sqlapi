package constraint_test

import (
	"context"
	"testing"

	_ "github.com/go-sql-driver/mysql"
	_ "github.com/jackc/pgx/v5/stdlib"
	"github.com/jackc/pgx/v5/tracelog"
	_ "github.com/lib/pq"
	_ "github.com/mattn/go-sqlite3"
	. "github.com/onsi/gomega"
	"github.com/rickb777/sqlapi"
	"github.com/rickb777/sqlapi/constraint"
	"github.com/rickb777/sqlapi/schema"
	"github.com/rickb777/sqlapi/support/testenv"
	"github.com/rickb777/sqlapi/types"
	"github.com/rickb777/sqlapi/vanilla"
	"github.com/rickb777/where/quote"
)

var gdb sqlapi.SqlDB

func TestCheckConstraint(t *testing.T) {
	g := NewGomegaWithT(t)

	cc0 := constraint.CheckConstraint{
		Expression: "role < 3",
	}

	persons := vanilla.NewRecordTable("persons", gdb).WithPrefix("constraint_").WithConstraint(cc0)
	fkc := persons.Constraints()[0]
	s := fkc.ConstraintSql(quote.AnsiQuoter, persons.Name(), 0)
	g.Expect(s).To(Equal(`CONSTRAINT "constraint_persons_c0" CHECK (role < 3)`), s)
}

func TestForeignKeyConstraint_withParentColumn(t *testing.T) {
	g := NewGomegaWithT(t)

	fkc0 := constraint.FkConstraint{
		ForeignKeyColumn: "addresspk",
		Parent:           constraint.Reference{TableName: "addresses", Column: "identity"},
		Update:           "restrict",
		Delete:           "cascade"}

	persons := vanilla.NewRecordTable("persons", gdb).WithPrefix("constraint_").WithConstraint(fkc0)
	fkc := persons.Constraints()[0]
	s := fkc.ConstraintSql(quote.AnsiQuoter, persons.Name(), 0)
	g.Expect(s).To(Equal(`CONSTRAINT "constraint_persons_c0" foreign key ("addresspk") references "constraint_addresses" ("identity") on update restrict on delete cascade`), s)
}

func TestForeignKeyConstraint_withoutParentColumn_withoutQuotes(t *testing.T) {
	g := NewGomegaWithT(t)

	fkc0 := constraint.FkConstraint{
		ForeignKeyColumn: "addresspk",
		Parent:           constraint.Reference{TableName: "addresses", Column: ""},
		Update:           "restrict",
		Delete:           "cascade"}

	persons := vanilla.NewRecordTable("persons", gdb).WithPrefix("constraint_").WithConstraint(fkc0)
	fkc := persons.Constraints().FkConstraints()[0]
	s := fkc.ConstraintSql(quote.NoQuoter, persons.Name(), 0)
	g.Expect(s).To(Equal(`CONSTRAINT constraint_persons_c0 foreign key (addresspk) references constraint_addresses on update restrict on delete cascade`), s)
}

func TestIdsUsedAsForeignKeys(t *testing.T) {
	g := NewGomegaWithT(t)

	aid1, aid2, aid3, aid4 := insertFixtures(t, gdb)

	fkc0 := constraint.FkConstraint{
		ForeignKeyColumn: "addressid",
		Parent:           constraint.Reference{TableName: "addresses", Column: "id"},
		Update:           "cascade",
		Delete:           "cascade"}

	persons := vanilla.NewRecordTable("persons", gdb).WithPrefix("constraint_").WithConstraint(fkc0)

	fkc := persons.Constraints().FkConstraints()[0]

	m1, err := fkc.RelationshipWith(persons.Name()).IdsUsedAsForeignKeys(persons)

	g.Expect(err).To(BeNil())
	g.Expect(m1.Slice()).To(ConsistOf(aid1, aid2))

	m2, err := fkc.RelationshipWith(persons.Name()).IdsUnusedAsForeignKeys(persons)

	g.Expect(err).To(BeNil())
	g.Expect(m2.Slice()).To(ConsistOf(aid3, aid4))
}

func TestFkConstraintOfField(t *testing.T) {
	g := NewGomegaWithT(t)

	i64 := schema.Type{Name: "int64", Base: types.Int64}
	field := &schema.Field{
		Node:    schema.Node{Name: "Cat", Type: i64},
		SqlName: "cat",
		Tags: &types.Tag{
			ForeignKey: "something.pk",
			OnUpdate:   "restrict",
			OnDelete:   "cascade",
		},
	}

	fkc := constraint.FkConstraintOfField(field)
	g.Expect(fkc).To(Equal(constraint.FkConstraint{
		ForeignKeyColumn: "cat",
		Parent: constraint.Reference{
			TableName: "something",
			Column:    "pk",
		},
		Update: "restrict",
		Delete: "cascade",
	}))
}

//-------------------------------------------------------------------------------------------------

func TestMain(m *testing.M) {
	testenv.SetDefaultDbDriver("sqlite3")
	testenv.Shebang(m, func(lgr tracelog.Logger, logLevel tracelog.LogLevel, tries int) (err error) {
		gdb, err = sqlapi.ConnectEnv(context.Background(), lgr, logLevel, tries)
		return err
	})
}
