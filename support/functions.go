package support

import (
	"database/sql"
	"fmt"
	"strings"

	"github.com/rickb777/sqlapi"
	"github.com/rickb777/sqlapi/require"
	"github.com/rickb777/where"
	"github.com/rickb777/where/quote"
)

// ReplaceTableName replaces all occurrences of "{TABLE}" with the table's name.
func ReplaceTableName(tbl sqlapi.Table, query string) string {
	return strings.Replace(query, "{TABLE}", tbl.Name().String(), -1)
}

// QueryOneNullThing queries for one cell of one record. Normally, the holder will be sql.NullString or similar.
// If required, the query can use "{TABLE}" in place of the table name.
func QueryOneNullThing(tbl sqlapi.Table, req require.Requirement, holder interface{}, query string, args ...interface{}) error {
	var n int64 = 0
	query = ReplaceTableName(tbl, query)

	rows, err := Query(tbl, query, args...)
	if err != nil {
		return err
	}
	defer rows.Close()

	if rows.Next() {
		err = rows.Scan(holder)

		if err == sql.ErrNoRows {
			return tbl.Logger().LogIfError(tbl.Ctx(), require.ErrorIfQueryNotSatisfiedBy(req, 0))
		} else {
			n++
		}

		if rows.Next() {
			n++ // not singular
		}
	}

	return tbl.Logger().LogIfError(tbl.Ctx(), require.ChainErrorIfQueryNotSatisfiedBy(rows.Err(), req, n))
}

//-------------------------------------------------------------------------------------------------

func sliceSql(tbl sqlapi.Table, column string, wh where.Expression, qc where.QueryConstraint) (string, []interface{}) {
	q := tbl.Dialect().Quoter()
	whs, args := where.Where(wh, q)
	orderBy := where.Build(qc, tbl.Dialect().Index())
	return fmt.Sprintf("SELECT %s FROM %s%s%s",
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
func Query(tbl sqlapi.Table, query string, args ...interface{}) (sqlapi.SqlRows, error) {
	q2 := tbl.Dialect().ReplacePlaceholders(query, args)
	lgr := tbl.Logger()
	lgr.LogQuery(tbl.Ctx(), q2, args...)
	rows, err := tbl.Execer().Query(tbl.Ctx(), q2, args...)
	return rows, lgr.LogIfError(tbl.Ctx(), err)
}

// Exec executes a modification query (insert, update, delete, etc) and returns the number of items affected.
//
// The query is logged using whatever logger is configured. If an error arises, this too is logged.
func Exec(tbl sqlapi.Table, req require.Requirement, query string, args ...interface{}) (int64, error) {
	q2 := tbl.Dialect().ReplacePlaceholders(query, args)
	n, err := doExec(tbl, q2, args...)
	return n, require.ChainErrorIfExecNotSatisfiedBy(err, req, n)
}

func doExec(tbl sqlapi.Table, query string, args ...interface{}) (int64, error) {
	lgr := tbl.Logger()
	lgr.LogQuery(tbl.Ctx(), query, args...)
	n, err := tbl.Execer().Exec(tbl.Ctx(), query, args...)
	if err != nil {
		return 0, lgr.LogError(tbl.Ctx(), err)
	}
	return n, err
}

// UpdateFields writes certain fields of all the records matching a 'where' expression.
func UpdateFields(tbl sqlapi.Table, req require.Requirement, wh where.Expression, fields ...sql.NamedArg) (int64, error) {
	query, args := updateFieldsSQL(tbl.Name().String(), tbl.Dialect().Quoter(), wh, fields...)
	return Exec(tbl, req, query, args...)
}

func updateFieldsSQL(tblName string, q quote.Quoter, wh where.Expression, fields ...sql.NamedArg) (string, []interface{}) {
	list := sqlapi.NamedArgList(fields)
	assignments := strings.Join(list.Assignments(q, 1), ", ")
	whs, wargs := where.Where(wh, q)
	query := fmt.Sprintf("UPDATE %s SET %s%s", q.Quote(tblName), assignments, whs)
	args := append(list.Values(), wargs...)
	return query, args
}

// DeleteByColumn deletes rows from the table, given some values and the name of the column they belong to.
// The list of values can be arbitrarily long.
func DeleteByColumn(tbl sqlapi.Table, req require.Requirement, column string, v ...interface{}) (int64, error) {
	const batch = 1000 // limited by Oracle DB
	const qt = "DELETE FROM %s WHERE %s IN (%s)"
	qName := tbl.Dialect().Quoter().Quote(tbl.Name().String())

	if req == require.All {
		req = require.Exactly(len(v))
	}

	var count, n int64
	var err error
	var max = batch
	if len(v) < batch {
		max = len(v)
	}
	d := tbl.Dialect()
	col := d.Quoter().Quote(column)
	args := make([]interface{}, max)

	if len(v) > batch {
		pl := d.Placeholders(batch)
		query := fmt.Sprintf(qt, qName, col, pl)

		for len(v) > batch {
			for i := 0; i < batch; i++ {
				args[i] = v[i]
			}

			n, err = doExec(tbl, query, args...)
			count += n
			if err != nil {
				return count, err
			}

			v = v[batch:]
		}
	}

	if len(v) > 0 {
		pl := d.Placeholders(len(v))
		query := fmt.Sprintf(qt, qName, col, pl)

		for i := 0; i < len(v); i++ {
			args[i] = v[i]
		}

		n, err = doExec(tbl, query, args...)
		count += n
	}

	return count, tbl.Logger().LogIfError(tbl.Ctx(), require.ChainErrorIfExecNotSatisfiedBy(err, req, n))
}

// GetIntIntIndex reads two integer columns from a specified database table and returns an index built from them.
func GetIntIntIndex(tbl sqlapi.Table, q quote.Quoter, keyColumn, valColumn string, wh where.Expression) (map[int64]int64, error) {
	whs, args := where.Where(wh)
	query := fmt.Sprintf("SELECT %s, %s FROM %s%s", q.Quote(keyColumn), q.Quote(valColumn), q.Quote(tbl.Name().String()), whs)
	rows, err := Query(tbl, query, args...)
	if err != nil {
		return nil, tbl.Logger().LogError(tbl.Ctx(), err)
	}
	defer rows.Close()

	index := make(map[int64]int64)
	for rows.Next() {
		var k, v int64
		err = rows.Scan(&k, &v)
		if err != nil {
			return nil, tbl.Logger().LogError(tbl.Ctx(), err)
		}
		index[k] = v
	}
	return index, nil
}

// GetStringIntIndex reads a string column and an integer column from a specified database table and returns an index built from them.
func GetStringIntIndex(tbl sqlapi.Table, q quote.Quoter, keyColumn, valColumn string, wh where.Expression) (map[string]int64, error) {
	whs, args := where.Where(wh)
	query := fmt.Sprintf("SELECT %s, %s FROM %s%s", q.Quote(keyColumn), q.Quote(valColumn), q.Quote(tbl.Name().String()), whs)
	rows, err := Query(tbl, query, args...)
	if err != nil {
		return nil, tbl.Logger().LogError(tbl.Ctx(), err)
	}
	defer rows.Close()

	index := make(map[string]int64)
	for rows.Next() {
		var k string
		var v int64
		err = rows.Scan(&k, &v)
		if err != nil {
			return nil, tbl.Logger().LogError(tbl.Ctx(), err)
		}
		index[k] = v
	}
	return index, nil
}

// GetIntStringIndex reads an integer column and a string column from a specified database table and returns an index built from them.
func GetIntStringIndex(tbl sqlapi.Table, q quote.Quoter, keyColumn, valColumn string, wh where.Expression) (map[int64]string, error) {
	whs, args := where.Where(wh)
	query := fmt.Sprintf("SELECT %s, %s FROM %s%s", q.Quote(keyColumn), q.Quote(valColumn), q.Quote(tbl.Name().String()), whs)
	rows, err := Query(tbl, query, args...)
	if err != nil {
		return nil, tbl.Logger().LogError(tbl.Ctx(), err)
	}
	defer rows.Close()

	index := make(map[int64]string)
	for rows.Next() {
		var k int64
		var v string
		err = rows.Scan(&k, &v)
		if err != nil {
			return nil, tbl.Logger().LogError(tbl.Ctx(), err)
		}
		index[k] = v
	}
	return index, nil
}
