package sqlapi

import (
	"bytes"
	"context"
	"errors"
	"flag"
	"fmt"
	"log"
	"net"
	"os"
	"strings"
	"testing"

	_ "github.com/go-sql-driver/mysql"
	"github.com/jackc/pgx/v4"
	"github.com/jackc/pgx/v4/log/testingadapter"
	_ "github.com/jackc/pgx/v4/stdlib"
	_ "github.com/lib/pq"
	_ "github.com/mattn/go-sqlite3"
	. "github.com/onsi/gomega"
	"github.com/rickb777/sqlapi/pgxapi/logadapter"
	"github.com/rickb777/sqlapi/support/testenv"
)

var gdb SqlDB

func TestLoggingOnOff(t *testing.T) {
	g := NewGomegaWithT(t)

	ctx := context.Background()
	buf := &bytes.Buffer{}
	logger := logadapter.NewLogger(log.New(buf, "X.", 0))
	tl := NewLogger(logger)

	tl.LogQuery(ctx, "one")
	tl.Log(ctx, pgx.LogLevelInfo, "two", nil)
	tl.TraceLogging(false)
	tl.LogQuery(ctx, "three")
	tl.TraceLogging(true)
	tl.LogQuery(ctx, "four")

	s := buf.String()
	g.Expect(s).To(Equal("X.one []\nX.two []\nX.four []\n"))
}

func TestLoggingError(t *testing.T) {
	g := NewGomegaWithT(t)

	ctx := context.Background()
	buf := &bytes.Buffer{}
	logger := logadapter.NewLogger(log.New(buf, "X.", 0))
	tl := NewLogger(logger)

	tl.LogError(ctx, fmt.Errorf("one"))
	tl.TraceLogging(false)
	tl.LogError(ctx, fmt.Errorf("two"))
	tl.TraceLogging(true)
	tl.LogError(ctx, fmt.Errorf("three"))
	tl.LogIfError(ctx, nil)
	tl.LogIfError(ctx, fmt.Errorf("four"))

	s := buf.String()
	g.Expect(s).To(Equal("X.Error [error:one]\nX.Error [error:three]\nX.Error [error:four]\n"))
}

func TestListTables(t *testing.T) {
	g := NewGomegaWithT(t)

	list, err := ListTables(gdb, nil)
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

	_, aid2, _, _ := insertFixtures(t, gdb)

	q := gdb.Dialect().ReplacePlaceholders("select xlines from pfx_addresses where id=?", nil)
	row := gdb.QueryRow(context.Background(), q, aid2)

	var xlines string
	err := row.Scan(&xlines)
	g.Expect(err).NotTo(HaveOccurred())
	g.Expect(xlines).To(Equal("2 Nutmeg Lane"))
}

func TestQueryContext(t *testing.T) {
	g := NewGomegaWithT(t)

	_, aid2, _, _ := insertFixtures(t, gdb)

	q := gdb.Dialect().ReplacePlaceholders("select xlines from pfx_addresses where id=?", nil)
	rows, err := gdb.Query(context.Background(), q, aid2)
	g.Expect(err).NotTo(HaveOccurred())
	g.Expect(rows.Next()).To(BeTrue())

	var xlines string
	err = rows.Scan(&xlines)
	g.Expect(err).NotTo(HaveOccurred())
	g.Expect(xlines).To(Equal("2 Nutmeg Lane"))

	g.Expect(rows.Next()).NotTo(BeTrue())
}

func TestSingleConnQuery(t *testing.T) {
	g := NewGomegaWithT(t)

	ctx := context.Background()
	_, aid2, _, _ := insertFixtures(t, gdb)

	q := gdb.Dialect().ReplacePlaceholders("select xlines from pfx_addresses where id=?", nil)
	e2 := gdb.SingleConn(ctx, func(ex Execer) error {
		rows, err := ex.Query(ctx, q, aid2)
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

	ctx := context.Background()
	insertFixtures(t, gdb)

	q := gdb.Dialect().ReplacePlaceholders("INSERT INTO pfx_addresses (xlines, postcode) VALUES (?, ?)", nil)
	err := gdb.Transact(ctx, nil, func(tx SqlTx) error {
		for i := 1; i <= 10; i++ {
			_, e2 := tx.Insert(ctx, "id", q, fmt.Sprintf("%d Pantagon Vale", i), "FX1 5EE")
			if e2 != nil {
				return e2
			}
		}
		return nil
	})
	g.Expect(err).NotTo(HaveOccurred())

	row := gdb.QueryRow(ctx, "select count(1) from pfx_addresses")

	var count int
	err = row.Scan(&count)
	g.Expect(err).NotTo(HaveOccurred())
	g.Expect(count).To(Equal(14))
}

func TestTransactCommitUsingExec(t *testing.T) {
	g := NewGomegaWithT(t)

	ctx := context.Background()
	_, aid2, aid3, _ := insertFixtures(t, gdb)

	q := gdb.Dialect().ReplacePlaceholders("delete from pfx_addresses where id in(?,?)", nil)
	err := gdb.Transact(ctx, nil, func(tx SqlTx) error {
		_, e2 := tx.Exec(ctx, q, aid2, aid3)
		return e2
	})
	g.Expect(err).NotTo(HaveOccurred())

	row := gdb.QueryRow(ctx, "select count(1) from pfx_addresses")

	var count int
	err = row.Scan(&count)
	g.Expect(err).NotTo(HaveOccurred())
	g.Expect(count).To(Equal(2))
}

func TestTransactRollback(t *testing.T) {
	g := NewGomegaWithT(t)

	ctx := context.Background()
	_, aid2, aid3, _ := insertFixtures(t, gdb)

	q := gdb.Dialect().ReplacePlaceholders("delete from pfx_addresses where id in(?,?)", nil)
	err := gdb.Transact(ctx, nil, func(tx SqlTx) error {
		tx.Exec(ctx, q, aid2, aid3)
		return errors.New("Bang")
	})
	g.Expect(err).To(HaveOccurred())
	g.Expect(err.Error()).To(Equal("Bang"))

	row := gdb.QueryRow(ctx, "select count(1) from pfx_addresses")

	var count int
	err = row.Scan(&count)
	g.Expect(err).NotTo(HaveOccurred())
	g.Expect(count).To(Equal(4))
}

func TestUserItemWrapper(t *testing.T) {
	g := NewGomegaWithT(t)

	d2 := gdb.With("hello")
	g.Expect(gdb.UserItem()).To(BeNil())
	g.Expect(d2.UserItem().(string)).To(Equal("hello"))
}

//-------------------------------------------------------------------------------------------------

type simpleLogger struct{}

func (l simpleLogger) Log(args ...interface{}) {
	if testing.Verbose() {
		log.Println(args...)
	}
}

//-------------------------------------------------------------------------------------------------

func TestMain(m *testing.M) {
	flag.Parse()

	var lvl pgx.LogLevel = pgx.LogLevelWarn
	if testing.Verbose() {
		lvl = pgx.LogLevelInfo
	}

	// first connection attempt: environment config for local DB
	testenv.SetEnvironmentForLocalPostgres()
	testUsingLocalDB(m, lvl)

	// second connection attempt: start up dockerised DB and use it
	testUsingDockertest(m, lvl)
}

func testUsingLocalDB(m *testing.M, lvl pgx.LogLevel) {
	lgr := testingadapter.NewLogger(simpleLogger{})
	var err error
	gdb, err = ConnectEnv(context.Background(), lgr, lvl)
	if err == nil {
		log.Printf("Test using local DB\n")
		os.Exit(m.Run())
	}

	var connErr *net.OpError
	if !errors.As(err, &connErr) {
		log.Fatalf("Cannot connect via env: %s", err)
	}
}

func testUsingDockertest(m *testing.M, lvl pgx.LogLevel) {
	testenv.SetUpDockerDbForTest(m, "postgres", func() {
		var err error
		lgr := testingadapter.NewLogger(simpleLogger{})
		gdb, err = ConnectEnv(context.Background(), lgr, lvl)
		if err != nil {
			log.Fatalf("Could not connect to DB in docker+postgres: %s", err)
		}
	})
}
