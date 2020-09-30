package dialect

import "fmt"

type pgx struct {
	postgres
}

var Pgx Dialect = pgx{}

func (d pgx) Index() int {
	return PgxIndex
}

func (d pgx) String() string {
	if d.q != nil {
		return fmt.Sprintf("Pgx/%s", d.q)
	}
	return "Pgx"
}

func (d pgx) Name() string {
	return "Pgx"
}

func (d pgx) Alias() string {
	return "pgx"
}
