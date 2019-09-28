package sqlapi_test

import (
	"bytes"
	"database/sql"
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

func TestLoggingOnOff(t *testing.T) {
	g := NewGomegaWithT(t)

	buf := &bytes.Buffer{}
	logger := log.New(buf, "X.", 0)

	db := sqlapi.NewDatabase(nil, dialect.Sqlite, logger, nil)
	db.Logger().LogQuery("one")
	db.Logger().TraceLogging(false)
	db.Logger().LogQuery("two")
	db.Logger().TraceLogging(true)
	db.Logger().LogQuery("three")

	s := buf.String()
	g.Expect(s).To(Equal("X.one\nX.three\n"))
}

func TestLoggingError(t *testing.T) {
	g := NewGomegaWithT(t)

	buf := &bytes.Buffer{}
	logger := log.New(buf, "X.", 0)

	db := sqlapi.NewDatabase(nil, dialect.Sqlite, logger, nil)
	db.Logger().LogError(fmt.Errorf("one"))
	db.Logger().TraceLogging(false)
	db.Logger().LogError(fmt.Errorf("two"))
	db.Logger().TraceLogging(true)
	db.Logger().LogError(fmt.Errorf("three"))
	db.Logger().LogIfError(nil)
	db.Logger().LogIfError(fmt.Errorf("four"))

	s := buf.String()
	g.Expect(s).To(Equal("X.Error: one\nX.Error: three\nX.Error: four\n"))
}

func TestListTables(t *testing.T) {
	g := NewGomegaWithT(t)

	db := newDatabase(t)
	defer cleanup(db.DB())

	list, err := db.ListTables(nil)
	g.Expect(err).NotTo(HaveOccurred())
	g.Expect(list.Filter(func(s string) bool {
		return strings.HasPrefix(s, "sql_")
	})).To(HaveLen(0))
	g.Expect(list.Filter(func(s string) bool {
		return strings.HasPrefix(s, "pg_")
	})).To(HaveLen(0))
}

func TestQueryContext(t *testing.T) {
	//g := NewGomegaWithT(t)
	//
	//db := newDatabase(t)
	//defer cleanup(db.DB())
	//
	//	list, err := db.DB().QueryContext(context.Background(), "")
	//	g.Expect(err).NotTo(HaveOccurred())
	//	g.Expect(list.Filter(func(s string) bool {
	//		return strings.HasPrefix(s, "sql_")
	//	})).To(HaveLen(0))
	//	g.Expect(list.Filter(func(s string) bool {
	//		return strings.HasPrefix(s, "pg_")
	//	})).To(HaveLen(0))
}
