package sqlapi

import (
	"context"
	"database/sql"
	"github.com/rickb777/sqlapi/dialect"
	"github.com/rickb777/sqlapi/util"
	"log"
	"regexp"
	"strings"
	"sync/atomic"
)

// Database typically wraps a *sql.DB with a dialect and (optionally) a logger.
// It's safe for concurrent use by multiple goroutines.
// See NewDatabase.
type Database interface {
	DB() Execer
	Dialect() dialect.Dialect
	Logger() Logger
	Wrapper() interface{}
	PingContext(ctx context.Context) error
	Ping() error
	Stats() sql.DBStats

	TraceLogging(on bool)
	LogQuery(query string, args ...interface{})
	LogIfError(err error) error
	LogError(err error) error

	//Exec(query string, args ...interface{}) (int64, error)
	//ExecContext(ctx context.Context, query string, args ...interface{}) (int64, error)
	//Prepare(query string) (SqlStmt, error)
	//PrepareContext(ctx context.Context, name, query string) (SqlStmt, error)
	//Query(query string, args ...interface{}) (SqlRows, error)
	//QueryContext(ctx context.Context, query string, args ...interface{}) (SqlRows, error)
	//QueryRow(query string, args ...interface{}) SqlRow
	//QueryRowContext(ctx context.Context, query string, args ...interface{}) SqlRow

	ListTables(re *regexp.Regexp) (util.StringList, error)
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
func NewDatabase(db SqlDB, dialect dialect.Dialect, logger *log.Logger, wrapper interface{}) Database {
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
// Panics if the Execer is not a SqlDB.
func (database *database) BeginTx(ctx context.Context, opts *sql.TxOptions) (SqlTx, error) {
	return database.db.(SqlDB).BeginTx(ctx, opts)
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
func (database *database) Logger() Logger {
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
	return database.db.(SqlDB).PingContext(ctx)
}

// Ping verifies a connection to the database is still alive,
// establishing a connection if necessary.
func (database *database) Ping() error {
	return database.PingContext(context.Background())
}

// Stats returns database statistics.
func (database *database) Stats() sql.DBStats {
	return database.db.(SqlDB).Stats()
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
			database.logger.Printf("%s %v\n", query, ss)
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

// ListTables gets all the table names in the database/schema.
// The regular expression supplies a filter: only names that match are returned.
// If the regular expression is nil, all tables names are returned.
func (database *database) ListTables(re *regexp.Regexp) (util.StringList, error) {
	ss := make(util.StringList, 0)
	rows, err := database.db.QueryContext(context.Background(), database.dialect.ShowTables())
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	for rows.Next() {
		var s string
		rows.Scan(&s)
		if re == nil || re.MatchString(s) {
			ss = append(ss, s)
		}
	}
	return ss, rows.Err()
}
