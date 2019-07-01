
// Slice{{.Type.U}}List requests a columnar slice of {{.Type}}s from a specified column.
func Slice{{.Type.U}}List(tbl sqlapi.Table, req require.Requirement, sqlname string, wh where.Expression, qc where.QueryConstraint) ([]{{.Type}}, error) {
	query, args := sliceSql(tbl, sqlname, wh, qc)
	rows, err := tbl.Query(query, args...)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	return doScan{{.Type.U}}List(req, rows, tbl.Database().Logger().LogIfError)
}

// doScan{{.Type.U}}List processes result rows to extract a list of {{.Type}}s.
// The result set should have been produced via a SELECT statement on just one column.
func doScan{{.Type.U}}List(req require.Requirement, rows sqlapi.SqlRows, qLog func(error) error) ([]{{.Type}}, error) {
	var v {{.Type}}
	list := make([]{{.Type}}, 0, 10)

	for rows.Next() {
		err := rows.Scan(&v)
		if err == sql.ErrNoRows {
			return list, qLog(require.ErrorIfQueryNotSatisfiedBy(req, int64(len(list))))
		} else {
			list = append(list, v)
		}
	}
	return list, qLog(require.ChainErrorIfQueryNotSatisfiedBy(rows.Err(), req, int64(len(list))))
}

// Slice{{.Type.U}}PtrList requests a columnar slice of {{.Type}}s from a specified nullable column.
func Slice{{.Type.U}}PtrList(tbl sqlapi.Table, req require.Requirement, sqlname string, wh where.Expression, qc where.QueryConstraint) ([]{{.Type}}, error) {
	query, args := sliceSql(tbl, sqlname, wh, qc)
	rows, err := tbl.Query(query, args...)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	return doScan{{.Type.U}}PtrList(req, rows, tbl.Database().Logger().LogIfError)
}

// doScan{{.Type.U}}PtrList processes result rows to extract a list of {{.Type}}s.
// The result set should have been produced via a SELECT statement on just one column.
func doScan{{.Type.U}}PtrList(req require.Requirement, rows sqlapi.SqlRows, qLog func(error) error) ([]{{.Type}}, error) {
	var v sql.Null{{.NT}}
	list := make([]{{.Type}}, 0, 10)

	for rows.Next() {
		err := rows.Scan(&v)
		if err == sql.ErrNoRows {
			return list, qLog(require.ErrorIfQueryNotSatisfiedBy(req, int64(len(list))))
		} else if v.Valid {
			list = append(list, {{.Type}}(v.{{.NT}}))
		}
	}
	return list, qLog(require.ChainErrorIfQueryNotSatisfiedBy(rows.Err(), req, int64(len(list))))
}
