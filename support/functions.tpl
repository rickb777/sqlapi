
//-------------------------------------------------------------------------------------------------
// {{.Type}}

// Slice{{.Type.U}}List requests a columnar slice of {{.Type}}s from a specified column.
func Slice{{.Type.U}}List(tbl {{.SqlApi}}.Table, req require.Requirement, sqlname string, wh where.Expression, qc where.QueryConstraint) ([]{{.Type}}, error) {
	query, args := sliceSql(tbl, sqlname, wh, qc)
	rows, err := Query(tbl, query, args...)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var v {{.Type}}
	list := make([]{{.Type}}, 0, 10)

	for rows.Next() {
		err := rows.Scan(&v)
		if err == sql.ErrNoRows {
			return list, tbl.Logger().LogIfError(require.ErrorIfQueryNotSatisfiedBy(req, int64(len(list))))
		} else {
			list = append(list, v)
		}
	}

	return list, tbl.Logger().LogIfError(require.ChainErrorIfQueryNotSatisfiedBy(rows.Err(), req, int64(len(list))))
}

// Slice{{.Type.U}}PtrList requests a columnar slice of {{.Type}}s from a specified nullable column.
func Slice{{.Type.U}}PtrList(tbl {{.SqlApi}}.Table, req require.Requirement, sqlname string, wh where.Expression, qc where.QueryConstraint) ([]{{.Type}}, error) {
	query, args := sliceSql(tbl, sqlname, wh, qc)
	rows, err := Query(tbl, query, args...)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var v sql.Null{{.NT}}
	list := make([]{{.Type}}, 0, 10)

	for rows.Next() {
		err := rows.Scan(&v)
		if err == sql.ErrNoRows {
			return list, tbl.Logger().LogIfError(require.ErrorIfQueryNotSatisfiedBy(req, int64(len(list))))
		} else if v.Valid {
			list = append(list, {{.Type}}(v.{{.NT}}))
		}
	}

	return list, tbl.Logger().LogIfError(require.ChainErrorIfQueryNotSatisfiedBy(rows.Err(), req, int64(len(list))))
}

// {{.Type.U}}AsInterfaceSlice adapts a slice of {{.Type}} to []interface{}.
func {{.Type.U}}AsInterfaceSlice(values []{{.Type}}) []interface{} {
	ii := make([]interface{}, len(values))
	for i, v := range values {
		ii[i] = v
	}
	return ii
}
