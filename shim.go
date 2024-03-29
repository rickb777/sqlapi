package sqlapi

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"log"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/tracelog"
	"github.com/rickb777/sqlapi/driver"
)

// WrapDB wraps a *sql.DB as SqlDB. The dialect is required.
// The logger is optional and can be nil, which disables logging.
func WrapDB(ex *sql.DB, di driver.Dialect, lgr Logger) SqlDB {
	return &shim{ex: ex, di: di, lgr: lgr, isTx: false}
}

type basicExecer interface {
	QueryContext(ctx context.Context, query string, args ...interface{}) (*sql.Rows, error)
	QueryRowContext(ctx context.Context, query string, args ...interface{}) *sql.Row
	ExecContext(ctx context.Context, query string, args ...interface{}) (sql.Result, error)
	PrepareContext(ctx context.Context, query string) (*sql.Stmt, error)
}

var _ basicExecer = new(sql.DB)
var _ basicExecer = new(sql.Tx)

//-------------------------------------------------------------------------------------------------

type shim struct {
	ex      basicExecer
	di      driver.Dialect
	lgr     Logger
	isTx    bool
	wrapped interface{}
}

var _ SqlDB = new(shim)
var _ SqlTx = new(shim)

//-------------------------------------------------------------------------------------------------

func (sh *shim) Query(ctx context.Context, query string, args ...interface{}) (SqlRows, error) {
	qr := sh.di.ReplacePlaceholders(query, nil)
	return sh.ex.QueryContext(defaultCtx(ctx), qr, args...)
}

func (sh *shim) QueryRow(ctx context.Context, query string, args ...interface{}) SqlRow {
	qr := sh.di.ReplacePlaceholders(query, nil)
	return sh.ex.QueryRowContext(defaultCtx(ctx), qr, args...)
}

func (sh *shim) Insert(ctx context.Context, pk, query string, args ...interface{}) (int64, error) {
	if sh.di.HasLastInsertId() {
		return sh.mysqlInsert(ctx, query, args...)
	}
	return sh.postgresInsert(ctx, pk, query, args...)
}

func (sh *shim) mysqlInsert(ctx context.Context, query string, args ...interface{}) (int64, error) {
	res, err := sh.ex.ExecContext(defaultCtx(ctx), query, args...)
	if err != nil {
		return 0, wrap(err, query, args)
	}
	id, err := res.LastInsertId()
	return id, wrap(err, query, args)
}

func (sh *shim) postgresInsert(ctx context.Context, pk, query string, args ...interface{}) (int64, error) {
	q2 := fmt.Sprintf("%s RETURNING %s", query, pk)
	qr := sh.di.ReplacePlaceholders(q2, nil)
	row := sh.ex.QueryRowContext(defaultCtx(ctx), qr, args...)
	var id int64
	err := row.Scan(&id)
	if err != nil && !errors.Is(err, pgx.ErrNoRows) {
		return 0, wrap(err, query, args)
	}
	return id, nil
}

func (sh *shim) Exec(ctx context.Context, query string, args ...interface{}) (int64, error) {
	qr := sh.di.ReplacePlaceholders(query, nil)
	res, err := sh.ex.ExecContext(defaultCtx(ctx), qr, args...)
	if err != nil {
		return 0, wrap(err, query, args)
	}
	n, err := res.RowsAffected()
	return n, wrap(err, query, args)
}

func (sh *shim) IsTx() bool {
	return sh.isTx
}

func (sh *shim) Logger() Logger {
	return sh.lgr
}

func (sh *shim) Dialect() driver.Dialect {
	return sh.di
}

func (sh *shim) With(userItem interface{}) SqlDB {
	cp := *sh
	cp.wrapped = userItem
	return &cp
}

func (sh *shim) UserItem() interface{} {
	return sh.wrapped
}

//-------------------------------------------------------------------------------------------------
// sql.DB specific methods

func convertTxOptions(pgopts *pgx.TxOptions) *sql.TxOptions {
	if pgopts == nil {
		return nil
	}

	iso := sql.LevelDefault
	switch pgopts.IsoLevel {
	case pgx.ReadCommitted:
		iso = sql.LevelReadCommitted
	case pgx.ReadUncommitted:
		iso = sql.LevelReadUncommitted
	case pgx.RepeatableRead:
		iso = sql.LevelRepeatableRead
	case pgx.Serializable:
		iso = sql.LevelSerializable
	}

	return &sql.TxOptions{
		Isolation: iso,
		ReadOnly:  pgopts.AccessMode == pgx.ReadOnly,
	}
}

func (sh *shim) beginTx(ctx context.Context, pgopts *pgx.TxOptions) (SqlTx, error) {
	opts := convertTxOptions(pgopts)
	tx, err := sh.ex.(*sql.DB).BeginTx(defaultCtx(ctx), opts)
	if err != nil {
		return nil, err
	}
	cp := *sh
	cp.ex = tx
	cp.isTx = true
	return &cp, nil
}

// Transact takes a function and executes it within a database transaction.
func (sh *shim) Transact(ctx context.Context, txOptions *pgx.TxOptions, fn func(SqlTx) error) (err error) {
	if sh.isTx {
		if _, isTx := sh.ex.(*sql.Tx); isTx {
			return fn(sh) // nested transactions are inlined
		}
	}

	var tx SqlTx
	tx, err = sh.beginTx(ctx, txOptions)
	if err != nil {
		return err
	}

	defer func() {
		if p := recover(); p != nil {
			_ = tx.Rollback(ctx)
			err = logPanicData(ctx, p, sh.lgr)

		} else if err != nil {
			_ = tx.Rollback(ctx)

		} else {
			err = tx.Commit(ctx)
		}
	}()

	return fn(tx)
}

func (sh *shim) SingleConn(ctx context.Context, fn func(ex Execer) error) (err error) {
	cp := sh.ex.(*sql.DB)
	var conn *sql.Conn
	conn, err = cp.Conn(defaultCtx(ctx))
	if err != nil {
		return err
	}

	defer func() {
		if p := recover(); p != nil {
			err = logPanicData(ctx, p, sh.lgr)
		}
		e2 := conn.Close()
		if err == nil {
			err = e2
		} // otherwise e2 is ignored
	}()

	ex := &shim{
		ex:      conn,
		lgr:     sh.lgr,
		di:      sh.di,
		isTx:    false,
		wrapped: sh.wrapped,
	}
	return fn(ex)
}

func logPanicData(ctx context.Context, p interface{}, lgr tracelog.Logger) error {
	// capture a stack trace using github.com/pkg/errors
	if e, ok := p.(error); ok {
		p = e
	} else {
		p = fmt.Errorf("%+v", p)
	}
	// using Sprintf so that the stack trace is printed (a feature of github.com/pkg/errors)
	if lgr != nil {
		lgr.Log(ctx, tracelog.LogLevelError, fmt.Sprintf("panic recovered: %+v", p), nil)
	} else {
		log.Printf("panic recovered: %+v", p)
	}
	return p.(error)
}

func (sh *shim) Close() error {
	return sh.ex.(*sql.DB).Close()
}

func (sh *shim) Ping(ctx context.Context) error {
	return sh.ex.(*sql.DB).PingContext(ctx)
}

func (sh *shim) Stats() sql.DBStats {
	return sh.ex.(*sql.DB).Stats()
}

//-------------------------------------------------------------------------------------------------
// TX-specific methods

func (sh *shim) Commit(_ context.Context) error {
	return sh.ex.(*sql.Tx).Commit()
}

func (sh *shim) Rollback(_ context.Context) error {
	return sh.ex.(*sql.Tx).Rollback()
}

func defaultCtx(ctx context.Context) context.Context {
	if ctx == nil {
		return context.Background()
	}
	return ctx
}

func wrap(err error, query string, args ...interface{}) error {
	if err != nil {
		return fmt.Errorf("%w %s %v", err, query, args)
	}
	return nil
}
