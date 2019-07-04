package pgxapi

import (
	"context"
	"fmt"
	"github.com/jackc/pgx"
	"github.com/pkg/errors"
	"github.com/rickb777/sqlapi/dialect"
	"log"
)

func WrapDB(pool *pgx.ConnPool, lgr pgx.Logger) SqlDB {
	if lgr == nil {
		return &shim{ex: pool, isTx: false}
	}
	return &shim{ex: pool, lgr: &toggleLogger{lgr: lgr, enabled: 1}, isTx: false}
}

type basicExecer interface {
	QueryEx(ctx context.Context, sql string, options *pgx.QueryExOptions, args ...interface{}) (*pgx.Rows, error)
	QueryRowEx(ctx context.Context, query string, options *pgx.QueryExOptions, args ...interface{}) *pgx.Row
	ExecEx(ctx context.Context, sql string, options *pgx.QueryExOptions, arguments ...interface{}) (pgx.CommandTag, error)
	PrepareEx(ctx context.Context, name, sql string, opts *pgx.PrepareExOptions) (*pgx.PreparedStatement, error)
	BeginBatch() *pgx.Batch
}

var _ basicExecer = new(pgx.ConnPool)
var _ basicExecer = new(pgx.Tx)

//-------------------------------------------------------------------------------------------------

type shim struct {
	ex   basicExecer
	lgr  *toggleLogger
	isTx bool
}

var _ SqlDB = new(shim)
var _ SqlTx = new(shim)

//-------------------------------------------------------------------------------------------------

func (sh *shim) QueryContext(ctx context.Context, query string, args ...interface{}) (SqlRows, error) {
	qr := dialect.Postgres.ReplacePlaceholders(query, nil)
	return sh.QueryExRaw(ctx, qr, nil, args...)
}

func (sh *shim) QueryExRaw(ctx context.Context, query string, options *pgx.QueryExOptions, args ...interface{}) (SqlRows, error) {
	rows, err := sh.ex.QueryEx(ctx, query, options, args...)
	if err != nil {
		return nil, errors.Wrapf(err, "%s %v", query, args)
	}
	return rows, nil
}

func (sh *shim) QueryRowContext(ctx context.Context, query string, args ...interface{}) SqlRow {
	qr := dialect.Postgres.ReplacePlaceholders(query, nil)
	return sh.QueryRowExRaw(ctx, qr, nil, args...)
}

func (sh *shim) QueryRowExRaw(ctx context.Context, query string, options *pgx.QueryExOptions, args ...interface{}) SqlRow {
	return sh.ex.QueryRowEx(ctx, query, options, args...)
}

func (sh *shim) InsertContext(ctx context.Context, query string, args ...interface{}) (int64, error) {
	row := sh.ex.QueryRowEx(ctx, query, nil, args...)
	var id int64
	err := row.Scan(&id)
	if err != nil && err != pgx.ErrNoRows {
		return 0, errors.Wrapf(err, "%s %v", query, args)
	}
	return id, errors.Wrapf(err, "%s %v", query, args)
}

func (sh *shim) ExecContext(ctx context.Context, query string, args ...interface{}) (int64, error) {
	tag, err := sh.ex.ExecEx(ctx, query, nil, args...)
	if err != nil {
		return 0, errors.Wrapf(err, "%s %v", query, args)
	}
	return tag.RowsAffected(), nil
}

func (sh *shim) PrepareContext(ctx context.Context, name, query string) (*pgx.PreparedStatement, error) {
	ps, err := sh.ex.PrepareEx(ctx, name, query, nil)
	return ps, errors.Wrapf(err, "%s %s", name, query)
}

func (sh *shim) BeginBatch() *pgx.Batch {
	return sh.ex.BeginBatch()
}

func (sh *shim) IsTx() bool {
	return sh.isTx
}

func (sh *shim) Logger() Logger {
	return sh.lgr
}

//-------------------------------------------------------------------------------------------------
// ConnPool-specific methods

func (sh *shim) beginTx(ctx context.Context, opts *pgx.TxOptions) (SqlTx, error) {
	tx, err := sh.ex.(*pgx.ConnPool).BeginEx(ctx, opts)
	if err != nil {
		return nil, errors.WithStack(err)
	}
	return &shim{ex: tx, lgr: sh.lgr, isTx: true}, err
}

// Transact takes a function and executes it within a database transaction.
func (sh *shim) Transact(ctx context.Context, txOptions *pgx.TxOptions, fn func(SqlTx) error) (err error) {
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
			// using Sprintf so that the stack trace is printed (a feature of github.com/pkg/errors)
			if sh.lgr != nil {
				sh.lgr.Log(pgx.LogLevelError, fmt.Sprintf("panic recovered: %+v", p), nil)
			} else {
				log.Printf("panic recovered: %+v", p)
			}
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

func (sh *shim) Close() {
	sh.ex.(*pgx.ConnPool).Close()
}

func (sh *shim) PingContext(ctx context.Context) error {
	cp := sh.ex.(*pgx.ConnPool)
	conn, err := cp.Acquire()
	if err != nil {
		return err
	}
	defer cp.Release(conn)
	return conn.Ping(ctx)
}

func (sh *shim) Stats() DBStats {
	return DBStats{} //sh.ex.(*pgx.ConnPool).Stats()
}

//-------------------------------------------------------------------------------------------------
// TX-specific methods

func (sh *shim) commit() error {
	return sh.ex.(*pgx.Tx).Commit()
}

func (sh *shim) rollback() error {
	return sh.ex.(*pgx.Tx).Rollback()
}
