package driver

import (
	"fmt"

	"github.com/rickb777/where/dialect"
)

type pgx struct {
	postgres
}

func Pgx(d ...dialect.DialectConfig) Dialect {
	return pgx{postgres: postgres{of(dialect.PostgresConfig, d...)}}
}

func (d pgx) Index() dialect.Dialect {
	return dialect.Postgres
}

func (d pgx) String() string {
	if d.d.Quoter != nil {
		return fmt.Sprintf("Pgx/%s", d.d.Quoter)
	}
	return "Pgx"
}

func (d pgx) Name() string {
	return "Pgx"
}

func (d pgx) Alias() string {
	return "pgx"
}
