package support

import (
	"context"
	"database/sql"
	"github.com/rickb777/sqlapi/pgxapi"
	"github.com/rickb777/sqlapi/require"
	"github.com/rickb777/where"
)

//-------------------------------------------------------------------------------------------------
// string

// SliceStringList requests a columnar slice of strings from a specified column.
//
// If the context ctx is nil, it defaults to context.Background().
func SliceStringList(ctx context.Context, tbl pgxapi.Table, req require.Requirement, sqlname string, wh where.Expression, qc where.QueryConstraint) ([]string, error) {
	query, args := sliceSql(tbl, sqlname, wh, qc)
	rows, err := Query(ctx, tbl, query, args...)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	return doScanStringList(req, rows, tbl.Database().Logger().LogIfError)
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
//
// If the context ctx is nil, it defaults to context.Background().
func SliceStringPtrList(ctx context.Context, tbl pgxapi.Table, req require.Requirement, sqlname string, wh where.Expression, qc where.QueryConstraint) ([]string, error) {
	query, args := sliceSql(tbl, sqlname, wh, qc)
	rows, err := Query(ctx, tbl, query, args...)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	return doScanStringPtrList(req, rows, tbl.Database().Logger().LogIfError)
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

// StringAsInterfaceSlice adapts a slice of string to []interface{}.
func StringAsInterfaceSlice(values []string) []interface{} {
	ii := make([]interface{}, len(values))
	for i, v := range values {
		ii[i] = v
	}
	return ii
}

//-------------------------------------------------------------------------------------------------
// int

// SliceIntList requests a columnar slice of ints from a specified column.
//
// If the context ctx is nil, it defaults to context.Background().
func SliceIntList(ctx context.Context, tbl pgxapi.Table, req require.Requirement, sqlname string, wh where.Expression, qc where.QueryConstraint) ([]int, error) {
	query, args := sliceSql(tbl, sqlname, wh, qc)
	rows, err := Query(ctx, tbl, query, args...)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	return doScanIntList(req, rows, tbl.Database().Logger().LogIfError)
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
//
// If the context ctx is nil, it defaults to context.Background().
func SliceIntPtrList(ctx context.Context, tbl pgxapi.Table, req require.Requirement, sqlname string, wh where.Expression, qc where.QueryConstraint) ([]int, error) {
	query, args := sliceSql(tbl, sqlname, wh, qc)
	rows, err := Query(ctx, tbl, query, args...)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	return doScanIntPtrList(req, rows, tbl.Database().Logger().LogIfError)
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

// IntAsInterfaceSlice adapts a slice of int to []interface{}.
func IntAsInterfaceSlice(values []int) []interface{} {
	ii := make([]interface{}, len(values))
	for i, v := range values {
		ii[i] = v
	}
	return ii
}

//-------------------------------------------------------------------------------------------------
// int64

// SliceInt64List requests a columnar slice of int64s from a specified column.
//
// If the context ctx is nil, it defaults to context.Background().
func SliceInt64List(ctx context.Context, tbl pgxapi.Table, req require.Requirement, sqlname string, wh where.Expression, qc where.QueryConstraint) ([]int64, error) {
	query, args := sliceSql(tbl, sqlname, wh, qc)
	rows, err := Query(ctx, tbl, query, args...)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	return doScanInt64List(req, rows, tbl.Database().Logger().LogIfError)
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
//
// If the context ctx is nil, it defaults to context.Background().
func SliceInt64PtrList(ctx context.Context, tbl pgxapi.Table, req require.Requirement, sqlname string, wh where.Expression, qc where.QueryConstraint) ([]int64, error) {
	query, args := sliceSql(tbl, sqlname, wh, qc)
	rows, err := Query(ctx, tbl, query, args...)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	return doScanInt64PtrList(req, rows, tbl.Database().Logger().LogIfError)
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

// Int64AsInterfaceSlice adapts a slice of int64 to []interface{}.
func Int64AsInterfaceSlice(values []int64) []interface{} {
	ii := make([]interface{}, len(values))
	for i, v := range values {
		ii[i] = v
	}
	return ii
}

//-------------------------------------------------------------------------------------------------
// int32

// SliceInt32List requests a columnar slice of int32s from a specified column.
//
// If the context ctx is nil, it defaults to context.Background().
func SliceInt32List(ctx context.Context, tbl pgxapi.Table, req require.Requirement, sqlname string, wh where.Expression, qc where.QueryConstraint) ([]int32, error) {
	query, args := sliceSql(tbl, sqlname, wh, qc)
	rows, err := Query(ctx, tbl, query, args...)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	return doScanInt32List(req, rows, tbl.Database().Logger().LogIfError)
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
//
// If the context ctx is nil, it defaults to context.Background().
func SliceInt32PtrList(ctx context.Context, tbl pgxapi.Table, req require.Requirement, sqlname string, wh where.Expression, qc where.QueryConstraint) ([]int32, error) {
	query, args := sliceSql(tbl, sqlname, wh, qc)
	rows, err := Query(ctx, tbl, query, args...)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	return doScanInt32PtrList(req, rows, tbl.Database().Logger().LogIfError)
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

// Int32AsInterfaceSlice adapts a slice of int32 to []interface{}.
func Int32AsInterfaceSlice(values []int32) []interface{} {
	ii := make([]interface{}, len(values))
	for i, v := range values {
		ii[i] = v
	}
	return ii
}

//-------------------------------------------------------------------------------------------------
// int16

// SliceInt16List requests a columnar slice of int16s from a specified column.
//
// If the context ctx is nil, it defaults to context.Background().
func SliceInt16List(ctx context.Context, tbl pgxapi.Table, req require.Requirement, sqlname string, wh where.Expression, qc where.QueryConstraint) ([]int16, error) {
	query, args := sliceSql(tbl, sqlname, wh, qc)
	rows, err := Query(ctx, tbl, query, args...)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	return doScanInt16List(req, rows, tbl.Database().Logger().LogIfError)
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
//
// If the context ctx is nil, it defaults to context.Background().
func SliceInt16PtrList(ctx context.Context, tbl pgxapi.Table, req require.Requirement, sqlname string, wh where.Expression, qc where.QueryConstraint) ([]int16, error) {
	query, args := sliceSql(tbl, sqlname, wh, qc)
	rows, err := Query(ctx, tbl, query, args...)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	return doScanInt16PtrList(req, rows, tbl.Database().Logger().LogIfError)
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

// Int16AsInterfaceSlice adapts a slice of int16 to []interface{}.
func Int16AsInterfaceSlice(values []int16) []interface{} {
	ii := make([]interface{}, len(values))
	for i, v := range values {
		ii[i] = v
	}
	return ii
}

//-------------------------------------------------------------------------------------------------
// int8

// SliceInt8List requests a columnar slice of int8s from a specified column.
//
// If the context ctx is nil, it defaults to context.Background().
func SliceInt8List(ctx context.Context, tbl pgxapi.Table, req require.Requirement, sqlname string, wh where.Expression, qc where.QueryConstraint) ([]int8, error) {
	query, args := sliceSql(tbl, sqlname, wh, qc)
	rows, err := Query(ctx, tbl, query, args...)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	return doScanInt8List(req, rows, tbl.Database().Logger().LogIfError)
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
//
// If the context ctx is nil, it defaults to context.Background().
func SliceInt8PtrList(ctx context.Context, tbl pgxapi.Table, req require.Requirement, sqlname string, wh where.Expression, qc where.QueryConstraint) ([]int8, error) {
	query, args := sliceSql(tbl, sqlname, wh, qc)
	rows, err := Query(ctx, tbl, query, args...)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	return doScanInt8PtrList(req, rows, tbl.Database().Logger().LogIfError)
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

// Int8AsInterfaceSlice adapts a slice of int8 to []interface{}.
func Int8AsInterfaceSlice(values []int8) []interface{} {
	ii := make([]interface{}, len(values))
	for i, v := range values {
		ii[i] = v
	}
	return ii
}

//-------------------------------------------------------------------------------------------------
// uint

// SliceUintList requests a columnar slice of uints from a specified column.
//
// If the context ctx is nil, it defaults to context.Background().
func SliceUintList(ctx context.Context, tbl pgxapi.Table, req require.Requirement, sqlname string, wh where.Expression, qc where.QueryConstraint) ([]uint, error) {
	query, args := sliceSql(tbl, sqlname, wh, qc)
	rows, err := Query(ctx, tbl, query, args...)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	return doScanUintList(req, rows, tbl.Database().Logger().LogIfError)
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
//
// If the context ctx is nil, it defaults to context.Background().
func SliceUintPtrList(ctx context.Context, tbl pgxapi.Table, req require.Requirement, sqlname string, wh where.Expression, qc where.QueryConstraint) ([]uint, error) {
	query, args := sliceSql(tbl, sqlname, wh, qc)
	rows, err := Query(ctx, tbl, query, args...)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	return doScanUintPtrList(req, rows, tbl.Database().Logger().LogIfError)
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

// UintAsInterfaceSlice adapts a slice of uint to []interface{}.
func UintAsInterfaceSlice(values []uint) []interface{} {
	ii := make([]interface{}, len(values))
	for i, v := range values {
		ii[i] = v
	}
	return ii
}

//-------------------------------------------------------------------------------------------------
// uint64

// SliceUint64List requests a columnar slice of uint64s from a specified column.
//
// If the context ctx is nil, it defaults to context.Background().
func SliceUint64List(ctx context.Context, tbl pgxapi.Table, req require.Requirement, sqlname string, wh where.Expression, qc where.QueryConstraint) ([]uint64, error) {
	query, args := sliceSql(tbl, sqlname, wh, qc)
	rows, err := Query(ctx, tbl, query, args...)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	return doScanUint64List(req, rows, tbl.Database().Logger().LogIfError)
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
//
// If the context ctx is nil, it defaults to context.Background().
func SliceUint64PtrList(ctx context.Context, tbl pgxapi.Table, req require.Requirement, sqlname string, wh where.Expression, qc where.QueryConstraint) ([]uint64, error) {
	query, args := sliceSql(tbl, sqlname, wh, qc)
	rows, err := Query(ctx, tbl, query, args...)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	return doScanUint64PtrList(req, rows, tbl.Database().Logger().LogIfError)
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

// Uint64AsInterfaceSlice adapts a slice of uint64 to []interface{}.
func Uint64AsInterfaceSlice(values []uint64) []interface{} {
	ii := make([]interface{}, len(values))
	for i, v := range values {
		ii[i] = v
	}
	return ii
}

//-------------------------------------------------------------------------------------------------
// uint32

// SliceUint32List requests a columnar slice of uint32s from a specified column.
//
// If the context ctx is nil, it defaults to context.Background().
func SliceUint32List(ctx context.Context, tbl pgxapi.Table, req require.Requirement, sqlname string, wh where.Expression, qc where.QueryConstraint) ([]uint32, error) {
	query, args := sliceSql(tbl, sqlname, wh, qc)
	rows, err := Query(ctx, tbl, query, args...)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	return doScanUint32List(req, rows, tbl.Database().Logger().LogIfError)
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
//
// If the context ctx is nil, it defaults to context.Background().
func SliceUint32PtrList(ctx context.Context, tbl pgxapi.Table, req require.Requirement, sqlname string, wh where.Expression, qc where.QueryConstraint) ([]uint32, error) {
	query, args := sliceSql(tbl, sqlname, wh, qc)
	rows, err := Query(ctx, tbl, query, args...)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	return doScanUint32PtrList(req, rows, tbl.Database().Logger().LogIfError)
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

// Uint32AsInterfaceSlice adapts a slice of uint32 to []interface{}.
func Uint32AsInterfaceSlice(values []uint32) []interface{} {
	ii := make([]interface{}, len(values))
	for i, v := range values {
		ii[i] = v
	}
	return ii
}

//-------------------------------------------------------------------------------------------------
// uint16

// SliceUint16List requests a columnar slice of uint16s from a specified column.
//
// If the context ctx is nil, it defaults to context.Background().
func SliceUint16List(ctx context.Context, tbl pgxapi.Table, req require.Requirement, sqlname string, wh where.Expression, qc where.QueryConstraint) ([]uint16, error) {
	query, args := sliceSql(tbl, sqlname, wh, qc)
	rows, err := Query(ctx, tbl, query, args...)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	return doScanUint16List(req, rows, tbl.Database().Logger().LogIfError)
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
//
// If the context ctx is nil, it defaults to context.Background().
func SliceUint16PtrList(ctx context.Context, tbl pgxapi.Table, req require.Requirement, sqlname string, wh where.Expression, qc where.QueryConstraint) ([]uint16, error) {
	query, args := sliceSql(tbl, sqlname, wh, qc)
	rows, err := Query(ctx, tbl, query, args...)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	return doScanUint16PtrList(req, rows, tbl.Database().Logger().LogIfError)
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

// Uint16AsInterfaceSlice adapts a slice of uint16 to []interface{}.
func Uint16AsInterfaceSlice(values []uint16) []interface{} {
	ii := make([]interface{}, len(values))
	for i, v := range values {
		ii[i] = v
	}
	return ii
}

//-------------------------------------------------------------------------------------------------
// uint8

// SliceUint8List requests a columnar slice of uint8s from a specified column.
//
// If the context ctx is nil, it defaults to context.Background().
func SliceUint8List(ctx context.Context, tbl pgxapi.Table, req require.Requirement, sqlname string, wh where.Expression, qc where.QueryConstraint) ([]uint8, error) {
	query, args := sliceSql(tbl, sqlname, wh, qc)
	rows, err := Query(ctx, tbl, query, args...)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	return doScanUint8List(req, rows, tbl.Database().Logger().LogIfError)
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
//
// If the context ctx is nil, it defaults to context.Background().
func SliceUint8PtrList(ctx context.Context, tbl pgxapi.Table, req require.Requirement, sqlname string, wh where.Expression, qc where.QueryConstraint) ([]uint8, error) {
	query, args := sliceSql(tbl, sqlname, wh, qc)
	rows, err := Query(ctx, tbl, query, args...)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	return doScanUint8PtrList(req, rows, tbl.Database().Logger().LogIfError)
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

// Uint8AsInterfaceSlice adapts a slice of uint8 to []interface{}.
func Uint8AsInterfaceSlice(values []uint8) []interface{} {
	ii := make([]interface{}, len(values))
	for i, v := range values {
		ii[i] = v
	}
	return ii
}

//-------------------------------------------------------------------------------------------------
// float64

// SliceFloat64List requests a columnar slice of float64s from a specified column.
//
// If the context ctx is nil, it defaults to context.Background().
func SliceFloat64List(ctx context.Context, tbl pgxapi.Table, req require.Requirement, sqlname string, wh where.Expression, qc where.QueryConstraint) ([]float64, error) {
	query, args := sliceSql(tbl, sqlname, wh, qc)
	rows, err := Query(ctx, tbl, query, args...)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	return doScanFloat64List(req, rows, tbl.Database().Logger().LogIfError)
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
//
// If the context ctx is nil, it defaults to context.Background().
func SliceFloat64PtrList(ctx context.Context, tbl pgxapi.Table, req require.Requirement, sqlname string, wh where.Expression, qc where.QueryConstraint) ([]float64, error) {
	query, args := sliceSql(tbl, sqlname, wh, qc)
	rows, err := Query(ctx, tbl, query, args...)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	return doScanFloat64PtrList(req, rows, tbl.Database().Logger().LogIfError)
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

// Float64AsInterfaceSlice adapts a slice of float64 to []interface{}.
func Float64AsInterfaceSlice(values []float64) []interface{} {
	ii := make([]interface{}, len(values))
	for i, v := range values {
		ii[i] = v
	}
	return ii
}

//-------------------------------------------------------------------------------------------------
// float32

// SliceFloat32List requests a columnar slice of float32s from a specified column.
//
// If the context ctx is nil, it defaults to context.Background().
func SliceFloat32List(ctx context.Context, tbl pgxapi.Table, req require.Requirement, sqlname string, wh where.Expression, qc where.QueryConstraint) ([]float32, error) {
	query, args := sliceSql(tbl, sqlname, wh, qc)
	rows, err := Query(ctx, tbl, query, args...)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	return doScanFloat32List(req, rows, tbl.Database().Logger().LogIfError)
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
//
// If the context ctx is nil, it defaults to context.Background().
func SliceFloat32PtrList(ctx context.Context, tbl pgxapi.Table, req require.Requirement, sqlname string, wh where.Expression, qc where.QueryConstraint) ([]float32, error) {
	query, args := sliceSql(tbl, sqlname, wh, qc)
	rows, err := Query(ctx, tbl, query, args...)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	return doScanFloat32PtrList(req, rows, tbl.Database().Logger().LogIfError)
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

// Float32AsInterfaceSlice adapts a slice of float32 to []interface{}.
func Float32AsInterfaceSlice(values []float32) []interface{} {
	ii := make([]interface{}, len(values))
	for i, v := range values {
		ii[i] = v
	}
	return ii
}
