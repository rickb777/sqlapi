package sqlapi

import (
	"context"
	"database/sql"
	"github.com/pkg/errors"
)

func WrapDB(ex basicExecer) Execer {
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

type shim struct {
	ex   basicExecer
	isTx bool
}

var _ SqlDB = new(shim)
var _ SqlTx = new(shim)

func (sh *shim) ExecContext(ctx context.Context, query string, args ...interface{}) (int64, error) {
	res, err := sh.ex.ExecContext(ctx, query, args...)
	if err != nil {
		return 0, errors.Wrapf(err, "%s %v", query, args)
	}
	return res.RowsAffected()
}

func (sh *shim) InsertContext(ctx context.Context, query string, args ...interface{}) (int64, error) {
	res, err := sh.ex.ExecContext(ctx, query, args...)
	if err != nil {
		return 0, errors.Wrapf(err, "%s %v", query, args)
	}
	return res.LastInsertId()
}

func (sh *shim) PrepareContext(ctx context.Context, name, query string) (SqlStmt, error) {
	return sh.ex.PrepareContext(ctx, query)
}

func (sh *shim) QueryContext(ctx context.Context, query string, args ...interface{}) (SqlRows, error) {
	return sh.ex.QueryContext(ctx, query, args...)
}

func (sh *shim) QueryRowContext(ctx context.Context, query string, args ...interface{}) SqlRow {
	return sh.ex.QueryRowContext(ctx, query, args...)
}

func (sh *shim) BeginTx(ctx context.Context, opts *sql.TxOptions) (SqlTx, error) {
	tx, err := sh.ex.(*sql.DB).BeginTx(ctx, opts)
	if err != nil {
		return nil, errors.WithStack(err)
	}
	return &shim{ex: tx, isTx: true}, nil
}

func (sh *shim) PingContext(ctx context.Context) error {
	return sh.ex.(*sql.DB).PingContext(ctx)
}

func (sh *shim) Stats() sql.DBStats {
	return sh.ex.(*sql.DB).Stats()
}

func (sh *shim) Commit() error {
	return sh.ex.(*sql.Tx).Commit()
}

func (sh *shim) Rollback() error {
	return sh.ex.(*sql.Tx).Rollback()
}

func (sh *shim) IsTx() bool {
	return sh.isTx
}
