package sqlapi

import (
	"context"
	"database/sql"
	"github.com/rickb777/collection"
	"github.com/rickb777/sqlapi/dialect"
	"log"
	"regexp"
)

type DBStats = sql.DBStats

// Database typically wraps a *sql.DB with a dialect and (optionally) a logger.
// It's safe for concurrent use by multiple goroutines.
// See NewDatabase.
type Database interface {
	DB() SqlDB
	Dialect() dialect.Dialect
	Logger() Logger
	Wrapper() interface{}
	ListTables(re *regexp.Regexp) (collection.StringList, error)
}

// database wraps a *sql.DB with a dialect and (optionally) a logger.
// It's safe for concurrent use by multiple goroutines.
type database struct {
	db      SqlDB
	dialect dialect.Dialect
	logger  *toggleLogger
	wrapper interface{}
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
	tl := &toggleLogger{} // not nil but inactive
	if logger != nil {
		tl = &toggleLogger{lgr: logger, enabled: 1}
	}
	return &database{
		db:      db,
		dialect: dialect,
		logger:  tl,
		wrapper: wrapper,
	}
}

// DB gets the Execer, which is a *sql.DB (except during testing using mocks).
func (database *database) DB() SqlDB {
	return database.db
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

// ListTables gets all the table names in the database/schema.
// The regular expression supplies a filter: only names that match are returned.
// If the regular expression is nil, all tables names are returned.
func (database *database) ListTables(re *regexp.Regexp) (collection.StringList, error) {
	ss := make(collection.StringList, 0)
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
