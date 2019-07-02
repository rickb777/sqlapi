package support

import (
	"database/sql"
	"github.com/rickb777/sqlapi/pgxapi"
	"github.com/rickb777/sqlapi/require"
	"github.com/rickb777/where"
)

// SliceStringList requests a columnar slice of strings from a specified column.
func SliceStringList(tbl pgxapi.Table, req require.Requirement, sqlname string, wh where.Expression, qc where.QueryConstraint) ([]string, error) {
	query, args := sliceSql(tbl, sqlname, wh, qc)
	rows, err := Query(tbl, query, args...)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	return doScanStringList(req, rows, tbl.Database().LogIfError)
}

// doScanStringList processes result rows to extract a list of strings.
// The result set should have been produced via a SELECT statement on just one column.
func doScanStringList(req require.Requirement, rows pgxapi.SqlRows, qLog func(error) error) ([]string, error) {
	var v string
	list := make([]string, 0, 10)

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

// SliceStringPtrList requests a columnar slice of strings from a specified nullable column.
func SliceStringPtrList(tbl pgxapi.Table, req require.Requirement, sqlname string, wh where.Expression, qc where.QueryConstraint) ([]string, error) {
	query, args := sliceSql(tbl, sqlname, wh, qc)
	rows, err := Query(tbl, query, args...)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	return doScanStringPtrList(req, rows, tbl.Database().LogIfError)
}

// doScanStringPtrList processes result rows to extract a list of strings.
// The result set should have been produced via a SELECT statement on just one column.
func doScanStringPtrList(req require.Requirement, rows pgxapi.SqlRows, qLog func(error) error) ([]string, error) {
	var v sql.NullString
	list := make([]string, 0, 10)

	for rows.Next() {
		err := rows.Scan(&v)
		if err == sql.ErrNoRows {
			return list, qLog(require.ErrorIfQueryNotSatisfiedBy(req, int64(len(list))))
		} else if v.Valid {
			list = append(list, string(v.String))
		}
	}
	return list, qLog(require.ChainErrorIfQueryNotSatisfiedBy(rows.Err(), req, int64(len(list))))
}

// SliceIntList requests a columnar slice of ints from a specified column.
func SliceIntList(tbl pgxapi.Table, req require.Requirement, sqlname string, wh where.Expression, qc where.QueryConstraint) ([]int, error) {
	query, args := sliceSql(tbl, sqlname, wh, qc)
	rows, err := Query(tbl, query, args...)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	return doScanIntList(req, rows, tbl.Database().LogIfError)
}

// doScanIntList processes result rows to extract a list of ints.
// The result set should have been produced via a SELECT statement on just one column.
func doScanIntList(req require.Requirement, rows pgxapi.SqlRows, qLog func(error) error) ([]int, error) {
	var v int
	list := make([]int, 0, 10)

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

// SliceIntPtrList requests a columnar slice of ints from a specified nullable column.
func SliceIntPtrList(tbl pgxapi.Table, req require.Requirement, sqlname string, wh where.Expression, qc where.QueryConstraint) ([]int, error) {
	query, args := sliceSql(tbl, sqlname, wh, qc)
	rows, err := Query(tbl, query, args...)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	return doScanIntPtrList(req, rows, tbl.Database().LogIfError)
}

// doScanIntPtrList processes result rows to extract a list of ints.
// The result set should have been produced via a SELECT statement on just one column.
func doScanIntPtrList(req require.Requirement, rows pgxapi.SqlRows, qLog func(error) error) ([]int, error) {
	var v sql.NullInt64
	list := make([]int, 0, 10)

	for rows.Next() {
		err := rows.Scan(&v)
		if err == sql.ErrNoRows {
			return list, qLog(require.ErrorIfQueryNotSatisfiedBy(req, int64(len(list))))
		} else if v.Valid {
			list = append(list, int(v.Int64))
		}
	}
	return list, qLog(require.ChainErrorIfQueryNotSatisfiedBy(rows.Err(), req, int64(len(list))))
}

// SliceInt64List requests a columnar slice of int64s from a specified column.
func SliceInt64List(tbl pgxapi.Table, req require.Requirement, sqlname string, wh where.Expression, qc where.QueryConstraint) ([]int64, error) {
	query, args := sliceSql(tbl, sqlname, wh, qc)
	rows, err := Query(tbl, query, args...)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	return doScanInt64List(req, rows, tbl.Database().LogIfError)
}

// doScanInt64List processes result rows to extract a list of int64s.
// The result set should have been produced via a SELECT statement on just one column.
func doScanInt64List(req require.Requirement, rows pgxapi.SqlRows, qLog func(error) error) ([]int64, error) {
	var v int64
	list := make([]int64, 0, 10)

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

// SliceInt64PtrList requests a columnar slice of int64s from a specified nullable column.
func SliceInt64PtrList(tbl pgxapi.Table, req require.Requirement, sqlname string, wh where.Expression, qc where.QueryConstraint) ([]int64, error) {
	query, args := sliceSql(tbl, sqlname, wh, qc)
	rows, err := Query(tbl, query, args...)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	return doScanInt64PtrList(req, rows, tbl.Database().LogIfError)
}

// doScanInt64PtrList processes result rows to extract a list of int64s.
// The result set should have been produced via a SELECT statement on just one column.
func doScanInt64PtrList(req require.Requirement, rows pgxapi.SqlRows, qLog func(error) error) ([]int64, error) {
	var v sql.NullInt64
	list := make([]int64, 0, 10)

	for rows.Next() {
		err := rows.Scan(&v)
		if err == sql.ErrNoRows {
			return list, qLog(require.ErrorIfQueryNotSatisfiedBy(req, int64(len(list))))
		} else if v.Valid {
			list = append(list, int64(v.Int64))
		}
	}
	return list, qLog(require.ChainErrorIfQueryNotSatisfiedBy(rows.Err(), req, int64(len(list))))
}

// SliceInt32List requests a columnar slice of int32s from a specified column.
func SliceInt32List(tbl pgxapi.Table, req require.Requirement, sqlname string, wh where.Expression, qc where.QueryConstraint) ([]int32, error) {
	query, args := sliceSql(tbl, sqlname, wh, qc)
	rows, err := Query(tbl, query, args...)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	return doScanInt32List(req, rows, tbl.Database().LogIfError)
}

// doScanInt32List processes result rows to extract a list of int32s.
// The result set should have been produced via a SELECT statement on just one column.
func doScanInt32List(req require.Requirement, rows pgxapi.SqlRows, qLog func(error) error) ([]int32, error) {
	var v int32
	list := make([]int32, 0, 10)

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

// SliceInt32PtrList requests a columnar slice of int32s from a specified nullable column.
func SliceInt32PtrList(tbl pgxapi.Table, req require.Requirement, sqlname string, wh where.Expression, qc where.QueryConstraint) ([]int32, error) {
	query, args := sliceSql(tbl, sqlname, wh, qc)
	rows, err := Query(tbl, query, args...)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	return doScanInt32PtrList(req, rows, tbl.Database().LogIfError)
}

// doScanInt32PtrList processes result rows to extract a list of int32s.
// The result set should have been produced via a SELECT statement on just one column.
func doScanInt32PtrList(req require.Requirement, rows pgxapi.SqlRows, qLog func(error) error) ([]int32, error) {
	var v sql.NullInt64
	list := make([]int32, 0, 10)

	for rows.Next() {
		err := rows.Scan(&v)
		if err == sql.ErrNoRows {
			return list, qLog(require.ErrorIfQueryNotSatisfiedBy(req, int64(len(list))))
		} else if v.Valid {
			list = append(list, int32(v.Int64))
		}
	}
	return list, qLog(require.ChainErrorIfQueryNotSatisfiedBy(rows.Err(), req, int64(len(list))))
}

// SliceInt16List requests a columnar slice of int16s from a specified column.
func SliceInt16List(tbl pgxapi.Table, req require.Requirement, sqlname string, wh where.Expression, qc where.QueryConstraint) ([]int16, error) {
	query, args := sliceSql(tbl, sqlname, wh, qc)
	rows, err := Query(tbl, query, args...)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	return doScanInt16List(req, rows, tbl.Database().LogIfError)
}

// doScanInt16List processes result rows to extract a list of int16s.
// The result set should have been produced via a SELECT statement on just one column.
func doScanInt16List(req require.Requirement, rows pgxapi.SqlRows, qLog func(error) error) ([]int16, error) {
	var v int16
	list := make([]int16, 0, 10)

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

// SliceInt16PtrList requests a columnar slice of int16s from a specified nullable column.
func SliceInt16PtrList(tbl pgxapi.Table, req require.Requirement, sqlname string, wh where.Expression, qc where.QueryConstraint) ([]int16, error) {
	query, args := sliceSql(tbl, sqlname, wh, qc)
	rows, err := Query(tbl, query, args...)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	return doScanInt16PtrList(req, rows, tbl.Database().LogIfError)
}

// doScanInt16PtrList processes result rows to extract a list of int16s.
// The result set should have been produced via a SELECT statement on just one column.
func doScanInt16PtrList(req require.Requirement, rows pgxapi.SqlRows, qLog func(error) error) ([]int16, error) {
	var v sql.NullInt64
	list := make([]int16, 0, 10)

	for rows.Next() {
		err := rows.Scan(&v)
		if err == sql.ErrNoRows {
			return list, qLog(require.ErrorIfQueryNotSatisfiedBy(req, int64(len(list))))
		} else if v.Valid {
			list = append(list, int16(v.Int64))
		}
	}
	return list, qLog(require.ChainErrorIfQueryNotSatisfiedBy(rows.Err(), req, int64(len(list))))
}

// SliceInt8List requests a columnar slice of int8s from a specified column.
func SliceInt8List(tbl pgxapi.Table, req require.Requirement, sqlname string, wh where.Expression, qc where.QueryConstraint) ([]int8, error) {
	query, args := sliceSql(tbl, sqlname, wh, qc)
	rows, err := Query(tbl, query, args...)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	return doScanInt8List(req, rows, tbl.Database().LogIfError)
}

// doScanInt8List processes result rows to extract a list of int8s.
// The result set should have been produced via a SELECT statement on just one column.
func doScanInt8List(req require.Requirement, rows pgxapi.SqlRows, qLog func(error) error) ([]int8, error) {
	var v int8
	list := make([]int8, 0, 10)

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

// SliceInt8PtrList requests a columnar slice of int8s from a specified nullable column.
func SliceInt8PtrList(tbl pgxapi.Table, req require.Requirement, sqlname string, wh where.Expression, qc where.QueryConstraint) ([]int8, error) {
	query, args := sliceSql(tbl, sqlname, wh, qc)
	rows, err := Query(tbl, query, args...)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	return doScanInt8PtrList(req, rows, tbl.Database().LogIfError)
}

// doScanInt8PtrList processes result rows to extract a list of int8s.
// The result set should have been produced via a SELECT statement on just one column.
func doScanInt8PtrList(req require.Requirement, rows pgxapi.SqlRows, qLog func(error) error) ([]int8, error) {
	var v sql.NullInt64
	list := make([]int8, 0, 10)

	for rows.Next() {
		err := rows.Scan(&v)
		if err == sql.ErrNoRows {
			return list, qLog(require.ErrorIfQueryNotSatisfiedBy(req, int64(len(list))))
		} else if v.Valid {
			list = append(list, int8(v.Int64))
		}
	}
	return list, qLog(require.ChainErrorIfQueryNotSatisfiedBy(rows.Err(), req, int64(len(list))))
}

// SliceUintList requests a columnar slice of uints from a specified column.
func SliceUintList(tbl pgxapi.Table, req require.Requirement, sqlname string, wh where.Expression, qc where.QueryConstraint) ([]uint, error) {
	query, args := sliceSql(tbl, sqlname, wh, qc)
	rows, err := Query(tbl, query, args...)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	return doScanUintList(req, rows, tbl.Database().LogIfError)
}

// doScanUintList processes result rows to extract a list of uints.
// The result set should have been produced via a SELECT statement on just one column.
func doScanUintList(req require.Requirement, rows pgxapi.SqlRows, qLog func(error) error) ([]uint, error) {
	var v uint
	list := make([]uint, 0, 10)

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

// SliceUintPtrList requests a columnar slice of uints from a specified nullable column.
func SliceUintPtrList(tbl pgxapi.Table, req require.Requirement, sqlname string, wh where.Expression, qc where.QueryConstraint) ([]uint, error) {
	query, args := sliceSql(tbl, sqlname, wh, qc)
	rows, err := Query(tbl, query, args...)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	return doScanUintPtrList(req, rows, tbl.Database().LogIfError)
}

// doScanUintPtrList processes result rows to extract a list of uints.
// The result set should have been produced via a SELECT statement on just one column.
func doScanUintPtrList(req require.Requirement, rows pgxapi.SqlRows, qLog func(error) error) ([]uint, error) {
	var v sql.NullInt64
	list := make([]uint, 0, 10)

	for rows.Next() {
		err := rows.Scan(&v)
		if err == sql.ErrNoRows {
			return list, qLog(require.ErrorIfQueryNotSatisfiedBy(req, int64(len(list))))
		} else if v.Valid {
			list = append(list, uint(v.Int64))
		}
	}
	return list, qLog(require.ChainErrorIfQueryNotSatisfiedBy(rows.Err(), req, int64(len(list))))
}

// SliceUint64List requests a columnar slice of uint64s from a specified column.
func SliceUint64List(tbl pgxapi.Table, req require.Requirement, sqlname string, wh where.Expression, qc where.QueryConstraint) ([]uint64, error) {
	query, args := sliceSql(tbl, sqlname, wh, qc)
	rows, err := Query(tbl, query, args...)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	return doScanUint64List(req, rows, tbl.Database().LogIfError)
}

// doScanUint64List processes result rows to extract a list of uint64s.
// The result set should have been produced via a SELECT statement on just one column.
func doScanUint64List(req require.Requirement, rows pgxapi.SqlRows, qLog func(error) error) ([]uint64, error) {
	var v uint64
	list := make([]uint64, 0, 10)

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

// SliceUint64PtrList requests a columnar slice of uint64s from a specified nullable column.
func SliceUint64PtrList(tbl pgxapi.Table, req require.Requirement, sqlname string, wh where.Expression, qc where.QueryConstraint) ([]uint64, error) {
	query, args := sliceSql(tbl, sqlname, wh, qc)
	rows, err := Query(tbl, query, args...)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	return doScanUint64PtrList(req, rows, tbl.Database().LogIfError)
}

// doScanUint64PtrList processes result rows to extract a list of uint64s.
// The result set should have been produced via a SELECT statement on just one column.
func doScanUint64PtrList(req require.Requirement, rows pgxapi.SqlRows, qLog func(error) error) ([]uint64, error) {
	var v sql.NullInt64
	list := make([]uint64, 0, 10)

	for rows.Next() {
		err := rows.Scan(&v)
		if err == sql.ErrNoRows {
			return list, qLog(require.ErrorIfQueryNotSatisfiedBy(req, int64(len(list))))
		} else if v.Valid {
			list = append(list, uint64(v.Int64))
		}
	}
	return list, qLog(require.ChainErrorIfQueryNotSatisfiedBy(rows.Err(), req, int64(len(list))))
}

// SliceUint32List requests a columnar slice of uint32s from a specified column.
func SliceUint32List(tbl pgxapi.Table, req require.Requirement, sqlname string, wh where.Expression, qc where.QueryConstraint) ([]uint32, error) {
	query, args := sliceSql(tbl, sqlname, wh, qc)
	rows, err := Query(tbl, query, args...)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	return doScanUint32List(req, rows, tbl.Database().LogIfError)
}

// doScanUint32List processes result rows to extract a list of uint32s.
// The result set should have been produced via a SELECT statement on just one column.
func doScanUint32List(req require.Requirement, rows pgxapi.SqlRows, qLog func(error) error) ([]uint32, error) {
	var v uint32
	list := make([]uint32, 0, 10)

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

// SliceUint32PtrList requests a columnar slice of uint32s from a specified nullable column.
func SliceUint32PtrList(tbl pgxapi.Table, req require.Requirement, sqlname string, wh where.Expression, qc where.QueryConstraint) ([]uint32, error) {
	query, args := sliceSql(tbl, sqlname, wh, qc)
	rows, err := Query(tbl, query, args...)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	return doScanUint32PtrList(req, rows, tbl.Database().LogIfError)
}

// doScanUint32PtrList processes result rows to extract a list of uint32s.
// The result set should have been produced via a SELECT statement on just one column.
func doScanUint32PtrList(req require.Requirement, rows pgxapi.SqlRows, qLog func(error) error) ([]uint32, error) {
	var v sql.NullInt64
	list := make([]uint32, 0, 10)

	for rows.Next() {
		err := rows.Scan(&v)
		if err == sql.ErrNoRows {
			return list, qLog(require.ErrorIfQueryNotSatisfiedBy(req, int64(len(list))))
		} else if v.Valid {
			list = append(list, uint32(v.Int64))
		}
	}
	return list, qLog(require.ChainErrorIfQueryNotSatisfiedBy(rows.Err(), req, int64(len(list))))
}

// SliceUint16List requests a columnar slice of uint16s from a specified column.
func SliceUint16List(tbl pgxapi.Table, req require.Requirement, sqlname string, wh where.Expression, qc where.QueryConstraint) ([]uint16, error) {
	query, args := sliceSql(tbl, sqlname, wh, qc)
	rows, err := Query(tbl, query, args...)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	return doScanUint16List(req, rows, tbl.Database().LogIfError)
}

// doScanUint16List processes result rows to extract a list of uint16s.
// The result set should have been produced via a SELECT statement on just one column.
func doScanUint16List(req require.Requirement, rows pgxapi.SqlRows, qLog func(error) error) ([]uint16, error) {
	var v uint16
	list := make([]uint16, 0, 10)

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

// SliceUint16PtrList requests a columnar slice of uint16s from a specified nullable column.
func SliceUint16PtrList(tbl pgxapi.Table, req require.Requirement, sqlname string, wh where.Expression, qc where.QueryConstraint) ([]uint16, error) {
	query, args := sliceSql(tbl, sqlname, wh, qc)
	rows, err := Query(tbl, query, args...)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	return doScanUint16PtrList(req, rows, tbl.Database().LogIfError)
}

// doScanUint16PtrList processes result rows to extract a list of uint16s.
// The result set should have been produced via a SELECT statement on just one column.
func doScanUint16PtrList(req require.Requirement, rows pgxapi.SqlRows, qLog func(error) error) ([]uint16, error) {
	var v sql.NullInt64
	list := make([]uint16, 0, 10)

	for rows.Next() {
		err := rows.Scan(&v)
		if err == sql.ErrNoRows {
			return list, qLog(require.ErrorIfQueryNotSatisfiedBy(req, int64(len(list))))
		} else if v.Valid {
			list = append(list, uint16(v.Int64))
		}
	}
	return list, qLog(require.ChainErrorIfQueryNotSatisfiedBy(rows.Err(), req, int64(len(list))))
}

// SliceUint8List requests a columnar slice of uint8s from a specified column.
func SliceUint8List(tbl pgxapi.Table, req require.Requirement, sqlname string, wh where.Expression, qc where.QueryConstraint) ([]uint8, error) {
	query, args := sliceSql(tbl, sqlname, wh, qc)
	rows, err := Query(tbl, query, args...)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	return doScanUint8List(req, rows, tbl.Database().LogIfError)
}

// doScanUint8List processes result rows to extract a list of uint8s.
// The result set should have been produced via a SELECT statement on just one column.
func doScanUint8List(req require.Requirement, rows pgxapi.SqlRows, qLog func(error) error) ([]uint8, error) {
	var v uint8
	list := make([]uint8, 0, 10)

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

// SliceUint8PtrList requests a columnar slice of uint8s from a specified nullable column.
func SliceUint8PtrList(tbl pgxapi.Table, req require.Requirement, sqlname string, wh where.Expression, qc where.QueryConstraint) ([]uint8, error) {
	query, args := sliceSql(tbl, sqlname, wh, qc)
	rows, err := Query(tbl, query, args...)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	return doScanUint8PtrList(req, rows, tbl.Database().LogIfError)
}

// doScanUint8PtrList processes result rows to extract a list of uint8s.
// The result set should have been produced via a SELECT statement on just one column.
func doScanUint8PtrList(req require.Requirement, rows pgxapi.SqlRows, qLog func(error) error) ([]uint8, error) {
	var v sql.NullInt64
	list := make([]uint8, 0, 10)

	for rows.Next() {
		err := rows.Scan(&v)
		if err == sql.ErrNoRows {
			return list, qLog(require.ErrorIfQueryNotSatisfiedBy(req, int64(len(list))))
		} else if v.Valid {
			list = append(list, uint8(v.Int64))
		}
	}
	return list, qLog(require.ChainErrorIfQueryNotSatisfiedBy(rows.Err(), req, int64(len(list))))
}

// SliceFloat64List requests a columnar slice of float64s from a specified column.
func SliceFloat64List(tbl pgxapi.Table, req require.Requirement, sqlname string, wh where.Expression, qc where.QueryConstraint) ([]float64, error) {
	query, args := sliceSql(tbl, sqlname, wh, qc)
	rows, err := Query(tbl, query, args...)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	return doScanFloat64List(req, rows, tbl.Database().LogIfError)
}

// doScanFloat64List processes result rows to extract a list of float64s.
// The result set should have been produced via a SELECT statement on just one column.
func doScanFloat64List(req require.Requirement, rows pgxapi.SqlRows, qLog func(error) error) ([]float64, error) {
	var v float64
	list := make([]float64, 0, 10)

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

// SliceFloat64PtrList requests a columnar slice of float64s from a specified nullable column.
func SliceFloat64PtrList(tbl pgxapi.Table, req require.Requirement, sqlname string, wh where.Expression, qc where.QueryConstraint) ([]float64, error) {
	query, args := sliceSql(tbl, sqlname, wh, qc)
	rows, err := Query(tbl, query, args...)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	return doScanFloat64PtrList(req, rows, tbl.Database().LogIfError)
}

// doScanFloat64PtrList processes result rows to extract a list of float64s.
// The result set should have been produced via a SELECT statement on just one column.
func doScanFloat64PtrList(req require.Requirement, rows pgxapi.SqlRows, qLog func(error) error) ([]float64, error) {
	var v sql.NullFloat64
	list := make([]float64, 0, 10)

	for rows.Next() {
		err := rows.Scan(&v)
		if err == sql.ErrNoRows {
			return list, qLog(require.ErrorIfQueryNotSatisfiedBy(req, int64(len(list))))
		} else if v.Valid {
			list = append(list, float64(v.Float64))
		}
	}
	return list, qLog(require.ChainErrorIfQueryNotSatisfiedBy(rows.Err(), req, int64(len(list))))
}

// SliceFloat32List requests a columnar slice of float32s from a specified column.
func SliceFloat32List(tbl pgxapi.Table, req require.Requirement, sqlname string, wh where.Expression, qc where.QueryConstraint) ([]float32, error) {
	query, args := sliceSql(tbl, sqlname, wh, qc)
	rows, err := Query(tbl, query, args...)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	return doScanFloat32List(req, rows, tbl.Database().LogIfError)
}

// doScanFloat32List processes result rows to extract a list of float32s.
// The result set should have been produced via a SELECT statement on just one column.
func doScanFloat32List(req require.Requirement, rows pgxapi.SqlRows, qLog func(error) error) ([]float32, error) {
	var v float32
	list := make([]float32, 0, 10)

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

// SliceFloat32PtrList requests a columnar slice of float32s from a specified nullable column.
func SliceFloat32PtrList(tbl pgxapi.Table, req require.Requirement, sqlname string, wh where.Expression, qc where.QueryConstraint) ([]float32, error) {
	query, args := sliceSql(tbl, sqlname, wh, qc)
	rows, err := Query(tbl, query, args...)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	return doScanFloat32PtrList(req, rows, tbl.Database().LogIfError)
}

// doScanFloat32PtrList processes result rows to extract a list of float32s.
// The result set should have been produced via a SELECT statement on just one column.
func doScanFloat32PtrList(req require.Requirement, rows pgxapi.SqlRows, qLog func(error) error) ([]float32, error) {
	var v sql.NullFloat64
	list := make([]float32, 0, 10)

	for rows.Next() {
		err := rows.Scan(&v)
		if err == sql.ErrNoRows {
			return list, qLog(require.ErrorIfQueryNotSatisfiedBy(req, int64(len(list))))
		} else if v.Valid {
			list = append(list, float32(v.Float64))
		}
	}
	return list, qLog(require.ChainErrorIfQueryNotSatisfiedBy(rows.Err(), req, int64(len(list))))
}
