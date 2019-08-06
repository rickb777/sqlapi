package dialect

import (
	"bytes"
	"fmt"
	"github.com/mercury-holidays/sqlapi/schema"
	"github.com/mercury-holidays/sqlapi/types"
	"github.com/rickb777/where/quote"
	"strconv"
)

type postgres struct {
	q quote.Quoter
}

var Postgres Dialect = postgres{q: quote.AnsiQuoter}

func (d postgres) Index() int {
	return PostgresIndex
}

func (d postgres) String() string {
	return "Postgres"
}

func (d postgres) Alias() string {
	return "PostgreSQL"
}

func (d postgres) Quoter() quote.Quoter {
	return d.q
}

func (d postgres) WithQuoter(q quote.Quoter) Dialect {
	d.q = q
	return d
}

// https://www.postgresql.org/docs/9.6/static/datatype.html
// https://www.convert-in.com/mysql-to-postgres-types-mapping.htm

func (dialect postgres) FieldAsColumn(field *schema.Field) string {
	tags := field.GetTags()

	switch field.Encode {
	case schema.ENCJSON:
		return "json"
	case schema.ENCTEXT:
		return "text"
	}

	column := "bytea"
	dflt := tags.Default

	switch field.Type.Base {
	case types.Int, types.Int64:
		column = "bigint"
	case types.Int8:
		column = "int8"
	case types.Int16:
		column = "smallint"
	case types.Int32:
		column = "integer"
	case types.Uint, types.Uint64:
		// Some DBs (including postgresql) do not support unsigned integers. Rejecting
		// uint64 >= 1<<63 prevents them becoming indistinguishable from int64s < 0. If you need
		// to insert such int64 anyway, you have to explicitly convert it on input to int64 and
		// convert it back to uint64 on output - and take care to never insert a signed integer
		// into the same column.
		column = "bigint" // incomplete number range but more storage efficiency
	case types.Uint8:
		column = "smallint"
	case types.Uint16:
		column = "integer"
	case types.Uint32:
		column = "bigint"
	case types.Float32:
		column = "real"
	case types.Float64:
		column = "double precision"
	case types.Bool:
		column = "boolean"
	case types.String:
		column = "text"
		dflt = fmt.Sprintf("'%s'", tags.Default)
	}

	// postgres uses a special column type
	// for autoincrementing keys.
	if tags.Auto {
		switch field.Type.Base {
		case types.Int, types.Int64, types.Uint64:
			column = "bigserial"
		default:
			column = "serial"
		}
	}

	return fieldTags(field.Type.IsPtr, tags, column, dflt)
}

func (dialect postgres) InsertHasReturningPhrase() bool {
	return true
}

func (dialect postgres) TruncateDDL(tableName string, force bool) []string {
	if force {
		return []string{fmt.Sprintf("TRUNCATE %s CASCADE", dialect.Quoter().Quote(tableName))}
	}

	return []string{fmt.Sprintf("TRUNCATE %s RESTRICT", dialect.Quoter().Quote(tableName))}
}

func (dialect postgres) ShowTables() string {
	return `SELECT tablename FROM pg_catalog.pg_tables`
}

//-------------------------------------------------------------------------------------------------

func (dialect postgres) HasNumberedPlaceholders() bool {
	return true
}

func (dialect postgres) Placeholders(n int) string {
	if n == 0 {
		return ""
	} else if n <= 9 {
		return postgresPlaceholders[:n*3-1]
	}
	buf := bytes.NewBufferString(postgresPlaceholders)
	for idx := 10; idx <= n; idx++ {
		if idx > 1 {
			buf.WriteByte(',')
		}
		buf.WriteByte('$')
		buf.WriteString(strconv.Itoa(idx))
	}
	return buf.String()
}

const postgresPlaceholders = "$1,$2,$3,$4,$5,$6,$7,$8,$9"

// ReplacePlaceholders converts a string containing '?' placeholders to
// the form used by PostgreSQL.
func (dialect postgres) ReplacePlaceholders(sql string, _ []interface{}) string {
	buf := &bytes.Buffer{}
	idx := 1
	for _, r := range sql {
		if r == '?' {
			buf.WriteByte('$')
			buf.WriteString(strconv.Itoa(idx))
			idx++
		} else {
			buf.WriteRune(r)
		}
	}
	return buf.String()
}

func (dialect postgres) CreateTableSettings() string {
	return ""
}
