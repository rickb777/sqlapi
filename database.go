package sqlapi

import (
	"context"
	"database/sql"
	"github.com/rickb777/sqlapi/dialect"
	"github.com/rickb777/sqlapi/require"
	"github.com/rickb777/sqlapi/util"
	"log"
	"strings"
	"sync/atomic"
)

// Database typically wraps a *sql.DB with a dialect and (optionally) a logger.
// It's safe for concurrent use by multiple goroutines.
// See NewDatabase.
type Database interface {
	DB() Execer
	BeginTx(ctx context.Context, opts *sql.TxOptions) (SqlTx, error)
	Begin() (SqlTx, error)
	Dialect() dialect.Dialect
	Logger() *log.Logger
	Wrapper() interface{}
	PingContext(ctx context.Context) error
	Ping() error
	Stats() sql.DBStats

	TraceLogging(on bool)
	LogQuery(query string, args ...interface{})
	LogIfError(err error) error
	LogError(err error) error

	Exec(query string, args ...interface{}) (sql.Result, error)
	ExecContext(ctx context.Context, query string, args ...interface{}) (sql.Result, error)
	Prepare(query string) (SqlStmt, error)
	PrepareContext(ctx context.Context, query string) (SqlStmt, error)
	Query(query string, args ...interface{}) (SqlRows, error)
	QueryContext(ctx context.Context, query string, args ...interface{}) (SqlRows, error)
	QueryRow(query string, args ...interface{}) SqlRow
	QueryRowContext(ctx context.Context, query string, args ...interface{}) SqlRow

	ScanStringList(req require.Requirement, rows SqlRows) ([]string, error)
	ScanIntList(req require.Requirement, rows SqlRows) ([]int, error)
	ScanUintList(req require.Requirement, rows SqlRows) ([]uint, error)
	ScanInt64List(req require.Requirement, rows SqlRows) ([]int64, error)
	ScanUint64List(req require.Requirement, rows SqlRows) ([]uint64, error)
	ScanInt32List(req require.Requirement, rows SqlRows) ([]int32, error)
	ScanUint32List(req require.Requirement, rows SqlRows) ([]uint32, error)
	ScanInt16List(req require.Requirement, rows SqlRows) ([]int16, error)
	ScanUint16List(req require.Requirement, rows SqlRows) ([]uint16, error)
	ScanInt8List(req require.Requirement, rows SqlRows) ([]int8, error)
	ScanUint8List(req require.Requirement, rows SqlRows) ([]uint8, error)
	ScanFloat32List(req require.Requirement, rows SqlRows) ([]float32, error)
	ScanFloat64List(req require.Requirement, rows SqlRows) ([]float64, error)

	TableExists(name TableName) (yes bool, err error)
	ListTables() (util.StringList, error)
}

// database wraps a *sql.DB with a dialect and (optionally) a logger.
// It's safe for concurrent use by multiple goroutines.
type database struct {
	db         Execer
	dialect    dialect.Dialect
	logger     *log.Logger
	lgrEnabled int32
	wrapper    interface{}
}

// NewDatabase creates a new database handler, which wraps the core *sql.DB along with
// the appropriate dialect.
//
// You can supply the logger you need, or else nil. If not nil, all queries will be logged
// and all database errors will be logged. Once constructed, the logger itself cannot be
// changed, but its output writer can (via the SetOutput method on Logger). Logging can
// be enabled and disabled as needed by using the TraceLogging method.
//
// The wrapper holds some associated data your application needs for this database, if any.
// Otherwise this should be nil. As with the logger, it cannot be changed after construction.
func NewDatabase(db Execer, dialect dialect.Dialect, logger *log.Logger, wrapper interface{}) Database {
	var enabled int32 = 0
	if logger != nil {
		enabled = 1
	}
	return &database{
		db:         db,
		dialect:    dialect,
		logger:     logger,
		lgrEnabled: enabled,
		wrapper:    wrapper,
	}
}

// DB gets the Execer, which is a *sql.DB (except during testing using mocks).
func (database *database) DB() Execer {
	return database.db
}

// BeginTx starts a transaction.
//
// The context is used until the transaction is committed or rolled back. If this
// context is cancelled, the sql package will roll back the transaction. In this
// case, Tx.Commit will then return an error.
//
// The provided TxOptions is optional and may be nil if defaults should be used.
// If a non-default isolation level is used that the driver doesn't support,
// an error will be returned.
//
// Panics if the Execer is not a TxStarter.
func (database *database) BeginTx(ctx context.Context, opts *sql.TxOptions) (SqlTx, error) {
	return database.db.(TxStarter).BeginTx(ctx, opts)
}

// Begin starts a transaction using default options. The default isolation level is
// dependent on the driver.
func (database *database) Begin() (SqlTx, error) {
	return database.BeginTx(context.Background(), nil)
}

// Dialect gets the current SQL dialect. This choice is determined when the database is
// constructed and doesn't subsequently change.
func (database *database) Dialect() dialect.Dialect {
	return database.dialect
}

// Logger gets the trace logger. Note that you can use this to rotate the output writer
// via its SetOutput method. Also, it can even disable it completely (via ioutil.Discard).
func (database *database) Logger() *log.Logger {
	return database.logger
}

// Wrapper gets whatever structure is present, as needed.
func (database *database) Wrapper() interface{} {
	return database.wrapper
}

//-------------------------------------------------------------------------------------------------

// PingContext verifies a connection to the database is still alive,
// establishing a connection if necessary.
func (database *database) PingContext(ctx context.Context) error {
	return database.db.(*sql.DB).PingContext(ctx)
}

// Ping verifies a connection to the database is still alive,
// establishing a connection if necessary.
func (database *database) Ping() error {
	return database.PingContext(context.Background())
}

// Exec executes a query without returning any rows.
// The args are for any placeholder parameters in the query.
func (database *database) Exec(query string, args ...interface{}) (sql.Result, error) {
	return database.ExecContext(context.Background(), query, args...)
}

// ExecContext executes a query without returning any rows.
// The args are for any placeholder parameters in the query.
func (database *database) ExecContext(ctx context.Context, query string, args ...interface{}) (sql.Result, error) {
	return database.db.ExecContext(ctx, query, args...)
}

// Prepare creates a prepared statement for later queries or executions.
// Multiple queries or executions may be run concurrently from the
// returned statement.
// The caller must call the statement's Close method
// when the statement is no longer needed.
func (database *database) Prepare(query string) (SqlStmt, error) {
	return database.PrepareContext(context.Background(), query)
}

// PrepareContext creates a prepared statement for later queries or executions.
// Multiple queries or executions may be run concurrently from the
// returned statement.
// The caller must call the statement's Close method
// when the statement is no longer needed.
//
// The provided context is used for the preparation of the statement, not for the
// execution of the statement.
func (database *database) PrepareContext(ctx context.Context, query string) (SqlStmt, error) {
	return database.db.PrepareContext(ctx, query)
}

// Query executes a query that returns rows, typically a SELECT.
// The args are for any placeholder parameters in the query.
func (database *database) Query(query string, args ...interface{}) (SqlRows, error) {
	return database.QueryContext(context.Background(), query, args...)
}

// QueryContext executes a query that returns rows, typically a SELECT.
// The args are for any placeholder parameters in the query.
func (database *database) QueryContext(ctx context.Context, query string, args ...interface{}) (SqlRows, error) {
	return database.db.QueryContext(ctx, query, args...)
}

// QueryRow executes a query that is expected to return at most one row.
// QueryRow always returns a non-nil value. Errors are deferred until
// Row's Scan method is called.
// If the query selects no rows, the *Row's Scan will return ErrNoRows.
// Otherwise, the *Row's Scan scans the first selected row and discards
// the rest.
func (database *database) QueryRow(query string, args ...interface{}) SqlRow {
	return database.QueryRowContext(context.Background(), query, args...)
}

// QueryRowContext executes a query that is expected to return at most one row.
// QueryRowContext always returns a non-nil value. Errors are deferred until
// Row's Scan method is called.
// If the query selects no rows, the *Row's Scan will return ErrNoRows.
// Otherwise, the *Row's Scan scans the first selected row and discards
// the rest.
func (database *database) QueryRowContext(ctx context.Context, query string, args ...interface{}) SqlRow {
	return database.db.QueryRowContext(ctx, query, args...)
}

// Stats returns database statistics.
func (database *database) Stats() sql.DBStats {
	return database.db.(*sql.DB).Stats()
}

//-------------------------------------------------------------------------------------------------

// TraceLogging turns query trace logging on or off. This has no effect unless the database was
// created with a non-nil logger.
func (database *database) TraceLogging(on bool) {
	if on && database.logger != nil {
		atomic.StoreInt32(&database.lgrEnabled, 1)
	} else {
		atomic.StoreInt32(&database.lgrEnabled, 0)
	}
}

func (database *database) loggingEnabled() bool {
	return atomic.LoadInt32(&database.lgrEnabled) != 0
}

// LogQuery writes query info to the logger, if it is not nil.
func (database *database) LogQuery(query string, args ...interface{}) {
	if database.loggingEnabled() {
		query = strings.TrimSpace(query)
		if len(args) > 0 {
			ss := make([]interface{}, len(args))
			for i, v := range args {
				ss[i] = derefArg(v)
			}
			database.logger.Printf(query+" %v\n", ss)
		} else {
			database.logger.Println(query)
		}
	}
}

func derefArg(arg interface{}) interface{} {
	switch v := arg.(type) {
	case *int:
		return *v
	case *int8:
		return *v
	case *int16:
		return *v
	case *int32:
		return *v
	case *int64:
		return *v
	case *uint:
		return *v
	case *uint8:
		return *v
	case *uint16:
		return *v
	case *uint32:
		return *v
	case *uint64:
		return *v
	case *float32:
		return *v
	case *float64:
		return *v
	case *bool:
		return *v
	case *string:
		return *v
	}
	return arg
}

// LogIfError writes error info to the logger, if both the logger and the error are non-nil.
// It returns the error.
func (database *database) LogIfError(err error) error {
	if database.loggingEnabled() && err != nil {
		database.logger.Printf("Error: %s\n", err)
	}
	return err
}

// LogError writes error info to the logger, if the logger is not nil. It returns the error.
func (database *database) LogError(err error) error {
	if database.loggingEnabled() {
		database.logger.Printf("Error: %s\n", err)
	}
	return err
}

//-------------------------------------------------------------------------------------------------

// ScanStringList processes result rows to extract a list of strings.
// The result set should have been produced via a SELECT statement on just one column.
func (database *database) ScanStringList(req require.Requirement, rows SqlRows) ([]string, error) {
	var v string
	list := make([]string, 0, 10)

	for rows.Next() {
		err := rows.Scan(&v)
		if err == sql.ErrNoRows {
			return list, database.LogIfError(require.ErrorIfQueryNotSatisfiedBy(req, int64(len(list))))
		} else {
			list = append(list, v)
		}
	}
	return list, database.LogIfError(require.ChainErrorIfQueryNotSatisfiedBy(rows.Err(), req, int64(len(list))))
}

// ScanIntList processes result rows to extract a list of ints.
// The result set should have been produced via a SELECT statement on just one column.
func (database *database) ScanIntList(req require.Requirement, rows SqlRows) ([]int, error) {
	var v int
	list := make([]int, 0, 10)

	for rows.Next() {
		err := rows.Scan(&v)
		if err == sql.ErrNoRows {
			return list, database.LogIfError(require.ErrorIfQueryNotSatisfiedBy(req, int64(len(list))))
		} else {
			list = append(list, v)
		}
	}
	return list, database.LogIfError(require.ChainErrorIfQueryNotSatisfiedBy(rows.Err(), req, int64(len(list))))
}

// ScanUintList processes result rows to extract a list of uints.
// The result set should have been produced via a SELECT statement on just one column.
func (database *database) ScanUintList(req require.Requirement, rows SqlRows) ([]uint, error) {
	var v uint
	list := make([]uint, 0, 10)

	for rows.Next() {
		err := rows.Scan(&v)
		if err == sql.ErrNoRows {
			return list, database.LogIfError(require.ErrorIfQueryNotSatisfiedBy(req, int64(len(list))))
		} else {
			list = append(list, v)
		}
	}
	return list, database.LogIfError(require.ChainErrorIfQueryNotSatisfiedBy(rows.Err(), req, int64(len(list))))
}

// ScanInt64List processes result rows to extract a list of int64s.
// The result set should have been produced via a SELECT statement on just one column.
func (database *database) ScanInt64List(req require.Requirement, rows SqlRows) ([]int64, error) {
	var v int64
	list := make([]int64, 0, 10)

	for rows.Next() {
		err := rows.Scan(&v)
		if err == sql.ErrNoRows {
			return list, database.LogIfError(require.ErrorIfQueryNotSatisfiedBy(req, int64(len(list))))
		} else {
			list = append(list, v)
		}
	}
	return list, database.LogIfError(require.ChainErrorIfQueryNotSatisfiedBy(rows.Err(), req, int64(len(list))))
}

// ScanUint64List processes result rows to extract a list of uint64s.
// The result set should have been produced via a SELECT statement on just one column.
func (database *database) ScanUint64List(req require.Requirement, rows SqlRows) ([]uint64, error) {
	var v uint64
	list := make([]uint64, 0, 10)

	for rows.Next() {
		err := rows.Scan(&v)
		if err == sql.ErrNoRows {
			return list, database.LogIfError(require.ErrorIfQueryNotSatisfiedBy(req, int64(len(list))))
		} else {
			list = append(list, v)
		}
	}
	return list, database.LogIfError(require.ChainErrorIfQueryNotSatisfiedBy(rows.Err(), req, int64(len(list))))
}

// ScanInt32List processes result rows to extract a list of int32s.
// The result set should have been produced via a SELECT statement on just one column.
func (database *database) ScanInt32List(req require.Requirement, rows SqlRows) ([]int32, error) {
	var v int32
	list := make([]int32, 0, 10)

	for rows.Next() {
		err := rows.Scan(&v)
		if err == sql.ErrNoRows {
			return list, database.LogIfError(require.ErrorIfQueryNotSatisfiedBy(req, int64(len(list))))
		} else {
			list = append(list, v)
		}
	}
	return list, database.LogIfError(require.ChainErrorIfQueryNotSatisfiedBy(rows.Err(), req, int64(len(list))))
}

// ScanUint32List processes result rows to extract a list of uint32s.
// The result set should have been produced via a SELECT statement on just one column.
func (database *database) ScanUint32List(req require.Requirement, rows SqlRows) ([]uint32, error) {
	var v uint32
	list := make([]uint32, 0, 10)

	for rows.Next() {
		err := rows.Scan(&v)
		if err == sql.ErrNoRows {
			return list, database.LogIfError(require.ErrorIfQueryNotSatisfiedBy(req, int64(len(list))))
		} else {
			list = append(list, v)
		}
	}
	return list, database.LogIfError(require.ChainErrorIfQueryNotSatisfiedBy(rows.Err(), req, int64(len(list))))
}

// ScanInt16List processes result rows to extract a list of int32s.
// The result set should have been produced via a SELECT statement on just one column.
func (database *database) ScanInt16List(req require.Requirement, rows SqlRows) ([]int16, error) {
	var v int16
	list := make([]int16, 0, 10)

	for rows.Next() {
		err := rows.Scan(&v)
		if err == sql.ErrNoRows {
			return list, database.LogIfError(require.ErrorIfQueryNotSatisfiedBy(req, int64(len(list))))
		} else {
			list = append(list, v)
		}
	}
	return list, database.LogIfError(require.ChainErrorIfQueryNotSatisfiedBy(rows.Err(), req, int64(len(list))))
}

// ScanUint16List processes result rows to extract a list of uint16s.
// The result set should have been produced via a SELECT statement on just one column.
func (database *database) ScanUint16List(req require.Requirement, rows SqlRows) ([]uint16, error) {
	var v uint16
	list := make([]uint16, 0, 10)

	for rows.Next() {
		err := rows.Scan(&v)
		if err == sql.ErrNoRows {
			return list, database.LogIfError(require.ErrorIfQueryNotSatisfiedBy(req, int64(len(list))))
		} else {
			list = append(list, v)
		}
	}
	return list, database.LogIfError(require.ChainErrorIfQueryNotSatisfiedBy(rows.Err(), req, int64(len(list))))
}

// ScanInt8List processes result rows to extract a list of int8s.
// The result set should have been produced via a SELECT statement on just one column.
func (database *database) ScanInt8List(req require.Requirement, rows SqlRows) ([]int8, error) {
	var v int8
	list := make([]int8, 0, 10)

	for rows.Next() {
		err := rows.Scan(&v)
		if err == sql.ErrNoRows {
			return list, database.LogIfError(require.ErrorIfQueryNotSatisfiedBy(req, int64(len(list))))
		} else {
			list = append(list, v)
		}
	}
	return list, database.LogIfError(require.ChainErrorIfQueryNotSatisfiedBy(rows.Err(), req, int64(len(list))))
}

// ScanUint8List processes result rows to extract a list of uint8s.
// The result set should have been produced via a SELECT statement on just one column.
func (database *database) ScanUint8List(req require.Requirement, rows SqlRows) ([]uint8, error) {
	var v uint8
	list := make([]uint8, 0, 10)

	for rows.Next() {
		err := rows.Scan(&v)
		if err == sql.ErrNoRows {
			return list, database.LogIfError(require.ErrorIfQueryNotSatisfiedBy(req, int64(len(list))))
		} else {
			list = append(list, v)
		}
	}
	return list, database.LogIfError(require.ChainErrorIfQueryNotSatisfiedBy(rows.Err(), req, int64(len(list))))
}

// ScanFloat32List processes result rows to extract a list of float32s.
// The result set should have been produced via a SELECT statement on just one column.
func (database *database) ScanFloat32List(req require.Requirement, rows SqlRows) ([]float32, error) {
	var v float32
	list := make([]float32, 0, 10)

	for rows.Next() {
		err := rows.Scan(&v)
		if err == sql.ErrNoRows {
			return list, database.LogIfError(require.ErrorIfQueryNotSatisfiedBy(req, int64(len(list))))
		} else {
			list = append(list, v)
		}
	}
	return list, database.LogIfError(require.ChainErrorIfQueryNotSatisfiedBy(rows.Err(), req, int64(len(list))))
}

// ScanFloat64List processes result rows to extract a list of float64s.
// The result set should have been produced via a SELECT statement on just one column.
func (database *database) ScanFloat64List(req require.Requirement, rows SqlRows) ([]float64, error) {
	var v float64
	list := make([]float64, 0, 10)

	for rows.Next() {
		err := rows.Scan(&v)
		if err == sql.ErrNoRows {
			return list, database.LogIfError(require.ErrorIfQueryNotSatisfiedBy(req, int64(len(list))))
		} else {
			list = append(list, v)
		}
	}
	return list, database.LogIfError(require.ChainErrorIfQueryNotSatisfiedBy(rows.Err(), req, int64(len(list))))
}

//-------------------------------------------------------------------------------------------------

// DoesTableExist gets all the table names in the database/schema.
func (database *database) TableExists(name TableName) (yes bool, err error) {
	wanted := name.String()
	rows, err := database.db.QueryContext(context.Background(), database.dialect.ShowTables())
	if err != nil {
		return false, err
	}
	defer rows.Close()

	for rows.Next() {
		var s string
		rows.Scan(&s)
		if s == wanted {
			return true, rows.Err()
		}
	}
	return false, rows.Err()
}

// ListTables gets all the table names in the database/schema.
func (database *database) ListTables() (util.StringList, error) {
	ss := make(util.StringList, 0)
	rows, err := database.db.QueryContext(context.Background(), database.dialect.ShowTables())
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	for rows.Next() {
		var s string
		rows.Scan(&s)
		ss = append(ss, s)
	}
	return ss, rows.Err()
}
