package constraint_test

import (
	"context"
	"database/sql"
	"fmt"
	_ "github.com/go-sql-driver/mysql"
	_ "github.com/lib/pq"
	_ "github.com/mattn/go-sqlite3"
	. "github.com/onsi/gomega"
	"github.com/mercury-holidays/sqlapi"
	"github.com/mercury-holidays/sqlapi/constraint"
	"github.com/mercury-holidays/sqlapi/dialect"
	"github.com/mercury-holidays/sqlapi/schema"
	"github.com/mercury-holidays/sqlapi/types"
	"github.com/mercury-holidays/sqlapi/vanilla"
	"github.com/rickb777/where/quote"
	"log"
	"os"
	"strings"
	"testing"
)

var db *sql.DB
var di dialect.Dialect

func connect() {
	dbDriver, ok := os.LookupEnv("GO_DRIVER")
	if !ok {
		dbDriver = "sqlite3"
	}
	di = dialect.PickDialect(dbDriver)
	dsn, ok := os.LookupEnv("GO_DSN")
	if !ok {
		dsn = ":memory:"
	}
	conn, err := sql.Open(dbDriver, dsn)
	if err != nil {
		panic(err)
	}
	db = conn
}

func newDatabase() sqlapi.Database {
	connect()
	if db == nil {
		return nil
	}

	ex := sqlapi.WrapDB(db)
	d := sqlapi.NewDatabase(ex, di, nil, nil)
	if testing.Verbose() {
		lgr := log.New(os.Stdout, "", log.LstdFlags)
		d = sqlapi.NewDatabase(ex, di, lgr, nil)
	}
	return d
}

func cleanup() {
	if db != nil {
		db.Close()
		db = nil
	}
}

func TestCheckConstraint(t *testing.T) {
	g := NewGomegaWithT(t)
	d := newDatabase()
	defer cleanup()

	cc0 := constraint.CheckConstraint{
		Expression: "role < 3",
	}

	persons := vanilla.NewRecordTable("persons", d).WithPrefix("pfx_").WithConstraint(cc0)
	fkc := persons.Constraints()[0]
	s := fkc.ConstraintSql(quote.AnsiQuoter, persons.Name(), 0)
	g.Expect(s).To(Equal(`CONSTRAINT "pfx_persons_c0" CHECK (role < 3)`), s)
}

func TestForeignKeyConstraint_withParentColumn(t *testing.T) {
	g := NewGomegaWithT(t)
	d := newDatabase()
	defer cleanup()

	fkc0 := constraint.FkConstraint{
		ForeignKeyColumn: "addresspk",
		Parent:           constraint.Reference{TableName: "addresses", Column: "identity"},
		Update:           "restrict",
		Delete:           "cascade"}

	persons := vanilla.NewRecordTable("persons", d).WithPrefix("pfx_").WithConstraint(fkc0)
	fkc := persons.Constraints()[0]
	s := fkc.ConstraintSql(quote.AnsiQuoter, persons.Name(), 0)
	g.Expect(s).To(Equal(`CONSTRAINT "pfx_persons_c0" foreign key ("addresspk") references "pfx_addresses" ("identity") on update restrict on delete cascade`), s)
}

func TestForeignKeyConstraint_withoutParentColumn_withoutQuotes(t *testing.T) {
	g := NewGomegaWithT(t)
	d := newDatabase()
	defer cleanup()

	fkc0 := constraint.FkConstraint{
		ForeignKeyColumn: "addresspk",
		Parent:           constraint.Reference{TableName: "addresses", Column: ""},
		Update:           "restrict",
		Delete:           "cascade"}

	persons := vanilla.NewRecordTable("persons", d).WithPrefix("pfx_").WithConstraint(fkc0)
	fkc := persons.Constraints().FkConstraints()[0]
	s := fkc.ConstraintSql(quote.NoQuoter, persons.Name(), 0)
	g.Expect(s).To(Equal(`CONSTRAINT pfx_persons_c0 foreign key (addresspk) references pfx_addresses on update restrict on delete cascade`), s)
}

func TestIdsUsedAsForeignKeys(t *testing.T) {
	g := NewGomegaWithT(t)
	d := newDatabase()
	defer cleanup()

	fkc0 := constraint.FkConstraint{
		ForeignKeyColumn: "addressid",
		Parent:           constraint.Reference{TableName: "addresses", Column: "id"},
		Update:           "cascade",
		Delete:           "cascade"}

	persons := vanilla.NewRecordTable("persons", d).WithPrefix("pfx_").WithConstraint(fkc0)

	setupSql := strings.Replace(createTables, "¬", "`", -1)
	_, err := d.DB().ExecContext(context.Background(), setupSql)
	g.Expect(err).To(BeNil())

	aid1 := insertOne(g, d, address1)
	aid2 := insertOne(g, d, address2)
	aid3 := insertOne(g, d, address3)
	aid4 := insertOne(g, d, address4)

	insertOne(g, d, fmt.Sprintf(person1a, aid1))
	insertOne(g, d, fmt.Sprintf(person1b, aid1))
	insertOne(g, d, fmt.Sprintf(person2a, aid2))

	fkc := persons.Constraints().FkConstraints()[0]

	m1, err := fkc.RelationshipWith(persons.Name()).IdsUsedAsForeignKeys(persons)

	g.Expect(err).To(BeNil())
	g.Expect(m1).To(HaveLen(2))
	g.Expect(m1.Contains(aid1)).To(BeTrue())
	g.Expect(m1.Contains(aid2)).To(BeTrue())

	m2, err := fkc.RelationshipWith(persons.Name()).IdsUnusedAsForeignKeys(persons)

	g.Expect(err).To(BeNil())
	g.Expect(m2).To(HaveLen(2))
	g.Expect(m2.Contains(aid3)).To(BeTrue())
	g.Expect(m2.Contains(aid4)).To(BeTrue())
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

func insertOne(g *GomegaWithT, d sqlapi.Database, query string) int64 {
	fmt.Fprintf(os.Stderr, "%s\n", query)
	id, err := d.DB().InsertContext(context.Background(), query)
	g.Expect(err).To(BeNil())
	return id
}

//-------------------------------------------------------------------------------------------------

const createTables = `
CREATE TABLE IF NOT EXISTS pfx_addresses (
 ¬id¬        integer primary key autoincrement,
 ¬lines¬     text,
 ¬postcode¬  text
);

CREATE TABLE IF NOT EXISTS pfx_persons (
 ¬uid¬       integer primary key autoincrement,
 ¬name¬      text,
 ¬addressid¬ integer default null
);

DELETE FROM pfx_persons;
DELETE FROM pfx_addresses;
`

const address1 = `INSERT INTO pfx_addresses (lines, postcode) VALUES ('Laurel Cottage', 'FX1 1AA')`
const address2 = `INSERT INTO pfx_addresses (lines, postcode) VALUES ('2 Nutmeg Lane', 'FX1 2BB')`
const address3 = `INSERT INTO pfx_addresses (lines, postcode) VALUES ('Corner Shop', 'FX1 3CC')`
const address4 = `INSERT INTO pfx_addresses (lines, postcode) VALUES ('4 The Oaks', 'FX1 5EE')`

const person1a = `INSERT INTO pfx_persons (name, addressid) VALUES ('John Brown', %d)`
const person1b = `INSERT INTO pfx_persons (name, addressid) VALUES ('Mary Brown', %d)`

const person2a = `INSERT INTO pfx_persons (name, addressid) VALUES ('Anne Bollin', %d)`
