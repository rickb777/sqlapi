package constraint_test

import (
	"context"
	"database/sql"
	"fmt"
	_ "github.com/go-sql-driver/mysql"
	_ "github.com/jackc/pgx/stdlib"
	_ "github.com/lib/pq"
	_ "github.com/mattn/go-sqlite3"
	. "github.com/onsi/gomega"
	"github.com/rickb777/sqlapi"
	"github.com/rickb777/sqlapi/constraint"
	"github.com/rickb777/sqlapi/dialect"
	"github.com/rickb777/sqlapi/schema"
	"github.com/rickb777/sqlapi/types"
	"github.com/rickb777/sqlapi/vanilla"
	"github.com/rickb777/where/quote"
	"io"
	"log"
	"os"
	"strings"
	"testing"
)

// Environment:
// GO_DRIVER  - the driver (sqlite3, mysql, postgres, pgx)
// GO_QUOTER  - the identifier quoter (ansi, mysql, none)
// GO_DSN     - the database DSN
// GO_VERBOSE - true for query logging

func skipIfNoPostgresDB(t *testing.T, di dialect.Dialect) {
	if (di.Index() == dialect.PostgresIndex || di.Index() == dialect.PgxIndex) && os.Getenv("PGHOST") == "" {
		t.Skip()
	}
}

func connect(t *testing.T) (*sql.DB, dialect.Dialect) {
	dbDriver, ok := os.LookupEnv("GO_DRIVER")
	if !ok {
		dbDriver = "sqlite3"
	}

	di := dialect.PickDialect(dbDriver) //.WithQuoter(dialect.NoQuoter)
	quoter, ok := os.LookupEnv("GO_QUOTER")
	if ok {
		switch strings.ToLower(quoter) {
		case "ansi":
			di = di.WithQuoter(quote.AnsiQuoter)
		case "mysql":
			di = di.WithQuoter(quote.MySqlQuoter)
		case "none":
			di = di.WithQuoter(quote.NoQuoter)
		default:
			t.Fatalf("Warning: unrecognised quoter %q.\n", quoter)
		}
	}

	skipIfNoPostgresDB(t, di)

	dsn, ok := os.LookupEnv("GO_DSN")
	if !ok {
		dsn = "file::memory:?mode=memory&cache=shared"
	}

	db, err := sql.Open(dbDriver, dsn)
	if err != nil {
		t.Fatalf("Error: Unable to connect to %s (%v); test is only partially complete.\n\n", dbDriver, err)
	}

	err = db.Ping()
	if err != nil {
		t.Fatalf("Error: Unable to ping %s (%v); test is only partially complete.\n\n", dbDriver, err)
	}

	fmt.Printf("Successfully connected to %s.\n", dbDriver)
	return db, di
}

func newDatabase(t *testing.T) sqlapi.Database {
	db, di := connect(t)

	var lgr *log.Logger
	goVerbose, ok := os.LookupEnv("GO_VERBOSE")
	if ok && strings.ToLower(goVerbose) == "true" {
		lgr = log.New(os.Stdout, "", log.LstdFlags)
	}

	return sqlapi.NewDatabase(sqlapi.WrapDB(db, di), di, lgr, nil)
}

func cleanup(db sqlapi.Execer) {
	if db != nil {
		if c, ok := db.(io.Closer); ok {
			c.Close()
		}
		os.Remove("test.db")
	}
}

func TestCheckConstraint(t *testing.T) {
	g := NewGomegaWithT(t)
	d := newDatabase(t)
	defer cleanup(d.DB())

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
	d := newDatabase(t)
	defer cleanup(d.DB())

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
	d := newDatabase(t)
	defer cleanup(d.DB())

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
	d := newDatabase(t)
	defer cleanup(d.DB())

	fkc0 := constraint.FkConstraint{
		ForeignKeyColumn: "addressid",
		Parent:           constraint.Reference{TableName: "addresses", Column: "id"},
		Update:           "cascade",
		Delete:           "cascade"}

	persons := vanilla.NewRecordTable("persons", d).WithPrefix("pfx_").WithConstraint(fkc0)

	for _, s := range createTablesSql(d.Dialect()) {
		_, err := d.DB().ExecContext(context.Background(), s)
		g.Expect(err).To(BeNil())
	}

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
	g.Expect(m1.ToSlice()).To(ConsistOf(aid1, aid2))

	m2, err := fkc.RelationshipWith(persons.Name()).IdsUnusedAsForeignKeys(persons)

	g.Expect(err).To(BeNil())
	g.Expect(m2.ToSlice()).To(ConsistOf(aid3, aid4))
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
	if !d.Dialect().HasLastInsertId() {
		query = query + " RETURNING id"
	}
	id, err := d.DB().InsertContext(context.Background(), query)
	g.Expect(err).To(BeNil())
	return id
}

//-------------------------------------------------------------------------------------------------

func createTablesSql(di dialect.Dialect) []string {
	switch di.Index() {
	case dialect.SqliteIndex:
		return createTablesSqlite
	case dialect.MysqlIndex:
		return createTablesMysql
	case dialect.PostgresIndex:
		return createTablesPostgresql
	}
	panic(di.String() + " unsupported")
}

var createTablesSqlite = []string{
	`DROP TABLE IF EXISTS pfx_addresses`,
	`DROP TABLE IF EXISTS pfx_persons`,

	`CREATE TABLE pfx_addresses (
	id        integer primary key autoincrement,
	xlines    text,
	postcode  text
	)`,

	`CREATE TABLE pfx_persons (
	id        integer primary key autoincrement,
	name      text,
	addressid integer default null
	)`,
}

var createTablesMysql = []string{
	`DROP TABLE IF EXISTS pfx_addresses`,
	`DROP TABLE IF EXISTS pfx_persons`,

	`CREATE TABLE pfx_addresses (
	id        int primary key auto_increment,
	xlines    text,
	postcode  text
	)`,

	`CREATE TABLE pfx_persons (
	id        int primary key auto_increment,
	name      text,
	addressid int default null
	)`,
}

var createTablesPostgresql = []string{
	`DROP TABLE IF EXISTS pfx_addresses`,
	`DROP TABLE IF EXISTS pfx_persons`,

	`CREATE TABLE pfx_addresses (
	id        serial primary key,
	xlines    text,
	postcode  text
	)`,

	`CREATE TABLE pfx_persons (
	id        serial primary key,
	name      text,
	addressid integer default null
	)`,
}

const address1 = `INSERT INTO pfx_addresses (xlines, postcode) VALUES ('Laurel Cottage', 'FX1 1AA')`
const address2 = `INSERT INTO pfx_addresses (xlines, postcode) VALUES ('2 Nutmeg Lane', 'FX1 2BB')`
const address3 = `INSERT INTO pfx_addresses (xlines, postcode) VALUES ('Corner Shop', 'FX1 3CC')`
const address4 = `INSERT INTO pfx_addresses (xlines, postcode) VALUES ('4 The Oaks', 'FX1 5EE')`

const person1a = `INSERT INTO pfx_persons (name, addressid) VALUES ('John Brown', %d)`
const person1b = `INSERT INTO pfx_persons (name, addressid) VALUES ('Mary Brown', %d)`
const person2a = `INSERT INTO pfx_persons (name, addressid) VALUES ('Anne Bollin', %d)`
