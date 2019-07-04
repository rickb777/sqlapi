package sqlapi

import (
	"context"
	"database/sql"
	"github.com/pkg/errors"
	"log"
)

func WrapDB(ex *sql.DB) SqlDB {
	return &shim{ex: ex}
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
	ex   basicExecer
	isTx bool
}

var _ SqlDB = new(shim)
var _ SqlTx = new(shim)

//-------------------------------------------------------------------------------------------------

func (sh *shim) QueryContext(ctx context.Context, query string, args ...interface{}) (SqlRows, error) {
	return sh.ex.QueryContext(ctx, query, args...)
}

func (sh *shim) QueryRowContext(ctx context.Context, query string, args ...interface{}) SqlRow {
	return sh.ex.QueryRowContext(ctx, query, args...)
}

func (sh *shim) InsertContext(ctx context.Context, query string, args ...interface{}) (int64, error) {
	res, err := sh.ex.ExecContext(ctx, query, args...)
	if err != nil {
		return 0, errors.Wrapf(err, "%s %v", query, args)
	}
	id, err := res.LastInsertId()
	return id, errors.Wrapf(err, "%s %v", query, args)
}

func (sh *shim) ExecContext(ctx context.Context, query string, args ...interface{}) (int64, error) {
	res, err := sh.ex.ExecContext(ctx, query, args...)
	if err != nil {
		return 0, errors.Wrapf(err, "%s %v", query, args)
	}
	n, err := res.RowsAffected()
	return n, errors.Wrapf(err, "%s %v", query, args)
}

func (sh *shim) PrepareContext(ctx context.Context, name, query string) (SqlStmt, error) {
	ps, err := sh.ex.PrepareContext(ctx, query)
	return ps, errors.Wrapf(err, "%s %s", name, query)
}

func (sh *shim) IsTx() bool {
	return sh.isTx
}

//-------------------------------------------------------------------------------------------------
// sql.DB specific methods

func (sh *shim) beginTx(ctx context.Context, opts *sql.TxOptions) (SqlTx, error) {
	tx, err := sh.ex.(*sql.DB).BeginTx(ctx, opts)
	if err != nil {
		return nil, errors.WithStack(err)
	}
	return &shim{ex: tx, isTx: true}, nil
}

func (sh *shim) Transact(ctx context.Context, txOptions *sql.TxOptions, fn func(SqlTx) error) (err error) {
	if tx, isTx := sh.ex.(SqlTx); isTx {
		return fn(tx) // nested transactions are inlined
	}

	var tx SqlTx
	tx, err = sh.beginTx(ctx, txOptions)
	if err != nil {
		return err
	}

	defer func() {
		if p := recover(); p != nil {
			// capture a stack trace using github.com/pkg/errors
			if e, ok := p.(error); ok {
				p = errors.WithStack(e)
			} else {
				p = errors.Errorf("%+v", p)
			}
			log.Printf("panic recovered: %+v", p)
			tx.rollback()
			err = errors.New("transaction was rolled back")

		} else if err != nil {
			tx.rollback()

		} else {
			err = tx.commit()
		}
	}()

	return fn(tx)
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

func (sh *shim) commit() error {
	return sh.ex.(*sql.Tx).Commit()
}

func (sh *shim) rollback() error {
	return sh.ex.(*sql.Tx).Rollback()
}
