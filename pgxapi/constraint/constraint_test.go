package constraint_test

import (
	"context"
	"testing"

	"github.com/jackc/pgx/v5/tracelog"
	. "github.com/onsi/gomega"
	"github.com/rickb777/sqlapi/pgxapi"
	"github.com/rickb777/sqlapi/pgxapi/constraint"
	"github.com/rickb777/sqlapi/pgxapi/vanilla"
	"github.com/rickb777/sqlapi/schema"
	"github.com/rickb777/sqlapi/support/testenv"
	"github.com/rickb777/sqlapi/types"
	"github.com/rickb777/where/quote"
)

var gdb pgxapi.SqlDB

func TestPgxCheckConstraint(t *testing.T) {
	g := NewGomegaWithT(t)

	cc0 := constraint.CheckConstraint{
		Expression: "role < 3",
	}

	persons := vanilla.NewRecordTable("persons", gdb).WithPrefix("constraint_").WithConstraint(cc0)
	fkc := persons.Constraints()[0]
	s := fkc.ConstraintSql(quote.AnsiQuoter, persons.Name(), 0)
	g.Expect(s).To(Equal(`CONSTRAINT "constraint_persons_c0" CHECK (role < 3)`), s)
}

func TestPgxForeignKeyConstraint_withParentColumn(t *testing.T) {
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

func TestPgxForeignKeyConstraint_withoutParentColumn_withoutQuotes(t *testing.T) {
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

func TestPgxIdsUsedAsForeignKeys(t *testing.T) {
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
	g.Expect(m1.ToSlice()).To(ConsistOf(aid1, aid2))

	m2, err := fkc.RelationshipWith(persons.Name()).IdsUnusedAsForeignKeys(persons)

	g.Expect(err).To(BeNil())
	g.Expect(m2.ToSlice()).To(ConsistOf(aid3, aid4))
}

func TestPgxFkConstraintOfField(t *testing.T) {
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
	testenv.SetDefaultDbDriver("pgx")
	testenv.Shebang(m, func(lgr tracelog.Logger, logLevel tracelog.LogLevel, tries int) (err error) {
		gdb, err = pgxapi.ConnectEnv(context.Background(), lgr, logLevel, tries)
		return err
	})
}
