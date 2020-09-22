package sqlapi_test

import (
	"bytes"
	"context"
	"database/sql"
	"errors"
	"fmt"
	"io"
	"log"
	"os"
	"strings"
	"testing"

	_ "github.com/go-sql-driver/mysql"
	_ "github.com/jackc/pgx/stdlib"
	_ "github.com/lib/pq"
	_ "github.com/mattn/go-sqlite3"
	. "github.com/onsi/gomega"
	"github.com/rickb777/sqlapi"
	"github.com/rickb777/sqlapi/dialect"
	"github.com/rickb777/where/quote"
)

// Environment:
// GO_DRIVER  - the driver (sqlite3, mysql, postgres, pgx)
// GO_QUOTER  - the identifier quoter (ansi, mysql, none)
// GO_DSN     - the database DSN
// GO_VERBOSE - true for query logging

var gdb *sql.DB
var gdi dialect.Dialect

func TestLoggingOnOff(t *testing.T) {
	g := NewGomegaWithT(t)

	buf := &bytes.Buffer{}
	logger := log.New(buf, "X.", 0)
	tl := sqlapi.NewLogger(logger)

	tl.LogQuery("one")
	tl.TraceLogging(false)
	tl.LogQuery("two")
	tl.TraceLogging(true)
	tl.LogQuery("three")

	s := buf.String()
	g.Expect(s).To(Equal("X.one\nX.three\n"))
}

func TestLoggingError(t *testing.T) {
	g := NewGomegaWithT(t)

	buf := &bytes.Buffer{}
	logger := log.New(buf, "X.", 0)
	tl := sqlapi.NewLogger(logger)

	tl.LogError(fmt.Errorf("one"))
	tl.TraceLogging(false)
	tl.LogError(fmt.Errorf("two"))
	tl.TraceLogging(true)
	tl.LogError(fmt.Errorf("three"))
	tl.LogIfError(nil)
	tl.LogIfError(fmt.Errorf("four"))

	s := buf.String()
	g.Expect(s).To(Equal("X.Error: one\nX.Error: three\nX.Error: four\n"))
}

func TestListTables(t *testing.T) {
	g := NewGomegaWithT(t)

	d := newDatabase(t)

	list, err := sqlapi.ListTables(d, nil)
	g.Expect(err).NotTo(HaveOccurred())
	g.Expect(list.Filter(func(s string) bool {
		return strings.HasPrefix(s, "sql_")
	})).To(HaveLen(0))
	g.Expect(list.Filter(func(s string) bool {
		return strings.HasPrefix(s, "pg_")
	})).To(HaveLen(0))
}

func TestQueryRowContext(t *testing.T) {
	g := NewGomegaWithT(t)

	d := newDatabase(t)

	_, aid2, _, _ := insertFixtures(t, d)

	q := d.Dialect().ReplacePlaceholders("select xlines from pfx_addresses where id=?", nil)
	row := d.QueryRowContext(context.Background(), q, aid2)

	var xlines string
	err := row.Scan(&xlines)
	g.Expect(err).NotTo(HaveOccurred())
	g.Expect(xlines).To(Equal("2 Nutmeg Lane"))
}

func TestQueryContext(t *testing.T) {
	g := NewGomegaWithT(t)

	d := newDatabase(t)

	_, aid2, _, _ := insertFixtures(t, d)

	q := d.Dialect().ReplacePlaceholders("select xlines from pfx_addresses where id=?", nil)
	rows, err := d.QueryContext(context.Background(), q, aid2)
	g.Expect(err).NotTo(HaveOccurred())
	g.Expect(rows.Next()).To(BeTrue())

	var xlines string
	err = rows.Scan(&xlines)
	g.Expect(err).NotTo(HaveOccurred())
	g.Expect(xlines).To(Equal("2 Nutmeg Lane"))

	g.Expect(rows.Next()).NotTo(BeTrue())
}

func TestSingleConnQueryContext(t *testing.T) {
	g := NewGomegaWithT(t)

	d := newDatabase(t)

	_, aid2, _, _ := insertFixtures(t, d)

	q := d.Dialect().ReplacePlaceholders("select xlines from pfx_addresses where id=?", nil)
	e2 := d.SingleConn(nil, func(ex sqlapi.Execer) error {
		rows, err := ex.QueryContext(context.Background(), q, aid2)
		g.Expect(err).NotTo(HaveOccurred())
		g.Expect(rows.Next()).To(BeTrue())

		var xlines string
		err = rows.Scan(&xlines)
		g.Expect(err).NotTo(HaveOccurred())
		g.Expect(xlines).To(Equal("2 Nutmeg Lane"))

		g.Expect(rows.Next()).NotTo(BeTrue())
		return err
	})
	g.Expect(e2).NotTo(HaveOccurred())
}

func TestTransactCommitUsingInsert(t *testing.T) {
	g := NewGomegaWithT(t)

	d := newDatabase(t)

	ctx := context.Background()
	insertFixtures(t, d)

	q := d.Dialect().ReplacePlaceholders("INSERT INTO pfx_addresses (xlines, postcode) VALUES (?, ?)", nil)
	err := d.Transact(ctx, nil, func(tx sqlapi.SqlTx) error {
		for i := 1; i <= 10; i++ {
			_, e2 := tx.InsertContext(ctx, "id", q, fmt.Sprintf("%d Pantagon Vale", i), "FX1 5EE")
			if e2 != nil {
				return e2
			}
		}
		return nil
	})
	g.Expect(err).NotTo(HaveOccurred())

	row := d.QueryRowContext(ctx, "select count(1) from pfx_addresses")

	var count int
	err = row.Scan(&count)
	g.Expect(err).NotTo(HaveOccurred())
	g.Expect(count).To(Equal(14))
}

func TestTransactCommitUsingExec(t *testing.T) {
	g := NewGomegaWithT(t)

	d := newDatabase(t)

	ctx := context.Background()
	_, aid2, aid3, _ := insertFixtures(t, d)

	q := d.Dialect().ReplacePlaceholders("delete from pfx_addresses where id in(?,?)", nil)
	err := d.Transact(ctx, nil, func(tx sqlapi.SqlTx) error {
		_, e2 := tx.ExecContext(ctx, q, aid2, aid3)
		return e2
	})
	g.Expect(err).NotTo(HaveOccurred())

	row := d.QueryRowContext(ctx, "select count(1) from pfx_addresses")

	var count int
	err = row.Scan(&count)
	g.Expect(err).NotTo(HaveOccurred())
	g.Expect(count).To(Equal(2))
}

func TestTransactRollback(t *testing.T) {
	g := NewGomegaWithT(t)

	d := newDatabase(t)

	ctx := context.Background()
	_, aid2, aid3, _ := insertFixtures(t, d)

	q := d.Dialect().ReplacePlaceholders("delete from pfx_addresses where id in(?,?)", nil)
	err := d.Transact(ctx, nil, func(tx sqlapi.SqlTx) error {
		tx.ExecContext(ctx, q, aid2, aid3)
		return errors.New("Bang")
	})
	g.Expect(err).To(HaveOccurred())
	g.Expect(err.Error()).To(Equal("Bang"))

	row := d.QueryRowContext(ctx, "select count(1) from pfx_addresses")

	var count int
	err = row.Scan(&count)
	g.Expect(err).NotTo(HaveOccurred())
	g.Expect(count).To(Equal(4))
}

//-------------------------------------------------------------------------------------------------

func TestMain(m *testing.M) {
	gdb, gdi = connect()
	code := m.Run()
	cleanup(gdb)
	os.Exit(code)
}

//-------------------------------------------------------------------------------------------------

func connect() (*sql.DB, dialect.Dialect) {
	dbDriver, ok := os.LookupEnv("GO_DRIVER")
	if !ok {
		dbDriver = "sqlite3"
	}

	di := dialect.PickDialect(dbDriver)
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
			log.Fatalf("Warning: unrecognised quoter %q.\n", quoter)
		}
	}

	dsn, ok := os.LookupEnv("GO_DSN")
	if !ok {
		dsn = "file::memory:?mode=memory&cache=shared"
	}

	db, err := sql.Open(dbDriver, dsn)
	if err != nil {
		log.Fatalf("Error: Unable to connect to %s (%v); test is only partially complete.\n\n", dbDriver, err)
	}

	err = db.Ping()
	if err != nil {
		log.Fatalf("Error: Unable to ping %s (%v); test is only partially complete.\n\n", dbDriver, err)
	}

	fmt.Printf("Successfully connected to %s.\n", dbDriver)
	return db, di
}

func newDatabase(t *testing.T) sqlapi.SqlDB {
	var lgr *log.Logger
	goVerbose, ok := os.LookupEnv("GO_VERBOSE")
	if ok && strings.ToLower(goVerbose) == "true" {
		lgr = log.New(os.Stdout, "", log.LstdFlags)
	}

	return sqlapi.WrapDB(gdb, sqlapi.NewLogger(lgr), gdi)
}

func cleanup(db io.Closer) {
	if db != nil {
		db.Close()
		os.Remove("test.db")
	}
}
