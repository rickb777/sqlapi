package support

import (
	"database/sql"
	"fmt"
	"github.com/pkg/errors"
	"github.com/rickb777/sqlapi/pgxapi"
	"github.com/rickb777/sqlapi/require"
	"github.com/rickb777/where"
	"github.com/rickb777/where/quote"
	"strings"
)

// ReplaceTableName replaces all occurrences of "{TABLE}" with the table's name.
func ReplaceTableName(tbl pgxapi.Table, query string) string {
	return strings.Replace(query, "{TABLE}", tbl.Name().String(), -1)
}

// QueryOneNullThing queries for one cell of one record. Normally, the holder will be sql.NullString or similar.
// If required, the query can use "{TABLE}" in place of the table name.
func QueryOneNullThing(tbl pgxapi.Table, req require.Requirement, holder interface{}, query string, args ...interface{}) error {
	var n int64 = 0
	query = ReplaceTableName(tbl, query)
	database := tbl.Database()

	rows, err := Query(tbl, query, args...)
	if err != nil {
		return err
	}
	defer rows.Close()

	if rows.Next() {
		err = rows.Scan(holder)

		if err == sql.ErrNoRows {
			return database.Logger().LogIfError(require.ErrorIfQueryNotSatisfiedBy(req, 0))
		} else {
			n++
		}

		if rows.Next() {
			n++ // not singular
		}
	}

	return database.Logger().LogIfError(require.ChainErrorIfQueryNotSatisfiedBy(rows.Err(), req, n))
}

//-------------------------------------------------------------------------------------------------

func sliceSql(tbl pgxapi.Table, column string, wh where.Expression, qc where.QueryConstraint) (string, []interface{}) {
	q := tbl.Dialect().Quoter()
	whs, args := where.Where(wh, q)
	orderBy := where.Build(qc, q)
	return fmt.Sprintf("SELECT %s FROM %s %s %s",
		q.Quote(column), q.Quote(tbl.Name().String()), whs, orderBy), args
}

//-------------------------------------------------------------------------------------------------

// Query is the low-level request method for this table.
//
// The query is logged using whatever logger is configured. If an error arises, this too is logged.
//
// The args are for any placeholder parameters in the query.
//
// The caller must call rows.Close() on the result.
func Query(tbl pgxapi.Table, query string, args ...interface{}) (pgxapi.SqlRows, error) {
	q2 := tbl.Dialect().ReplacePlaceholders(query, args)
	lgr := tbl.Database().Logger()
	lgr.LogQuery(q2, args...)
	rows, err := tbl.Execer().QueryContext(tbl.Ctx(), q2, args...)
	return rows, lgr.LogIfError(errors.Wrapf(err, "%s %+v", q2, args))
}

// Exec executes a modification query (insert, update, delete, etc) and returns the number of items affected.
//
// The query is logged using whatever logger is configured. If an error arises, this too is logged.
func Exec(tbl pgxapi.Table, req require.Requirement, query string, args ...interface{}) (int64, error) {
	q2 := tbl.Dialect().ReplacePlaceholders(query, args)
	lgr := tbl.Database().Logger()
	lgr.LogQuery(q2, args...)
	n, err := tbl.Execer().ExecContext(tbl.Ctx(), q2, args...)
	if err != nil {
		return 0, lgr.LogError(errors.Wrapf(err, "%s %+v", q2, args))
	}
	return n, require.ChainErrorIfExecNotSatisfiedBy(err, req, n)
}

// UpdateFields writes certain fields of all the records matching a 'where' expression.
func UpdateFields(tbl pgxapi.Table, req require.Requirement, wh where.Expression, fields ...sql.NamedArg) (int64, error) {
	query, args := updateFieldsSQL(tbl.Name().String(), tbl.Dialect().Quoter(), wh, fields...)
	return Exec(tbl, req, query, args...)
}

func updateFieldsSQL(tblName string, q quote.Quoter, wh where.Expression, fields ...sql.NamedArg) (string, []interface{}) {
	list := pgxapi.NamedArgList(fields)
	assignments := strings.Join(list.Assignments(q, 1), ", ")
	whs, wargs := where.Where(wh, q)
	query := fmt.Sprintf("UPDATE %s SET %s %s", q.Quote(tblName), assignments, whs)
	args := append(list.Values(), wargs...)
	return query, args
}

// GetIntIntIndex reads two integer columns from a specified database table and returns an index built from them.
func GetIntIntIndex(tbl pgxapi.Table, q quote.Quoter, keyColumn, valColumn string, wh where.Expression) (map[int64]int64, error) {
	whs, args := where.Where(wh)
	query := fmt.Sprintf("SELECT %s, %s FROM %s %s", q.Quote(keyColumn), q.Quote(valColumn), q.Quote(tbl.Name().String()), whs)
	q2 := tbl.Dialect().ReplacePlaceholders(query, args)
	lgr := tbl.Database().Logger()
	lgr.LogQuery(q2, args...)
	rows, err := tbl.Execer().QueryContext(tbl.Ctx(), q2, args...)
	if err != nil {
		return nil, lgr.LogError(errors.Wrapf(err, "%s %+v", q2, args))
	}
	defer rows.Close()

	index := make(map[int64]int64)
	for rows.Next() {
		var k, v int64
		err = rows.Scan(&k, &v)
		if err != nil {
			return nil, lgr.LogError(errors.Wrapf(err, "%s %+v", q2, args))
		}
		index[k] = v
	}
	return index, nil
}

// GetStringIntIndex reads a string column and an integer column from a specified database table and returns an index built from them.
func GetStringIntIndex(tbl pgxapi.Table, q quote.Quoter, keyColumn, valColumn string, wh where.Expression) (map[string]int64, error) {
	whs, args := where.Where(wh)
	query := fmt.Sprintf("SELECT %s, %s FROM %s %s", q.Quote(keyColumn), q.Quote(valColumn), q.Quote(tbl.Name().String()), whs)
	q2 := tbl.Dialect().ReplacePlaceholders(query, args)
	lgr := tbl.Database().Logger()
	lgr.LogQuery(q2, args...)
	rows, err := tbl.Execer().QueryContext(tbl.Ctx(), q2, args...)
	if err != nil {
		return nil, lgr.LogError(errors.Wrapf(err, "%s %+v", q2, args))
	}
	defer rows.Close()

	index := make(map[string]int64)
	for rows.Next() {
		var k string
		var v int64
		err = rows.Scan(&k, &v)
		if err != nil {
			return nil, lgr.LogError(errors.Wrapf(err, "%s %+v", q2, args))
		}
		index[k] = v
	}
	return index, nil
}

// GetIntStringIndex reads an integer column and a string column from a specified database table and returns an index built from them.
func GetIntStringIndex(tbl pgxapi.Table, q quote.Quoter, keyColumn, valColumn string, wh where.Expression) (map[int64]string, error) {
	whs, args := where.Where(wh)
	query := fmt.Sprintf("SELECT %s, %s FROM %s %s", q.Quote(keyColumn), q.Quote(valColumn), q.Quote(tbl.Name().String()), whs)
	q2 := tbl.Dialect().ReplacePlaceholders(query, args)
	lgr := tbl.Database().Logger()
	lgr.LogQuery(q2, args...)
	rows, err := tbl.Execer().QueryContext(tbl.Ctx(), q2, args...)
	if err != nil {
		return nil, lgr.LogError(errors.Wrapf(err, "%s %+v", q2, args))
	}
	defer rows.Close()

	index := make(map[int64]string)
	for rows.Next() {
		var k int64
		var v string
		err = rows.Scan(&k, &v)
		if err != nil {
			return nil, lgr.LogError(errors.Wrapf(err, "%s %+v", q2, args))
		}
		index[k] = v
	}
	return index, nil
}
