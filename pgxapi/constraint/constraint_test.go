package constraint_test

import (
	"context"
	"testing"

	"github.com/jackc/pgx/v5/tracelog"
	"github.com/rickb777/expect"
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
	cc0 := constraint.CheckConstraint{
		Expression: "role < 3",
	}

	persons := vanilla.NewRecordTable("persons", gdb).WithPrefix("constraint_").WithConstraint(cc0)
	fkc := persons.Constraints()[0]
	s := fkc.ConstraintSql(quote.AnsiQuoter, persons.Name(), 0)
	expect.String(s).I(s).ToBe(t, `CONSTRAINT "constraint_persons_c0" CHECK (role < 3)`)
}

func TestPgxForeignKeyConstraint_withParentColumn(t *testing.T) {
	fkc0 := constraint.FkConstraint{
		ForeignKeyColumn: "addresspk",
		Parent:           constraint.Reference{TableName: "addresses", Column: "identity"},
		Update:           "restrict",
		Delete:           "cascade"}

	persons := vanilla.NewRecordTable("persons", gdb).WithPrefix("constraint_").WithConstraint(fkc0)
	fkc := persons.Constraints()[0]
	s := fkc.ConstraintSql(quote.AnsiQuoter, persons.Name(), 0)
	expect.String(s).I(s).ToBe(t, `CONSTRAINT "constraint_persons_c0" foreign key ("addresspk") references "constraint_addresses" ("identity") on update restrict on delete cascade`)
}

func TestPgxForeignKeyConstraint_withoutParentColumn_withoutQuotes(t *testing.T) {
	fkc0 := constraint.FkConstraint{
		ForeignKeyColumn: "addresspk",
		Parent:           constraint.Reference{TableName: "addresses", Column: ""},
		Update:           "restrict",
		Delete:           "cascade"}

	persons := vanilla.NewRecordTable("persons", gdb).WithPrefix("constraint_").WithConstraint(fkc0)
	fkc := persons.Constraints().FkConstraints()[0]
	s := fkc.ConstraintSql(quote.NoQuoter, persons.Name(), 0)
	expect.String(s).I(s).ToBe(t, `CONSTRAINT constraint_persons_c0 foreign key (addresspk) references constraint_addresses on update restrict on delete cascade`)
}

func TestPgxIdsUsedAsForeignKeys(t *testing.T) {
	aid1, aid2, aid3, aid4 := insertFixtures(t, gdb)

	fkc0 := constraint.FkConstraint{
		ForeignKeyColumn: "addressid",
		Parent:           constraint.Reference{TableName: "addresses", Column: "id"},
		Update:           "cascade",
		Delete:           "cascade"}

	persons := vanilla.NewRecordTable("persons", gdb).WithPrefix("constraint_").WithConstraint(fkc0)

	fkc := persons.Constraints().FkConstraints()[0]

	m1, err := fkc.RelationshipWith(persons.Name()).IdsUsedAsForeignKeys(persons)

	expect.Error(err).ToBeNil(t)
	expect.Slice(m1.Slice()).ToContainAll(t, aid1, aid2)

	m2, err := fkc.RelationshipWith(persons.Name()).IdsUnusedAsForeignKeys(persons)

	expect.Error(err).ToBeNil(t)
	expect.Slice(m2.Slice()).ToContainAll(t, aid3, aid4)
}

func TestPgxFkConstraintOfField(t *testing.T) {
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
	expect.Any(fkc).ToBe(t, constraint.FkConstraint{
		ForeignKeyColumn: "cat",
		Parent: constraint.Reference{
			TableName: "something",
			Column:    "pk",
		},
		Update: "restrict",
		Delete: "cascade",
	})
}

//-------------------------------------------------------------------------------------------------

func TestMain(m *testing.M) {
	testenv.SetDefaultDbDriver("pgx")
	testenv.Shebang(m, func(lgr tracelog.Logger, logLevel tracelog.LogLevel, tries int) (err error) {
		gdb, err = pgxapi.ConnectEnv(context.Background(), lgr, logLevel, tries)
		return err
	})
}
