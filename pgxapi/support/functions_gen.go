package support

import (
	"database/sql"
	"github.com/rickb777/sqlapi/pgxapi"
	"github.com/rickb777/sqlapi/require"
	"github.com/rickb777/where"
)

//-------------------------------------------------------------------------------------------------
// string

// SliceStringList requests a columnar slice of strings from a specified column.
func SliceStringList(tbl pgxapi.Table, req require.Requirement, sqlname string, wh where.Expression, qc where.QueryConstraint) ([]string, error) {
	query, args := sliceSql(tbl, sqlname, wh, qc)
	rows, err := Query(tbl, query, args...)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var v string
	list := make([]string, 0, 10)

	for rows.Next() {
		err := rows.Scan(&v)
		if err == sql.ErrNoRows {
			return list, tbl.Logger().LogIfError(tbl.Ctx(), require.ErrorIfQueryNotSatisfiedBy(req, int64(len(list))))
		} else {
			list = append(list, v)
		}
	}

	return list, tbl.Logger().LogIfError(tbl.Ctx(), require.ChainErrorIfQueryNotSatisfiedBy(rows.Err(), req, int64(len(list))))
}

// SliceStringPtrList requests a columnar slice of strings from a specified nullable column.
func SliceStringPtrList(tbl pgxapi.Table, req require.Requirement, sqlname string, wh where.Expression, qc where.QueryConstraint) ([]string, error) {
	query, args := sliceSql(tbl, sqlname, wh, qc)
	rows, err := Query(tbl, query, args...)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var v sql.NullString
	list := make([]string, 0, 10)

	for rows.Next() {
		err := rows.Scan(&v)
		if err == sql.ErrNoRows {
			return list, tbl.Logger().LogIfError(tbl.Ctx(), require.ErrorIfQueryNotSatisfiedBy(req, int64(len(list))))
		} else if v.Valid {
			list = append(list, string(v.String))
		}
	}

	return list, tbl.Logger().LogIfError(tbl.Ctx(), require.ChainErrorIfQueryNotSatisfiedBy(rows.Err(), req, int64(len(list))))
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
func SliceIntList(tbl pgxapi.Table, req require.Requirement, sqlname string, wh where.Expression, qc where.QueryConstraint) ([]int, error) {
	query, args := sliceSql(tbl, sqlname, wh, qc)
	rows, err := Query(tbl, query, args...)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var v int
	list := make([]int, 0, 10)

	for rows.Next() {
		err := rows.Scan(&v)
		if err == sql.ErrNoRows {
			return list, tbl.Logger().LogIfError(tbl.Ctx(), require.ErrorIfQueryNotSatisfiedBy(req, int64(len(list))))
		} else {
			list = append(list, v)
		}
	}

	return list, tbl.Logger().LogIfError(tbl.Ctx(), require.ChainErrorIfQueryNotSatisfiedBy(rows.Err(), req, int64(len(list))))
}

// SliceIntPtrList requests a columnar slice of ints from a specified nullable column.
func SliceIntPtrList(tbl pgxapi.Table, req require.Requirement, sqlname string, wh where.Expression, qc where.QueryConstraint) ([]int, error) {
	query, args := sliceSql(tbl, sqlname, wh, qc)
	rows, err := Query(tbl, query, args...)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var v sql.NullInt64
	list := make([]int, 0, 10)

	for rows.Next() {
		err := rows.Scan(&v)
		if err == sql.ErrNoRows {
			return list, tbl.Logger().LogIfError(tbl.Ctx(), require.ErrorIfQueryNotSatisfiedBy(req, int64(len(list))))
		} else if v.Valid {
			list = append(list, int(v.Int64))
		}
	}

	return list, tbl.Logger().LogIfError(tbl.Ctx(), require.ChainErrorIfQueryNotSatisfiedBy(rows.Err(), req, int64(len(list))))
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
func SliceInt64List(tbl pgxapi.Table, req require.Requirement, sqlname string, wh where.Expression, qc where.QueryConstraint) ([]int64, error) {
	query, args := sliceSql(tbl, sqlname, wh, qc)
	rows, err := Query(tbl, query, args...)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var v int64
	list := make([]int64, 0, 10)

	for rows.Next() {
		err := rows.Scan(&v)
		if err == sql.ErrNoRows {
			return list, tbl.Logger().LogIfError(tbl.Ctx(), require.ErrorIfQueryNotSatisfiedBy(req, int64(len(list))))
		} else {
			list = append(list, v)
		}
	}

	return list, tbl.Logger().LogIfError(tbl.Ctx(), require.ChainErrorIfQueryNotSatisfiedBy(rows.Err(), req, int64(len(list))))
}

// SliceInt64PtrList requests a columnar slice of int64s from a specified nullable column.
func SliceInt64PtrList(tbl pgxapi.Table, req require.Requirement, sqlname string, wh where.Expression, qc where.QueryConstraint) ([]int64, error) {
	query, args := sliceSql(tbl, sqlname, wh, qc)
	rows, err := Query(tbl, query, args...)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var v sql.NullInt64
	list := make([]int64, 0, 10)

	for rows.Next() {
		err := rows.Scan(&v)
		if err == sql.ErrNoRows {
			return list, tbl.Logger().LogIfError(tbl.Ctx(), require.ErrorIfQueryNotSatisfiedBy(req, int64(len(list))))
		} else if v.Valid {
			list = append(list, int64(v.Int64))
		}
	}

	return list, tbl.Logger().LogIfError(tbl.Ctx(), require.ChainErrorIfQueryNotSatisfiedBy(rows.Err(), req, int64(len(list))))
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
func SliceInt32List(tbl pgxapi.Table, req require.Requirement, sqlname string, wh where.Expression, qc where.QueryConstraint) ([]int32, error) {
	query, args := sliceSql(tbl, sqlname, wh, qc)
	rows, err := Query(tbl, query, args...)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var v int32
	list := make([]int32, 0, 10)

	for rows.Next() {
		err := rows.Scan(&v)
		if err == sql.ErrNoRows {
			return list, tbl.Logger().LogIfError(tbl.Ctx(), require.ErrorIfQueryNotSatisfiedBy(req, int64(len(list))))
		} else {
			list = append(list, v)
		}
	}

	return list, tbl.Logger().LogIfError(tbl.Ctx(), require.ChainErrorIfQueryNotSatisfiedBy(rows.Err(), req, int64(len(list))))
}

// SliceInt32PtrList requests a columnar slice of int32s from a specified nullable column.
func SliceInt32PtrList(tbl pgxapi.Table, req require.Requirement, sqlname string, wh where.Expression, qc where.QueryConstraint) ([]int32, error) {
	query, args := sliceSql(tbl, sqlname, wh, qc)
	rows, err := Query(tbl, query, args...)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var v sql.NullInt64
	list := make([]int32, 0, 10)

	for rows.Next() {
		err := rows.Scan(&v)
		if err == sql.ErrNoRows {
			return list, tbl.Logger().LogIfError(tbl.Ctx(), require.ErrorIfQueryNotSatisfiedBy(req, int64(len(list))))
		} else if v.Valid {
			list = append(list, int32(v.Int64))
		}
	}

	return list, tbl.Logger().LogIfError(tbl.Ctx(), require.ChainErrorIfQueryNotSatisfiedBy(rows.Err(), req, int64(len(list))))
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
func SliceInt16List(tbl pgxapi.Table, req require.Requirement, sqlname string, wh where.Expression, qc where.QueryConstraint) ([]int16, error) {
	query, args := sliceSql(tbl, sqlname, wh, qc)
	rows, err := Query(tbl, query, args...)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var v int16
	list := make([]int16, 0, 10)

	for rows.Next() {
		err := rows.Scan(&v)
		if err == sql.ErrNoRows {
			return list, tbl.Logger().LogIfError(tbl.Ctx(), require.ErrorIfQueryNotSatisfiedBy(req, int64(len(list))))
		} else {
			list = append(list, v)
		}
	}

	return list, tbl.Logger().LogIfError(tbl.Ctx(), require.ChainErrorIfQueryNotSatisfiedBy(rows.Err(), req, int64(len(list))))
}

// SliceInt16PtrList requests a columnar slice of int16s from a specified nullable column.
func SliceInt16PtrList(tbl pgxapi.Table, req require.Requirement, sqlname string, wh where.Expression, qc where.QueryConstraint) ([]int16, error) {
	query, args := sliceSql(tbl, sqlname, wh, qc)
	rows, err := Query(tbl, query, args...)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var v sql.NullInt64
	list := make([]int16, 0, 10)

	for rows.Next() {
		err := rows.Scan(&v)
		if err == sql.ErrNoRows {
			return list, tbl.Logger().LogIfError(tbl.Ctx(), require.ErrorIfQueryNotSatisfiedBy(req, int64(len(list))))
		} else if v.Valid {
			list = append(list, int16(v.Int64))
		}
	}

	return list, tbl.Logger().LogIfError(tbl.Ctx(), require.ChainErrorIfQueryNotSatisfiedBy(rows.Err(), req, int64(len(list))))
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
func SliceInt8List(tbl pgxapi.Table, req require.Requirement, sqlname string, wh where.Expression, qc where.QueryConstraint) ([]int8, error) {
	query, args := sliceSql(tbl, sqlname, wh, qc)
	rows, err := Query(tbl, query, args...)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var v int8
	list := make([]int8, 0, 10)

	for rows.Next() {
		err := rows.Scan(&v)
		if err == sql.ErrNoRows {
			return list, tbl.Logger().LogIfError(tbl.Ctx(), require.ErrorIfQueryNotSatisfiedBy(req, int64(len(list))))
		} else {
			list = append(list, v)
		}
	}

	return list, tbl.Logger().LogIfError(tbl.Ctx(), require.ChainErrorIfQueryNotSatisfiedBy(rows.Err(), req, int64(len(list))))
}

// SliceInt8PtrList requests a columnar slice of int8s from a specified nullable column.
func SliceInt8PtrList(tbl pgxapi.Table, req require.Requirement, sqlname string, wh where.Expression, qc where.QueryConstraint) ([]int8, error) {
	query, args := sliceSql(tbl, sqlname, wh, qc)
	rows, err := Query(tbl, query, args...)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var v sql.NullInt64
	list := make([]int8, 0, 10)

	for rows.Next() {
		err := rows.Scan(&v)
		if err == sql.ErrNoRows {
			return list, tbl.Logger().LogIfError(tbl.Ctx(), require.ErrorIfQueryNotSatisfiedBy(req, int64(len(list))))
		} else if v.Valid {
			list = append(list, int8(v.Int64))
		}
	}

	return list, tbl.Logger().LogIfError(tbl.Ctx(), require.ChainErrorIfQueryNotSatisfiedBy(rows.Err(), req, int64(len(list))))
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
func SliceUintList(tbl pgxapi.Table, req require.Requirement, sqlname string, wh where.Expression, qc where.QueryConstraint) ([]uint, error) {
	query, args := sliceSql(tbl, sqlname, wh, qc)
	rows, err := Query(tbl, query, args...)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var v uint
	list := make([]uint, 0, 10)

	for rows.Next() {
		err := rows.Scan(&v)
		if err == sql.ErrNoRows {
			return list, tbl.Logger().LogIfError(tbl.Ctx(), require.ErrorIfQueryNotSatisfiedBy(req, int64(len(list))))
		} else {
			list = append(list, v)
		}
	}

	return list, tbl.Logger().LogIfError(tbl.Ctx(), require.ChainErrorIfQueryNotSatisfiedBy(rows.Err(), req, int64(len(list))))
}

// SliceUintPtrList requests a columnar slice of uints from a specified nullable column.
func SliceUintPtrList(tbl pgxapi.Table, req require.Requirement, sqlname string, wh where.Expression, qc where.QueryConstraint) ([]uint, error) {
	query, args := sliceSql(tbl, sqlname, wh, qc)
	rows, err := Query(tbl, query, args...)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var v sql.NullInt64
	list := make([]uint, 0, 10)

	for rows.Next() {
		err := rows.Scan(&v)
		if err == sql.ErrNoRows {
			return list, tbl.Logger().LogIfError(tbl.Ctx(), require.ErrorIfQueryNotSatisfiedBy(req, int64(len(list))))
		} else if v.Valid {
			list = append(list, uint(v.Int64))
		}
	}

	return list, tbl.Logger().LogIfError(tbl.Ctx(), require.ChainErrorIfQueryNotSatisfiedBy(rows.Err(), req, int64(len(list))))
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
func SliceUint64List(tbl pgxapi.Table, req require.Requirement, sqlname string, wh where.Expression, qc where.QueryConstraint) ([]uint64, error) {
	query, args := sliceSql(tbl, sqlname, wh, qc)
	rows, err := Query(tbl, query, args...)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var v uint64
	list := make([]uint64, 0, 10)

	for rows.Next() {
		err := rows.Scan(&v)
		if err == sql.ErrNoRows {
			return list, tbl.Logger().LogIfError(tbl.Ctx(), require.ErrorIfQueryNotSatisfiedBy(req, int64(len(list))))
		} else {
			list = append(list, v)
		}
	}

	return list, tbl.Logger().LogIfError(tbl.Ctx(), require.ChainErrorIfQueryNotSatisfiedBy(rows.Err(), req, int64(len(list))))
}

// SliceUint64PtrList requests a columnar slice of uint64s from a specified nullable column.
func SliceUint64PtrList(tbl pgxapi.Table, req require.Requirement, sqlname string, wh where.Expression, qc where.QueryConstraint) ([]uint64, error) {
	query, args := sliceSql(tbl, sqlname, wh, qc)
	rows, err := Query(tbl, query, args...)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var v sql.NullInt64
	list := make([]uint64, 0, 10)

	for rows.Next() {
		err := rows.Scan(&v)
		if err == sql.ErrNoRows {
			return list, tbl.Logger().LogIfError(tbl.Ctx(), require.ErrorIfQueryNotSatisfiedBy(req, int64(len(list))))
		} else if v.Valid {
			list = append(list, uint64(v.Int64))
		}
	}

	return list, tbl.Logger().LogIfError(tbl.Ctx(), require.ChainErrorIfQueryNotSatisfiedBy(rows.Err(), req, int64(len(list))))
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
func SliceUint32List(tbl pgxapi.Table, req require.Requirement, sqlname string, wh where.Expression, qc where.QueryConstraint) ([]uint32, error) {
	query, args := sliceSql(tbl, sqlname, wh, qc)
	rows, err := Query(tbl, query, args...)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var v uint32
	list := make([]uint32, 0, 10)

	for rows.Next() {
		err := rows.Scan(&v)
		if err == sql.ErrNoRows {
			return list, tbl.Logger().LogIfError(tbl.Ctx(), require.ErrorIfQueryNotSatisfiedBy(req, int64(len(list))))
		} else {
			list = append(list, v)
		}
	}

	return list, tbl.Logger().LogIfError(tbl.Ctx(), require.ChainErrorIfQueryNotSatisfiedBy(rows.Err(), req, int64(len(list))))
}

// SliceUint32PtrList requests a columnar slice of uint32s from a specified nullable column.
func SliceUint32PtrList(tbl pgxapi.Table, req require.Requirement, sqlname string, wh where.Expression, qc where.QueryConstraint) ([]uint32, error) {
	query, args := sliceSql(tbl, sqlname, wh, qc)
	rows, err := Query(tbl, query, args...)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var v sql.NullInt64
	list := make([]uint32, 0, 10)

	for rows.Next() {
		err := rows.Scan(&v)
		if err == sql.ErrNoRows {
			return list, tbl.Logger().LogIfError(tbl.Ctx(), require.ErrorIfQueryNotSatisfiedBy(req, int64(len(list))))
		} else if v.Valid {
			list = append(list, uint32(v.Int64))
		}
	}

	return list, tbl.Logger().LogIfError(tbl.Ctx(), require.ChainErrorIfQueryNotSatisfiedBy(rows.Err(), req, int64(len(list))))
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
func SliceUint16List(tbl pgxapi.Table, req require.Requirement, sqlname string, wh where.Expression, qc where.QueryConstraint) ([]uint16, error) {
	query, args := sliceSql(tbl, sqlname, wh, qc)
	rows, err := Query(tbl, query, args...)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var v uint16
	list := make([]uint16, 0, 10)

	for rows.Next() {
		err := rows.Scan(&v)
		if err == sql.ErrNoRows {
			return list, tbl.Logger().LogIfError(tbl.Ctx(), require.ErrorIfQueryNotSatisfiedBy(req, int64(len(list))))
		} else {
			list = append(list, v)
		}
	}

	return list, tbl.Logger().LogIfError(tbl.Ctx(), require.ChainErrorIfQueryNotSatisfiedBy(rows.Err(), req, int64(len(list))))
}

// SliceUint16PtrList requests a columnar slice of uint16s from a specified nullable column.
func SliceUint16PtrList(tbl pgxapi.Table, req require.Requirement, sqlname string, wh where.Expression, qc where.QueryConstraint) ([]uint16, error) {
	query, args := sliceSql(tbl, sqlname, wh, qc)
	rows, err := Query(tbl, query, args...)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var v sql.NullInt64
	list := make([]uint16, 0, 10)

	for rows.Next() {
		err := rows.Scan(&v)
		if err == sql.ErrNoRows {
			return list, tbl.Logger().LogIfError(tbl.Ctx(), require.ErrorIfQueryNotSatisfiedBy(req, int64(len(list))))
		} else if v.Valid {
			list = append(list, uint16(v.Int64))
		}
	}

	return list, tbl.Logger().LogIfError(tbl.Ctx(), require.ChainErrorIfQueryNotSatisfiedBy(rows.Err(), req, int64(len(list))))
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
func SliceUint8List(tbl pgxapi.Table, req require.Requirement, sqlname string, wh where.Expression, qc where.QueryConstraint) ([]uint8, error) {
	query, args := sliceSql(tbl, sqlname, wh, qc)
	rows, err := Query(tbl, query, args...)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var v uint8
	list := make([]uint8, 0, 10)

	for rows.Next() {
		err := rows.Scan(&v)
		if err == sql.ErrNoRows {
			return list, tbl.Logger().LogIfError(tbl.Ctx(), require.ErrorIfQueryNotSatisfiedBy(req, int64(len(list))))
		} else {
			list = append(list, v)
		}
	}

	return list, tbl.Logger().LogIfError(tbl.Ctx(), require.ChainErrorIfQueryNotSatisfiedBy(rows.Err(), req, int64(len(list))))
}

// SliceUint8PtrList requests a columnar slice of uint8s from a specified nullable column.
func SliceUint8PtrList(tbl pgxapi.Table, req require.Requirement, sqlname string, wh where.Expression, qc where.QueryConstraint) ([]uint8, error) {
	query, args := sliceSql(tbl, sqlname, wh, qc)
	rows, err := Query(tbl, query, args...)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var v sql.NullInt64
	list := make([]uint8, 0, 10)

	for rows.Next() {
		err := rows.Scan(&v)
		if err == sql.ErrNoRows {
			return list, tbl.Logger().LogIfError(tbl.Ctx(), require.ErrorIfQueryNotSatisfiedBy(req, int64(len(list))))
		} else if v.Valid {
			list = append(list, uint8(v.Int64))
		}
	}

	return list, tbl.Logger().LogIfError(tbl.Ctx(), require.ChainErrorIfQueryNotSatisfiedBy(rows.Err(), req, int64(len(list))))
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
func SliceFloat64List(tbl pgxapi.Table, req require.Requirement, sqlname string, wh where.Expression, qc where.QueryConstraint) ([]float64, error) {
	query, args := sliceSql(tbl, sqlname, wh, qc)
	rows, err := Query(tbl, query, args...)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var v float64
	list := make([]float64, 0, 10)

	for rows.Next() {
		err := rows.Scan(&v)
		if err == sql.ErrNoRows {
			return list, tbl.Logger().LogIfError(tbl.Ctx(), require.ErrorIfQueryNotSatisfiedBy(req, int64(len(list))))
		} else {
			list = append(list, v)
		}
	}

	return list, tbl.Logger().LogIfError(tbl.Ctx(), require.ChainErrorIfQueryNotSatisfiedBy(rows.Err(), req, int64(len(list))))
}

// SliceFloat64PtrList requests a columnar slice of float64s from a specified nullable column.
func SliceFloat64PtrList(tbl pgxapi.Table, req require.Requirement, sqlname string, wh where.Expression, qc where.QueryConstraint) ([]float64, error) {
	query, args := sliceSql(tbl, sqlname, wh, qc)
	rows, err := Query(tbl, query, args...)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var v sql.NullFloat64
	list := make([]float64, 0, 10)

	for rows.Next() {
		err := rows.Scan(&v)
		if err == sql.ErrNoRows {
			return list, tbl.Logger().LogIfError(tbl.Ctx(), require.ErrorIfQueryNotSatisfiedBy(req, int64(len(list))))
		} else if v.Valid {
			list = append(list, float64(v.Float64))
		}
	}

	return list, tbl.Logger().LogIfError(tbl.Ctx(), require.ChainErrorIfQueryNotSatisfiedBy(rows.Err(), req, int64(len(list))))
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
func SliceFloat32List(tbl pgxapi.Table, req require.Requirement, sqlname string, wh where.Expression, qc where.QueryConstraint) ([]float32, error) {
	query, args := sliceSql(tbl, sqlname, wh, qc)
	rows, err := Query(tbl, query, args...)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var v float32
	list := make([]float32, 0, 10)

	for rows.Next() {
		err := rows.Scan(&v)
		if err == sql.ErrNoRows {
			return list, tbl.Logger().LogIfError(tbl.Ctx(), require.ErrorIfQueryNotSatisfiedBy(req, int64(len(list))))
		} else {
			list = append(list, v)
		}
	}

	return list, tbl.Logger().LogIfError(tbl.Ctx(), require.ChainErrorIfQueryNotSatisfiedBy(rows.Err(), req, int64(len(list))))
}

// SliceFloat32PtrList requests a columnar slice of float32s from a specified nullable column.
func SliceFloat32PtrList(tbl pgxapi.Table, req require.Requirement, sqlname string, wh where.Expression, qc where.QueryConstraint) ([]float32, error) {
	query, args := sliceSql(tbl, sqlname, wh, qc)
	rows, err := Query(tbl, query, args...)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var v sql.NullFloat64
	list := make([]float32, 0, 10)

	for rows.Next() {
		err := rows.Scan(&v)
		if err == sql.ErrNoRows {
			return list, tbl.Logger().LogIfError(tbl.Ctx(), require.ErrorIfQueryNotSatisfiedBy(req, int64(len(list))))
		} else if v.Valid {
			list = append(list, float32(v.Float64))
		}
	}

	return list, tbl.Logger().LogIfError(tbl.Ctx(), require.ChainErrorIfQueryNotSatisfiedBy(rows.Err(), req, int64(len(list))))
}

// Float32AsInterfaceSlice adapts a slice of float32 to []interface{}.
func Float32AsInterfaceSlice(values []float32) []interface{} {
	ii := make([]interface{}, len(values))
	for i, v := range values {
		ii[i] = v
	}
	return ii
}
