package pgxapi

import (
	"bytes"
	"context"
	"fmt"
	"log"
	"net"
	"os"
	"strings"
	"sync"
	"testing"
	"time"

	"github.com/ory/dockertest/v3"
	"github.com/ory/dockertest/v3/docker"
	"github.com/pkg/errors"

	"github.com/jackc/pgx"
	"github.com/jackc/pgx/log/testingadapter"
	. "github.com/onsi/gomega"
	"github.com/rickb777/sqlapi/dialect"
	"github.com/rickb777/sqlapi/pgxapi/logadapter"
)

var lock = sync.Mutex{}
var db SqlDB

func TestLoggingOnOff(t *testing.T) {
	g := NewGomegaWithT(t)

	buf := &bytes.Buffer{}
	logger := logadapter.NewLogger(log.New(buf, "X.", 0))
	sh := &shim{lgr: &toggleLogger{lgr: logger, enabled: 1}}

	d := NewDatabase(sh, dialect.Sqlite, nil)
	lgr := d.Logger()
	lgr.LogError(errors.New("one"))
	lgr.LogError(errors.New("two"))

	s := buf.String()
	g.Expect(s).To(Equal("X.Error [error:one]\nX.Error [error:two]\n"))
}

func TestLoggingError(t *testing.T) {
	g := NewGomegaWithT(t)

	buf := &bytes.Buffer{}
	logger := logadapter.NewLogger(log.New(buf, "X.", 0))
	sh := &shim{lgr: &toggleLogger{lgr: logger, enabled: 1}}

	d := NewDatabase(sh, dialect.Sqlite, nil)
	lgr := d.Logger()
	lgr.LogIfError(nil)
	lgr.LogIfError(fmt.Errorf("four"))
	lgr.LogIfError(nil)

	s := buf.String()
	g.Expect(s).To(Equal("X.Error [error:four]\n"))
}

func TestListTables(t *testing.T) {
	g := NewGomegaWithT(t)

	d := newDatabase(t)

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

func TestSingleConnQueryContext(t *testing.T) {
	g := NewGomegaWithT(t)

	d := newDatabase(t)

	_, aid2, _, _ := insertFixtures(t, d)

	q := d.Dialect().ReplacePlaceholders("select xlines from pfx_addresses where id=?", nil)
	e2 := d.DB().SingleConn(nil, func(ex Execer) error {
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

	err := d.DB().Transact(ctx, nil, func(tx SqlTx) error {
		q := d.Dialect().ReplacePlaceholders("INSERT INTO pfx_addresses (xlines, postcode) VALUES (?, ?)", nil)
		_, e2 := tx.InsertContext(ctx, "id", q, "5 Pantagon Vale", "FX1 5EE")
		return e2
	})
	g.Expect(err).NotTo(HaveOccurred())

	row := d.DB().QueryRowContext(ctx, "select count(1) from pfx_addresses")

	var count int
	err = row.Scan(&count)
	g.Expect(err).NotTo(HaveOccurred())
	g.Expect(count).To(Equal(5))
}

func TestTransactCommitUsingExec(t *testing.T) {
	g := NewGomegaWithT(t)

	d := newDatabase(t)

	ctx := context.Background()
	_, aid2, aid3, _ := insertFixtures(t, d)

	err := d.DB().Transact(ctx, nil, func(tx SqlTx) error {
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

	ctx := context.Background()
	_, aid2, aid3, _ := insertFixtures(t, d)

	err := d.DB().Transact(ctx, nil, func(tx SqlTx) error {
		q := d.Dialect().ReplacePlaceholders("delete from pfx_addresses where id in(?,?)", nil)
		tx.ExecContext(ctx, q, aid2, aid3)
		return errors.New("Bang")
	})
	g.Expect(err).To(HaveOccurred())
	g.Expect(err.Error()).To(Equal("Bang"))

	row := d.DB().QueryRowContext(ctx, "select count(1) from pfx_addresses")

	var count int
	err = row.Scan(&count)
	g.Expect(err).NotTo(HaveOccurred())
	g.Expect(count).To(Equal(4))
}

//-------------------------------------------------------------------------------------------------

func newDatabase(t *testing.T) Database {
	if db == nil {
		return nil
	}

	d := NewDatabase(db, dialect.Postgres, nil)
	if !testing.Verbose() {
		d.Logger().TraceLogging(false)
	}
	return d
}

type simpleLogger struct{}

func (l simpleLogger) Log(args ...interface{}) {
	log.Println(args...)
}

//-------------------------------------------------------------------------------------------------

func TestMain(m *testing.M) {
	// first connection attempt: environment config for local DB
	testUsingLocalDB(m)

	// second connection attempt: connect to DB provided by TravisCI
	testUsingTravisCiDB(m)

	// third connection attempt: start up dockerised DB and use it
	testUsingDockertest(m)
}

func testUsingLocalDB(m *testing.M) {
	log.Println("Attempting to connect to local postgresql")

	lgr := testingadapter.NewLogger(simpleLogger{})
	poolConfig := ParseEnvConfig()
	poolConfig.Logger = lgr
	poolConfig.LogLevel = pgx.LogLevelInfo

	pgxdb, err := pgx.NewConnPool(poolConfig)
	if err == nil {
		db = WrapDB(pgxdb, lgr)
		os.Exit(m.Run())
	}

	var connErr *net.OpError
	if !errors.As(err, &connErr) {
		log.Fatalf("Cannot connect via env: %s", err)
	}
}

func testUsingTravisCiDB(m *testing.M) {
	log.Println("Attempting to connect to local postgresql (TravisCI)")

	lgr := testingadapter.NewLogger(simpleLogger{})
	poolConfig := pgx.ConnPoolConfig{
		ConnConfig: pgx.ConnConfig{
			Host:     "localhost",
			Port:     5432,
			Database: "postgres",
			User:     "postgres",
			Password: "",
			Logger:   lgr,
			LogLevel: pgx.LogLevelInfo,
		},
	}

	pgxdb, err := pgx.NewConnPool(poolConfig)
	if err == nil {
		db = WrapDB(pgxdb, lgr)
		os.Exit(m.Run())
	}

	var connErr *net.OpError
	if !errors.As(err, &connErr) {
		log.Fatalf("Cannot connect: %s", err)
	}
}

func testUsingDockertest(m *testing.M) {
	log.Println("Attempting to connect to docker postgresql")

	// uses a sensible default on windows (tcp/http) and linux/osx (socket)
	pool, err := dockertest.NewPool("")
	if err != nil {
		log.Fatalf("Could not connect to docker: %s", err)
	}

	// pulls an image, creates a container based on it and runs it
	opts := &dockertest.RunOptions{
		Repository: "postgres",
		Tag:        "12-alpine",
		PortBindings: map[docker.Port][]docker.PortBinding{
			"5432/tcp": {{HostPort: "15432/tcp"}},
		},
		Env: []string{"PGPASSWORD=x", "POSTGRES_PASSWORD=x"},
	}
	resource, err := pool.RunWithOptions(opts)
	if err != nil {
		log.Fatalf("Could not start resource: %s", err)
	}

	// docker always takes some time to start
	time.Sleep(1950 * time.Millisecond)

	poolConfig := pgx.ConnPoolConfig{
		ConnConfig: pgx.ConnConfig{
			Host:     "localhost",
			Port:     15432,
			Database: "postgres",
			User:     "postgres",
			Password: "x",
			Logger:   testingadapter.NewLogger(simpleLogger{}),
			LogLevel: pgx.LogLevelInfo,
		},
	}

	pool.MaxWait = 10 * time.Second
	// exponential backoff-retry, because the application in the container might not be ready to accept connections yet
	if err = pool.Retry(func() error {
		db, err = Connect(poolConfig)
		if err != nil {
			return err
		}
		return db.PingContext(context.Background())
	}); err != nil {
		log.Fatalf("Could not connect to docker: %s", err)
	}

	code := m.Run()

	// You can't defer this because os.Exit doesn't care for defer
	if err = pool.Purge(resource); err != nil {
		log.Fatalf("Could not purge resource: %s", err)
	}

	os.Exit(code)
}
