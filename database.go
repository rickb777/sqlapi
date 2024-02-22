package sqlapi

import (
	"context"
	"database/sql"
	"regexp"
)

type DBStats = sql.DBStats

//-------------------------------------------------------------------------------------------------

// ListTables gets all the table names in the database/schema.
// The regular expression supplies a filter: only names that match are returned.
// If the regular expression is nil, all table names are returned.
func ListTables(ex Execer, re *regexp.Regexp) ([]string, error) {
	ss := make([]string, 0)
	rows, err := ex.Query(context.Background(), ex.Dialect().ShowTables())
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
