package sqlapi

import (
	"context"
	"database/sql"
	"fmt"
	"log"

	"github.com/jackc/pgx"

	"github.com/pkg/errors"
	"github.com/rickb777/sqlapi/dialect"
)

// WrapDB wraps a *sql.DB as SqlDB. The dialect is required.
// The logger is optional and can be nil, which disables logging.
func WrapDB(ex *sql.DB, di dialect.Dialect, lgr Logger) SqlDB {
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
	di      dialect.Dialect
	lgr     Logger
	isTx    bool
	wrapped interface{}
}

var _ SqlDB = new(shim)
var _ SqlTx = new(shim)

//-------------------------------------------------------------------------------------------------

func (sh *shim) QueryContext(ctx context.Context, query string, args ...interface{}) (SqlRows, error) {
	qr := sh.di.ReplacePlaceholders(query, nil)
	return sh.ex.QueryContext(defaultCtx(ctx), qr, args...)
}

func (sh *shim) QueryRowContext(ctx context.Context, query string, args ...interface{}) SqlRow {
	qr := sh.di.ReplacePlaceholders(query, nil)
	return sh.ex.QueryRowContext(defaultCtx(ctx), qr, args...)
}

func (sh *shim) InsertContext(ctx context.Context, pk, query string, args ...interface{}) (int64, error) {
	if sh.di.HasLastInsertId() {
		return sh.mysqlInsertContext(ctx, query, args...)
	}
	return sh.postgresInsertContext(ctx, pk, query, args...)
}

func (sh *shim) mysqlInsertContext(ctx context.Context, query string, args ...interface{}) (int64, error) {
	res, err := sh.ex.ExecContext(defaultCtx(ctx), query, args...)
	if err != nil {
		return 0, errors.Wrapf(err, "%s %v", query, args)
	}
	id, err := res.LastInsertId()
	return id, errors.Wrapf(err, "%s %v", query, args)
}

func (sh *shim) postgresInsertContext(ctx context.Context, pk, query string, args ...interface{}) (int64, error) {
	q2 := fmt.Sprintf("%s RETURNING %s", query, pk)
	qr := sh.di.ReplacePlaceholders(q2, nil)
	row := sh.ex.QueryRowContext(defaultCtx(ctx), qr, args...)
	var id int64
	err := row.Scan(&id)
	if err != nil && err != pgx.ErrNoRows {
		return 0, errors.Wrapf(err, "%s %v", query, args)
	}
	return id, errors.Wrapf(err, "%s %v", query, args)
}

func (sh *shim) ExecContext(ctx context.Context, query string, args ...interface{}) (int64, error) {
	qr := sh.di.ReplacePlaceholders(query, nil)
	res, err := sh.ex.ExecContext(defaultCtx(ctx), qr, args...)
	if err != nil {
		return 0, errors.Wrapf(err, "%s %v", query, args)
	}
	n, err := res.RowsAffected()
	return n, errors.Wrapf(err, "%s %v", query, args)
}

func (sh *shim) PrepareContext(ctx context.Context, name, query string) (SqlStmt, error) {
	qr := sh.di.ReplacePlaceholders(query, nil)
	ps, err := sh.ex.PrepareContext(defaultCtx(ctx), qr)
	return ps, errors.Wrapf(err, "%s %s", name, query)
}

func (sh *shim) IsTx() bool {
	return sh.isTx
}

func (sh *shim) Logger() Logger {
	return sh.lgr
}

func (sh *shim) Dialect() dialect.Dialect {
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

func (sh *shim) beginTx(ctx context.Context, opts *sql.TxOptions) (SqlTx, error) {
	tx, err := sh.ex.(*sql.DB).BeginTx(defaultCtx(ctx), opts)
	if err != nil {
		return nil, errors.WithStack(err)
	}
	cp := *sh
	cp.ex = tx
	cp.isTx = true
	return &cp, nil
}

// Transact takes a function and executes it within a database transaction.
func (sh *shim) Transact(ctx context.Context, txOptions *sql.TxOptions, fn func(SqlTx) error) (err error) {
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
			_ = tx.Rollback()
			err = logPanicData(p, sh.lgr)

		} else if err != nil {
			_ = tx.Rollback()

		} else {
			err = tx.Commit()
		}
	}()

	return fn(tx)
}

func (sh *shim) SingleConn(ctx context.Context, fn func(ex Execer) error) (err error) {
	cp := sh.ex.(*sql.DB)
	var conn *sql.Conn
	conn, err = cp.Conn(defaultCtx(ctx))
	if err != nil {
		return errors.WithStack(err)
	}

	defer func() {
		if p := recover(); p != nil {
			err = logPanicData(p, sh.lgr)
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

func logPanicData(p interface{}, lgr Logger) error {
	// capture a stack trace using github.com/pkg/errors
	if e, ok := p.(error); ok {
		p = errors.WithStack(e)
	} else {
		p = errors.Errorf("%+v", p)
	}
	// using Sprintf so that the stack trace is printed (a feature of github.com/pkg/errors)
	if lgr != nil {
		lgr.Log(fmt.Sprintf("panic recovered: %+v", p), nil)
	} else {
		log.Printf("panic recovered: %+v", p)
	}
	return p.(error)
}

func (sh *shim) Close() error {
	return sh.ex.(*sql.DB).Close()
}

func (sh *shim) PingContext(ctx context.Context) error {
	return sh.ex.(*sql.DB).PingContext(ctx)
}

func (sh *shim) Stats() sql.DBStats {
	return sh.ex.(*sql.DB).Stats()
}

//-------------------------------------------------------------------------------------------------
// TX-specific methods

func (sh *shim) Commit() error {
	return sh.ex.(*sql.Tx).Commit()
}

func (sh *shim) Rollback() error {
	return sh.ex.(*sql.Tx).Rollback()
}

func defaultCtx(ctx context.Context) context.Context {
	if ctx == nil {
		return context.Background()
	}
	return ctx
}
