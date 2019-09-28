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
	"sync"
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

var lock = sync.Mutex{}

func connect(t *testing.T) (*sql.DB, dialect.Dialect) {
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
			t.Fatalf("Warning: unrecognised quoter %q.\n", quoter)
		}
	}

	dsn, ok := os.LookupEnv("GO_DSN")
	if !ok {
		dsn = "file::memory:?mode=memory&cache=shared"
	}

	db, err := sql.Open(dbDriver, dsn)
	if err != nil {
		t.Fatalf("Error: Unable to connect to %s (%v); test is only partially complete.\n\n", dbDriver, err)
	}

	lock.Lock()

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
		lock.Unlock()
		os.Remove("test.db")
	}
}

func TestLoggingOnOff(t *testing.T) {
	g := NewGomegaWithT(t)

	buf := &bytes.Buffer{}
	logger := log.New(buf, "X.", 0)

	d := sqlapi.NewDatabase(nil, dialect.Sqlite, logger, nil)
	lgr := d.Logger()
	lgr.LogQuery("one")
	lgr.TraceLogging(false)
	lgr.LogQuery("two")
	lgr.TraceLogging(true)
	lgr.LogQuery("three")

	s := buf.String()
	g.Expect(s).To(Equal("X.one\nX.three\n"))
}

func TestLoggingError(t *testing.T) {
	g := NewGomegaWithT(t)

	buf := &bytes.Buffer{}
	logger := log.New(buf, "X.", 0)

	d := sqlapi.NewDatabase(nil, dialect.Sqlite, logger, nil)
	lgr := d.Logger()
	lgr.LogError(fmt.Errorf("one"))
	lgr.TraceLogging(false)
	lgr.LogError(fmt.Errorf("two"))
	lgr.TraceLogging(true)
	lgr.LogError(fmt.Errorf("three"))
	lgr.LogIfError(nil)
	lgr.LogIfError(fmt.Errorf("four"))

	s := buf.String()
	g.Expect(s).To(Equal("X.Error: one\nX.Error: three\nX.Error: four\n"))
}

func TestListTables(t *testing.T) {
	g := NewGomegaWithT(t)

	d := newDatabase(t)
	defer cleanup(d.DB())

	list, err := d.ListTables(nil)
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
	defer cleanup(d.DB())

	_, aid2, _, _ := insertFixtures(t, d)

	q := d.Dialect().ReplacePlaceholders("select xlines from pfx_addresses where id=?", nil)
	row := d.DB().QueryRowContext(context.Background(), q, aid2)

	var xlines string
	err := row.Scan(&xlines)
	g.Expect(err).NotTo(HaveOccurred())
	g.Expect(xlines).To(Equal("2 Nutmeg Lane"))
}

func TestQueryContext(t *testing.T) {
	g := NewGomegaWithT(t)

	d := newDatabase(t)
	defer cleanup(d.DB())

	_, aid2, _, _ := insertFixtures(t, d)

	q := d.Dialect().ReplacePlaceholders("select xlines from pfx_addresses where id=?", nil)
	rows, err := d.DB().QueryContext(context.Background(), q, aid2)
	g.Expect(err).NotTo(HaveOccurred())
	g.Expect(rows.Next()).To(BeTrue())

	var xlines string
	err = rows.Scan(&xlines)
	g.Expect(err).NotTo(HaveOccurred())
	g.Expect(xlines).To(Equal("2 Nutmeg Lane"))

	g.Expect(rows.Next()).NotTo(BeTrue())
}

func TestTransactCommit(t *testing.T) {
	g := NewGomegaWithT(t)

	d := newDatabase(t)
	defer cleanup(d.DB())

	ctx := context.Background()
	_, aid2, aid3, _ := insertFixtures(t, d)

	err := d.DB().Transact(ctx, nil, func(tx sqlapi.SqlTx) error {
		q := d.Dialect().ReplacePlaceholders("delete from pfx_addresses where id in(?,?)", nil)
		_, e2 := tx.ExecContext(ctx, q, aid2, aid3)
		return e2
	})
	g.Expect(err).NotTo(HaveOccurred())

	row := d.DB().QueryRowContext(ctx, "select count(1) from pfx_addresses")

	var count int
	err = row.Scan(&count)
	g.Expect(err).NotTo(HaveOccurred())
	g.Expect(count).To(Equal(2))
}

func TestTransactRollback(t *testing.T) {
	g := NewGomegaWithT(t)

	d := newDatabase(t)
	defer cleanup(d.DB())

	ctx := context.Background()
	_, aid2, aid3, _ := insertFixtures(t, d)

	err := d.DB().Transact(ctx, nil, func(tx sqlapi.SqlTx) error {
		q := d.Dialect().ReplacePlaceholders("delete from pfx_addresses where id in(?,?)", nil)
		tx.ExecContext(ctx, q, aid2, aid3)
		return errors.New("Bang")
	})
	g.Expect(err.Error()).To(Equal("Bang"))

	row := d.DB().QueryRowContext(ctx, "select count(1) from pfx_addresses")

	var count int
	err = row.Scan(&count)
	g.Expect(err).NotTo(HaveOccurred())
	g.Expect(count).To(Equal(4))
}
