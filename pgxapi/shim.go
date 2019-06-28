package pgxapi

import (
	"context"
	"fmt"
	"github.com/rickb777/where"
	"time"

	"github.com/jackc/pgx"
	"github.com/pkg/errors"
	"github.com/rickb777/sqlapi/dialect"
)

type basicExecer interface {
	QueryEx(ctx context.Context, sql string, options *pgx.QueryExOptions, args ...interface{}) (*pgx.Rows, error)
	QueryRowEx(ctx context.Context, query string, options *pgx.QueryExOptions, args ...interface{}) *pgx.Row
	ExecEx(ctx context.Context, sql string, options *pgx.QueryExOptions, arguments ...interface{}) (pgx.CommandTag, error)
	PrepareEx(ctx context.Context, name, sql string, opts *pgx.PrepareExOptions) (*pgx.PreparedStatement, error)
	BeginBatch() *pgx.Batch
}

var _ basicExecer = new(pgx.ConnPool)
var _ basicExecer = new(pgx.Tx)

type shim struct {
	ex   basicExecer
	lgr  *toggleLogger
	isTx bool
}

var _ SqlDB = new(shim)
var _ SqlTx = new(shim)

// QueryContext executes any query returning data rows,
// having first replaced all '?' placeholders with numbered ones.
func (sh *shim) QueryContext(ctx context.Context, query string, args ...interface{}) (SqlRows, error) {
	qr := dialect.Postgres.ReplacePlaceholders(query, nil)
	return sh.QueryExRaw(ctx, qr, nil, args...)
}

// QueryExRaw executes any query returning data rows.
func (sh *shim) QueryExRaw(ctx context.Context, query string, options *pgx.QueryExOptions, args ...interface{}) (SqlRows, error) {
	rows, err := sh.ex.QueryEx(ctx, query, options, args...)
	if err != nil {
		return nil, errors.Wrapf(err, "%s %v", query, args)
	}
	return rows, nil
}

// QueryRowContext executes any query returning one data row,
// having first replaced all '?' placeholders with numbered ones.
func (sh *shim) QueryRowContext(ctx context.Context, query string, args ...interface{}) SqlRow {
	qr := dialect.Postgres.ReplacePlaceholders(query, nil)
	return sh.QueryRowExRaw(ctx, qr, nil, args...)
}

// QueryRowExRaw executes any query returning one data row.
func (sh *shim) QueryRowExRaw(ctx context.Context, query string, options *pgx.QueryExOptions, args ...interface{}) SqlRow {
	return sh.ex.QueryRowEx(ctx, query, options, args...)
}

// InsertContext executes an insert query returning an ID.
func (sh *shim) InsertContext(ctx context.Context, query string, args ...interface{}) (int64, error) {
	row := sh.ex.QueryRowEx(ctx, query, nil, args...)
	var id int64
	err := row.Scan(&id)
	if err != nil && err != pgx.ErrNoRows {
		return 0, errors.Wrapf(err, "%s %v", query, args)
	}
	return id, errors.Wrapf(err, "%s %v", query, args)
}

// ExecContext executes any query returning nothing.
func (sh *shim) ExecContext(ctx context.Context, query string, args ...interface{}) (int64, error) {
	tag, err := sh.ex.ExecEx(ctx, query, nil, args...)
	if err != nil {
		return 0, errors.Wrapf(err, "%s %v", query, args)
	}
	return tag.RowsAffected(), nil
}

// PrepareContext prepares a statement for repeated use.
func (sh *shim) PrepareContext(ctx context.Context, name, sql string) (*pgx.PreparedStatement, error) {
	ps, err := sh.ex.PrepareEx(ctx, name, sql, nil)
	return ps, errors.Wrapf(err, "%s %s", name, sql)
}

// BeginBatch begins a batch operation.
func (sh *shim) BeginBatch() *pgx.Batch {
	return sh.ex.BeginBatch()
}

func (sh *shim) IsTx() bool {
	return sh.isTx
}

func (sh *shim) Log(level pgx.LogLevel, msg string, data map[string]interface{}) {
	sh.lgr.Log(level, msg, data)
}

func (sh *shim) TraceLogging(on bool) {
	sh.lgr.TraceLogging(on)
}

//-------------------------------------------------------------------------------------------------
// ConnPool-specific methods

func (sh *shim) BeginTx(ctx context.Context, opts *pgx.TxOptions) (SqlTx, error) {
	tx, err := sh.ex.(*pgx.ConnPool).BeginEx(ctx, opts)
	return &shim{ex: tx, lgr: sh.lgr, isTx: true}, err
}

// Transact takes a function and executes it within a DB transaction.
func (sh *shim) Transact(ctx context.Context, txOptions *pgx.TxOptions, fn func(Execer) error) (err error) {
	var pgxTx SqlTx
	pgxTx, err = sh.BeginTx(ctx, txOptions)
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
			sh.lgr.Log(pgx.LogLevelError, fmt.Sprintf("panic recovered: %+v", p), nil)
			pgxTx.Rollback()
			err = errors.New("transaction was rolled back")

		} else if err != nil {
			pgxTx.Rollback()

		} else {
			err = pgxTx.Commit()
		}
	}()

	return fn(pgxTx)
}

// GetIntIntIndex reads two integer columns from a specified database table and returns an index built from them.
func (sh *shim) GetIntIntIndex(ctx context.Context, tableName, column1, column2 string, wh where.Expression) (map[int64]int64, error) {
	whs, args := where.Where(wh)
	q := fmt.Sprintf("SELECT %s, %s from %s%s", column1, column2, tableName, whs)
	rows, err := sh.ex.QueryEx(ctx, q, nil, args...)
	if err != nil {
		return nil, errors.Wrapf(err, "%s %+v", q, args)
	}
	defer rows.Close()

	index := make(map[int64]int64)
	for rows.Next() {
		var k, v int64
		err = rows.Scan(&k, &v)
		if err != nil {
			return nil, errors.Wrapf(err, "%s %+v", q, args)
		}
		index[k] = v
	}
	return index, nil
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

// Log emits a log event, supporting an elapsed-time calculation and providing an easier
// way to supply data parameters as name,value pairs.
func (sh *shim) LogT(level pgx.LogLevel, msg string, startTime *time.Time, data ...interface{}) {
	sh.lgr.LogT(level, msg, startTime, data...)
}

//-------------------------------------------------------------------------------------------------
// TX-specific methods

func (sh *shim) Commit() error {
	return sh.ex.(*pgx.Tx).Commit()
}

func (sh *shim) Rollback() error {
	return sh.ex.(*pgx.Tx).Rollback()
}
