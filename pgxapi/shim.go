package pgxapi

import (
	"context"
	"errors"
	"fmt"
	"log"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgconn"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/jackc/pgx/v5/tracelog"
	"github.com/rickb777/sqlapi/driver"
	"github.com/rickb777/where/dialect"
	"github.com/rickb777/where/quote"
)

// WrapDB wraps a *pgx.ConnPool as SqlDB.
// The logger is optional and can be nil, which disables logging.
// The quoter is optional and can be nil, defaulting to no quotes.
func WrapDB(pool *pgxpool.Pool, lgr tracelog.Logger, quoter quote.Quoter) SqlDB {
	if quoter == nil {
		quoter = quote.NoQuoter
	}
	di := driver.Postgres().WithQuoter(quoter)
	return &shim{ex: pool, di: di, lgr: NewLogger(lgr), isTx: false}
}

type basicExecer interface {
	Query(ctx context.Context, query string, args ...any) (pgx.Rows, error)
	QueryRow(ctx context.Context, query string, args ...any) pgx.Row
	Exec(ctx context.Context, query string, args ...any) (pgconn.CommandTag, error)
	//PrepareEx(ctx context.Context, name, query string, opts *pgx.PrepareExOptions) (*pgx.PreparedStatement, error)
	//BeginBatch() *pgx.Batch
}

var _ basicExecer = new(pgx.Conn)
var _ basicExecer = new(pgxpool.Pool)

//var _ basicExecer = new(pgx.Tx)

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

func (sh *shim) Query(ctx context.Context, query string, args ...any) (SqlRows, error) {
	qr := driver.Postgres().ReplacePlaceholders(query, nil)
	rows, err := sh.ex.Query(defaultCtx(ctx), qr, args...)
	if err != nil {
		return nil, wrap(err, query, args)
	}
	return rows, nil
}

func (sh *shim) QueryRow(ctx context.Context, query string, args ...any) SqlRow {
	qr := driver.Postgres().ReplacePlaceholders(query, nil)
	return sh.ex.QueryRow(defaultCtx(ctx), qr, args...)
}

func (sh *shim) Insert(ctx context.Context, pk, query string, args ...any) (int64, error) {
	q2 := fmt.Sprintf("%s RETURNING %s", query, pk)
	qr := driver.Postgres().ReplacePlaceholders(q2, nil)
	row := sh.ex.QueryRow(defaultCtx(ctx), qr, args...)
	var id int64
	err := row.Scan(&id)
	if err != nil && !errors.Is(err, pgx.ErrNoRows) {
		return 0, wrap(err, query, args)
	}
	return id, nil
}

func (sh *shim) Exec(ctx context.Context, query string, args ...any) (int64, error) {
	qr := dialect.ReplacePlaceholdersWithNumbers(query, "$")
	tag, err := sh.ex.Exec(defaultCtx(ctx), qr, args...)
	if err != nil {
		return 0, wrap(err, query, args)
	}
	return tag.RowsAffected(), nil
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
// pgx-specific methods

func (sh *shim) beginTx(ctx context.Context, opts *pgx.TxOptions) (SqlTx, error) {
	if opts == nil {
		opts = &pgx.TxOptions{}
	}
	tx, err := sh.ex.(*pgxpool.Pool).BeginTx(defaultCtx(ctx), *opts)
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
		return fn(sh) // nested transactions are inlined
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
	cp := sh.ex.(*pgxpool.Pool)
	conn, err := cp.Acquire(ctx)
	if err != nil {
		return err
	}

	defer func() {
		if p := recover(); p != nil {
			err = logPanicData(ctx, p, sh.lgr)
		}
		conn.Release()
	}()

	ex := &shim{
		ex:      conn,
		lgr:     sh.lgr,
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
	sh.ex.(*pgxpool.Pool).Close()
	return nil
}

func (sh *shim) Ping(ctx context.Context) error {
	ctx = defaultCtx(ctx)
	return sh.SingleConn(ctx, func(ex Execer) error {
		return ex.(*shim).ex.(*pgx.Conn).Ping(ctx)
	})
}

func (sh *shim) Stats() DBStats {
	return DBStats{} //sh.ex.(*pgx.ConnPool).Stats()
}

//-------------------------------------------------------------------------------------------------
// TX-specific methods

func (sh *shim) Commit(ctx context.Context) error {
	return sh.ex.(pgx.Tx).Commit(ctx)
}

func (sh *shim) Rollback(ctx context.Context) error {
	return sh.ex.(pgx.Tx).Rollback(ctx)
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
