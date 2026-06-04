package driver

import (
	"fmt"

	"github.com/rickb777/where/v2/dialect"
	"github.com/rickb777/where/v2/quote"
)

type pgx struct {
	postgres
}

func Pgx(q ...quote.Quoter) Dialect {
	return pgx{postgres: postgres{d: dialect.Postgres, o: dialect.Dollar, q: of(dialect.PostgresQuoter, q...)}}
}

func (d pgx) Index() dialect.Dialect {
	return dialect.Postgres
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
