package pgxapi

import (
	"context"
	"database/sql"
	"regexp"

	"github.com/rickb777/collection"
	"github.com/rickb777/sqlapi/dialect"
)

type DBStats = sql.DBStats

// Database typically wraps a *pgx.ConnPool with a dialect and (optionally) a logger.
// It's safe for concurrent use by multiple goroutines.
// See NewDatabase.
type Database interface {
	DB() SqlDB
	Dialect() dialect.Dialect
	Logger() Logger
	Wrapper() interface{}
}

// database wraps a *sql.DB with a dialect and (optionally) a logger.
// It's safe for concurrent use by multiple goroutines.
type database struct {
	db      SqlDB
	dialect dialect.Dialect
	wrapper interface{}
}

// NewDatabase creates a new database handler, which wraps the core *sql.DB along with
// the appropriate dialect.
//
// The wrapper holds some associated data your application needs for this database, if any.
// Otherwise this should be nil. As with the logger, it cannot be changed after construction.
func NewDatabase(db SqlDB, dialect dialect.Dialect, wrapper interface{}) Database {
	return &database{
		db:      db,
		dialect: dialect,
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

// Logger gets the trace logger.
func (database *database) Logger() Logger {
	return database.db.Logger()
}

// Wrapper gets whatever structure is present, as needed.
func (database *database) Wrapper() interface{} {
	return database.wrapper
}

//-------------------------------------------------------------------------------------------------

// ListTables gets all the table names in the database/schema.
// The regular expression supplies a filter: only names that match are returned.
// If the regular expression is nil, all table names are returned.
func ListTables(ex Execer, re *regexp.Regexp) (collection.StringList, error) {
	ss := make(collection.StringList, 0)
	rows, err := ex.QueryContext(context.Background(), ex.Dialect().ShowTables())
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
