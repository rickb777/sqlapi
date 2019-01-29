package schema

type pgx struct {
	postgres
}

var Pgx Dialect = pgx{}

func (d pgx) Index() int {
	return PgxIndex
}

func (d pgx) String() string {
	return "Pgx"
}

func (d pgx) Alias() string {
	return "pgx"
}
