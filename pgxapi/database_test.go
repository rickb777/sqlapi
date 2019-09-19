package pgxapi

import (
	"bytes"
	"errors"
	"fmt"
	"github.com/jackc/pgx"
	"github.com/jackc/pgx/log/testingadapter"
	. "github.com/onsi/gomega"
	"github.com/rickb777/sqlapi/dialect"
	"github.com/rickb777/sqlapi/pgxapi/logadapter"
	"log"
	"strings"
	"testing"
)

// Environment:
// PGHOST, PGPORT, PGUSER, PGPASSWORD, PGDATABASE, PGCONNECT_TIMEOUT,
// PGSSLMODE, PGSSLKEY, PGSSLCERT, PGSSLROOTCERT.
// (see https://www.postgresql.org/docs/11/libpq-envars.html)

func connect(t *testing.T) SqlDB {
	lgr := testingadapter.NewLogger(t)
	db, err := ConnectEnv(lgr, pgx.LogLevelInfo)
	if err != nil {
		t.Log(err)
		t.Skip()
	}
	return db
}

func newDatabase(t *testing.T) Database {
	db := connect(t)
	if db == nil {
		return nil
	}

	d := NewDatabase(db, dialect.Postgres, nil)
	if !testing.Verbose() {
		d.Logger().TraceLogging(false)
	}
	return d
}

func cleanup(db SqlDB) {
	if db != nil {
		db.Close()
		db = nil
	}
}

func TestLoggingOnOff(t *testing.T) {
	g := NewGomegaWithT(t)

	buf := &bytes.Buffer{}
	logger := logadapter.NewLogger(log.New(buf, "X.", 0))
	sh := &shim{lgr: &toggleLogger{lgr: logger, enabled: 1}}

	db := NewDatabase(sh, dialect.Sqlite, nil)
	db.Logger().LogError(errors.New("one"))
	db.Logger().LogError(errors.New("two"))

	s := buf.String()
	g.Expect(s).To(Equal("X.Error [error:one]\nX.Error [error:two]\n"))
}

func TestLoggingError(t *testing.T) {
	g := NewGomegaWithT(t)

	buf := &bytes.Buffer{}
	logger := logadapter.NewLogger(log.New(buf, "X.", 0))
	sh := &shim{lgr: &toggleLogger{lgr: logger, enabled: 1}}

	db := NewDatabase(sh, dialect.Sqlite, nil)
	db.Logger().LogIfError(nil)
	db.Logger().LogIfError(fmt.Errorf("four"))
	db.Logger().LogIfError(nil)

	s := buf.String()
	g.Expect(s).To(Equal("X.Error [error:four]\n"))
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
