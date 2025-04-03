package pgxapi

import (
	"bytes"
	"context"
	"errors"
	"fmt"
	"log"
	"slices"
	"strings"
	"testing"

	"github.com/jackc/pgx/v5/tracelog"
	"github.com/rickb777/expect"
	"github.com/rickb777/sqlapi"
	"github.com/rickb777/sqlapi/pgxapi/logadapter"
	"github.com/rickb777/sqlapi/support/testenv"
)

var gdb SqlDB

func TestLoggingOnOff(t *testing.T) {
	ctx := context.Background()
	buf := &bytes.Buffer{}
	logger := logadapter.NewLogger(log.New(buf, "X.", 0))
	tl := NewLogger(logger)

	tl.LogQuery(ctx, "one") // silently dropped
	tl.Log(ctx, tracelog.LogLevelInfo, "two", nil)
	tl.TraceLogging(false)
	tl.Log(ctx, tracelog.LogLevelInfo, "three", nil)
	tl.TraceLogging(true)
	tl.Log(ctx, tracelog.LogLevelInfo, "four", nil)

	s := buf.String()
	expect.String(s).ToBe(t, "X.two []\nX.four []\n")
}

func TestLoggingError(t *testing.T) {
	ctx := context.Background()
	buf := &bytes.Buffer{}
	logger := logadapter.NewLogger(log.New(buf, "X.", 0))
	tl := sqlapi.NewLogger(logger)

	tl.LogError(ctx, fmt.Errorf("one"))
	tl.TraceLogging(false)
	tl.LogError(ctx, fmt.Errorf("two"))
	tl.TraceLogging(true)
	tl.LogError(ctx, fmt.Errorf("three"))
	tl.LogIfError(ctx, nil)
	tl.LogIfError(ctx, fmt.Errorf("four"))

	s := buf.String()
	expect.String(s).ToBe(t, "X.Error [error:one]\nX.Error [error:three]\nX.Error [error:four]\n")
}

func TestListTables(t *testing.T) {
	list, err := ListTables(gdb, nil)
	expect.Error(err).Not().ToHaveOccurred(t)
	expect.Slice(slices.DeleteFunc(list, func(s string) bool {
		return !strings.HasPrefix(s, "sql_")
	})).ToBeEmpty(t)
	expect.Slice(slices.DeleteFunc(list, func(s string) bool {
		return !strings.HasPrefix(s, "pg_")
	})).ToBeEmpty(t)
}

func TestQueryRowContext(t *testing.T) {
	_, aid2, _, _ := insertFixtures(t, gdb)

	q := gdb.Dialect().ReplacePlaceholders("select xlines from pfx_addresses where id=?", nil)
	row := gdb.QueryRow(context.Background(), q, aid2)

	var xlines string
	err := row.Scan(&xlines)
	expect.Error(err).Not().ToHaveOccurred(t)
	expect.String(xlines).ToBe(t, "2 Nutmeg Lane")
}

func TestQueryContext(t *testing.T) {
	_, aid2, _, _ := insertFixtures(t, gdb)

	q := gdb.Dialect().ReplacePlaceholders("select xlines from pfx_addresses where id=?", nil)
	rows, err := gdb.Query(context.Background(), q, aid2)
	expect.Error(err).Not().ToHaveOccurred(t)
	expect.Bool(rows.Next()).ToBeTrue(t)

	var xlines string
	err = rows.Scan(&xlines)
	expect.Error(err).Not().ToHaveOccurred(t)
	expect.String(xlines).ToBe(t, "2 Nutmeg Lane")

	expect.Bool(rows.Next()).Not().ToBeTrue(t)
}

func TestSingleConnQuery(t *testing.T) {
	ctx := context.Background()
	_, aid2, _, _ := insertFixtures(t, gdb)

	q := gdb.Dialect().ReplacePlaceholders("select xlines from pfx_addresses where id=?", nil)
	e2 := gdb.SingleConn(ctx, func(ex Execer) error {
		rows, err := ex.Query(ctx, q, aid2)
		expect.Error(err).Not().ToHaveOccurred(t)
		expect.Bool(rows.Next()).ToBeTrue(t)

		var xlines string
		err = rows.Scan(&xlines)
		expect.Error(err).Not().ToHaveOccurred(t)
		expect.String(xlines).ToBe(t, "2 Nutmeg Lane")

		expect.Bool(rows.Next()).Not().ToBeTrue(t)
		return err
	})
	expect.Error(e2).Not().ToHaveOccurred(t)
}

func TestTransactCommitUsingInsert(t *testing.T) {
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
	expect.Error(err).Not().ToHaveOccurred(t)

	row := gdb.QueryRow(ctx, "select count(1) from pfx_addresses")

	var count int
	err = row.Scan(&count)
	expect.Error(err).Not().ToHaveOccurred(t)
	expect.Number(count).ToBe(t, 14)
}

func TestTransactCommitUsingExec(t *testing.T) {
	ctx := context.Background()
	_, aid2, aid3, _ := insertFixtures(t, gdb)

	q := gdb.Dialect().ReplacePlaceholders("delete from pfx_addresses where id in(?,?)", nil)
	err := gdb.Transact(ctx, nil, func(tx SqlTx) error {
		_, e2 := tx.Exec(ctx, q, aid2, aid3)
		return e2
	})
	expect.Error(err).Not().ToHaveOccurred(t)

	row := gdb.QueryRow(ctx, "select count(1) from pfx_addresses")

	var count int
	err = row.Scan(&count)
	expect.Error(err).Not().ToHaveOccurred(t)
	expect.Number(count).ToBe(t, 2)
}

func TestTransactRollback(t *testing.T) {
	ctx := context.Background()
	_, aid2, aid3, _ := insertFixtures(t, gdb)

	q := gdb.Dialect().ReplacePlaceholders("delete from pfx_addresses where id in(?,?)", nil)
	err := gdb.Transact(ctx, nil, func(tx SqlTx) error {
		tx.Exec(ctx, q, aid2, aid3)
		return errors.New("Bang")
	})
	expect.Error(err).ToContain(t, "Bang")

	row := gdb.QueryRow(ctx, "select count(1) from pfx_addresses")

	var count int
	err = row.Scan(&count)
	expect.Error(err).Not().ToHaveOccurred(t)
	expect.Number(count).ToBe(t, 4)
}

func TestUserItemWrapper(t *testing.T) {
	d2 := gdb.With("hello")
	expect.Any(gdb.UserItem()).ToBeNil(t)
	expect.String(d2.UserItem().(string)).ToBe(t, "hello")
}

//-------------------------------------------------------------------------------------------------

func TestMain(m *testing.M) {
	testenv.SetDefaultDbDriver("pgx")
	testenv.Shebang(m, func(lgr tracelog.Logger, logLevel tracelog.LogLevel, tries int) (err error) {
		gdb, err = ConnectEnv(context.Background(), lgr, logLevel, tries)
		return err
	})
}
