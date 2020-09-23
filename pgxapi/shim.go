package pgxapi

import (
	"context"
	"fmt"
	"log"

	"github.com/jackc/pgx"
	"github.com/pkg/errors"
	"github.com/rickb777/sqlapi/dialect"
	"github.com/rickb777/where/quote"
)

// WrapDB wraps a *pgx.ConnPool as SqlDB.
// The logger is optional and can be nil, which disables logging.
// The quoter is optional and can be nil, defaulting to no quotes.
func WrapDB(pool *pgx.ConnPool, lgr pgx.Logger, quoter quote.Quoter) SqlDB {
	if quoter == nil {
		quoter = quote.NoQuoter
	}
	di := dialect.Postgres.WithQuoter(quoter)
	return &shim{ex: pool, di: di, lgr: NewLogger(lgr), isTx: false}
}

type basicExecer interface {
	QueryEx(ctx context.Context, sql string, options *pgx.QueryExOptions, args ...interface{}) (*pgx.Rows, error)
	QueryRowEx(ctx context.Context, query string, options *pgx.QueryExOptions, args ...interface{}) *pgx.Row
	ExecEx(ctx context.Context, sql string, options *pgx.QueryExOptions, arguments ...interface{}) (pgx.CommandTag, error)
	PrepareEx(ctx context.Context, name, sql string, opts *pgx.PrepareExOptions) (*pgx.PreparedStatement, error)
	BeginBatch() *pgx.Batch
}

var _ basicExecer = new(pgx.Conn)
var _ basicExecer = new(pgx.ConnPool)
var _ basicExecer = new(pgx.Tx)

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
	qr := dialect.Postgres.ReplacePlaceholders(query, nil)
	return sh.QueryExRaw(ctx, qr, nil, args...)
}

func (sh *shim) QueryRowContext(ctx context.Context, query string, args ...interface{}) SqlRow {
	qr := dialect.Postgres.ReplacePlaceholders(query, nil)
	return sh.QueryRowExRaw(ctx, qr, nil, args...)
}

func (sh *shim) QueryExRaw(ctx context.Context, query string, options *pgx.QueryExOptions, args ...interface{}) (SqlRows, error) {
	qr := dialect.Postgres.ReplacePlaceholders(query, nil)
	rows, err := sh.ex.QueryEx(defaultCtx(ctx), qr, options, args...)
	if err != nil {
		return nil, errors.Wrapf(err, "%s %v", query, args)
	}
	return rows, nil
}

func (sh *shim) QueryRowExRaw(ctx context.Context, query string, options *pgx.QueryExOptions, args ...interface{}) SqlRow {
	qr := dialect.Postgres.ReplacePlaceholders(query, nil)
	return sh.ex.QueryRowEx(defaultCtx(ctx), qr, options, args...)
}

func (sh *shim) InsertContext(ctx context.Context, pk, query string, args ...interface{}) (int64, error) {
	q2 := fmt.Sprintf("%s RETURNING %s", query, pk)
	qr := dialect.Postgres.ReplacePlaceholders(q2, nil)
	row := sh.ex.QueryRowEx(defaultCtx(ctx), qr, nil, args...)
	var id int64
	err := row.Scan(&id)
	if err != nil && err != pgx.ErrNoRows {
		return 0, errors.Wrapf(err, "%s %v", query, args)
	}
	return id, errors.Wrapf(err, "%s %v", query, args)
}

func (sh *shim) ExecContext(ctx context.Context, query string, args ...interface{}) (int64, error) {
	qr := dialect.Postgres.ReplacePlaceholders(query, nil)
	tag, err := sh.ex.ExecEx(defaultCtx(ctx), qr, nil, args...)
	if err != nil {
		return 0, errors.Wrapf(err, "%s %v", query, args)
	}
	return tag.RowsAffected(), nil
}

func (sh *shim) PrepareContext(ctx context.Context, name, query string) (*pgx.PreparedStatement, error) {
	qr := dialect.Postgres.ReplacePlaceholders(query, nil)
	ps, err := sh.ex.PrepareEx(defaultCtx(ctx), name, qr, nil)
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
// ConnPool-specific methods

func (sh *shim) beginTx(ctx context.Context, opts *pgx.TxOptions) (SqlTx, error) {
	tx, err := sh.ex.(*pgx.ConnPool).BeginEx(defaultCtx(ctx), opts)
	if err != nil {
		return nil, errors.WithStack(err)
	}
	cp := *sh
	cp.ex = tx
	cp.isTx = true
	return &cp, nil
}

// Transact takes a function and executes it within a database transaction.
func (sh *shim) Transact(ctx context.Context, txOptions *pgx.TxOptions, fn func(SqlTx) error) (err error) {
	if sh.isTx {
		if tx, isTx := sh.ex.(SqlTx); isTx {
			return fn(tx) // nested transactions are inlined
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
	cp := sh.ex.(*pgx.ConnPool)
	var conn *pgx.Conn
	conn, err = cp.AcquireEx(defaultCtx(ctx))
	if err != nil {
		return errors.WithStack(err)
	}

	defer func() {
		if p := recover(); p != nil {
			err = logPanicData(p, sh.lgr)
		}
		cp.Release(conn)
	}()

	ex := &shim{
		ex:      conn,
		lgr:     sh.lgr,
		isTx:    false,
		wrapped: sh.wrapped,
	}
	return fn(ex)
}

func logPanicData(p interface{}, lgr pgx.Logger) error {
	// capture a stack trace using github.com/pkg/errors
	if e, ok := p.(error); ok {
		p = errors.WithStack(e)
	} else {
		p = errors.Errorf("%+v", p)
	}
	// using Sprintf so that the stack trace is printed (a feature of github.com/pkg/errors)
	if lgr != nil {
		lgr.Log(pgx.LogLevelError, fmt.Sprintf("panic recovered: %+v", p), nil)
	} else {
		log.Printf("panic recovered: %+v", p)
	}
	return p.(error)
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
	return conn.Ping(defaultCtx(ctx))
}

func (sh *shim) Stats() DBStats {
	return DBStats{} //sh.ex.(*pgx.ConnPool).Stats()
}

//-------------------------------------------------------------------------------------------------
// TX-specific methods

func (sh *shim) Commit() error {
	return sh.ex.(*pgx.Tx).Commit()
}

func (sh *shim) Rollback() error {
	return sh.ex.(*pgx.Tx).Rollback()
}

func defaultCtx(ctx context.Context) context.Context {
	if ctx == nil {
		return context.Background()
	}
	return ctx
}
