package pgxapi

import (
	"context"
	"fmt"
	"time"

	"github.com/jackc/pgx"
	"github.com/pkg/errors"
	"github.com/rickb777/sqlapi/dialect"
	"github.com/rickb777/where"
)

type postgres struct {
	db     *pgx.ConnPool
	logger pgx.Logger
}

//var _ DB = new(postgres)

// Logger access the underlying logger.
func (pg *postgres) Logger() pgx.Logger {
	return pg.logger
}

// Log emits a log event, supporting an elapsed-time calculation and providing an easier
// way to supply data parameters as name,value pairs.
func (pg *postgres) Log(level pgx.LogLevel, msg string, startTime *time.Time, data ...interface{}) {
	m := make(map[string]interface{})
	if startTime != nil {
		took := time.Now().Sub(*startTime)
		m["took"] = took
	}
	for i := 1; i < len(data); i += 2 {
		k := data[i-1].(string)
		v := data[i]
		m[k] = v
	}
	pg.logger.Log(level, msg, m)
}

// QueryExReplacePlaceholders executes any query returning data rows,
// having first replaced all '?' placeholders with numbered ones.
func (pg *postgres) QueryExReplacePlaceholders(ctx context.Context, query string, options *pgx.QueryExOptions, args ...interface{}) (*pgx.Rows, error) {
	qr := dialect.Postgres.ReplacePlaceholders(query, nil)
	return pg.QueryEx(ctx, qr, options, args...)
}

// QueryEx executes any query returning data rows.
func (pg *postgres) QueryEx(ctx context.Context, query string, options *pgx.QueryExOptions, args ...interface{}) (*pgx.Rows, error) {
	return pg.db.QueryEx(ctx, query, options, args...)
}

// QueryRowExReplacePlaceholders executes any query returning one data row.
// having first replaced all '?' placeholders with numbered ones.
func (pg *postgres) QueryRowExReplacePlaceholders(ctx context.Context, query string, options *pgx.QueryExOptions, args ...interface{}) *pgx.Row {
	qr := dialect.Postgres.ReplacePlaceholders(query, nil)
	return pg.db.QueryRowEx(ctx, qr, options, args...)
}

// QueryRowEx executes any query returning one data row.
func (pg *postgres) QueryRowEx(ctx context.Context, query string, options *pgx.QueryExOptions, args ...interface{}) *pgx.Row {
	return pg.db.QueryRowEx(ctx, query, options, args...)
}

// Transact takes a function and executes it within a DB transaction.
func (pg *postgres) Transact(ctx context.Context, txOptions *pgx.TxOptions, fn func(Execer) error) (err error) {
	var pgxTx *pgx.Tx
	pgxTx, err = pg.db.BeginEx(ctx, txOptions)
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
			pg.Log(pgx.LogLevelError, fmt.Sprintf("panic recovered: %+v", p), nil)
			pgxTx.Rollback()
			err = errors.New("transaction was rolled back")

		} else if err != nil {
			pgxTx.Rollback()

		} else {
			err = pgxTx.Commit()
		}
	}()

	return fn(WrapTX(pgxTx))
}

// GetIntIntIndex reads two integer columns from a specified database table and returns an index built from them.
func (pg *postgres) GetIntIntIndex(ctx context.Context, tableName, column1, column2 string, wh where.Expression) (map[int64]int64, error) {
	whs, args := where.Where(wh)
	q := fmt.Sprintf("SELECT %s, %s from %s%s", column1, column2, tableName, whs)
	rows, err := pg.db.QueryEx(ctx, q, nil, args...)
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

func (pg *postgres) BeginBatch() *pgx.Batch {
	return pg.db.BeginBatch()
}

func (pg *postgres) Close() {
	pg.db.Close()
}
