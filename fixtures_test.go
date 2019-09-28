package sqlapi_test

import (
	"context"
	"fmt"
	"testing"

	"github.com/onsi/gomega"
	"github.com/rickb777/sqlapi"
	"github.com/rickb777/sqlapi/dialect"
)

func insertFixtures(t *testing.T, d sqlapi.Database) (aid1, aid2, aid3, aid4 int64) {
	g := gomega.NewGomegaWithT(t)

	for _, s := range createTablesSql(d.Dialect()) {
		_, err := d.DB().ExecContext(context.Background(), s)
		g.Expect(err).To(gomega.BeNil())
	}

	aid1 = insertOne(g, d, address1)
	aid2 = insertOne(g, d, address2)
	aid3 = insertOne(g, d, address3)
	aid4 = insertOne(g, d, address4)

	insertOne(g, d, fmt.Sprintf(person1a, aid1))
	insertOne(g, d, fmt.Sprintf(person1b, aid1))
	insertOne(g, d, fmt.Sprintf(person2a, aid2))

	return aid1, aid2, aid3, aid4
}

func insertOne(g *gomega.GomegaWithT, d sqlapi.Database, query string) int64 {
	if !d.Dialect().HasLastInsertId() {
		query = query + " RETURNING id"
	}
	id, err := d.DB().InsertContext(context.Background(), query)
	g.Expect(err).To(gomega.BeNil())
	return id
}

func createTablesSql(di dialect.Dialect) []string {
	switch di.Index() {
	case dialect.SqliteIndex:
		return createTablesSqlite
	case dialect.MysqlIndex:
		return createTablesMysql
	case dialect.PostgresIndex:
		return createTablesPostgresql
	}
	panic(di.String() + " unsupported")
}

var createTablesSqlite = []string{
	`DROP TABLE IF EXISTS pfx_addresses`,
	`DROP TABLE IF EXISTS pfx_persons`,

	`CREATE TABLE pfx_addresses (
	id        integer primary key autoincrement,
	xlines    text,
	postcode  text
	)`,

	`CREATE TABLE pfx_persons (
	id        integer primary key autoincrement,
	name      text,
	addressid integer default null
	)`,
}

var createTablesMysql = []string{
	`DROP TABLE IF EXISTS pfx_addresses`,
	`DROP TABLE IF EXISTS pfx_persons`,

	`CREATE TABLE pfx_addresses (
	id        int primary key auto_increment,
	xlines    text,
	postcode  text
	)`,

	`CREATE TABLE pfx_persons (
	id        int primary key auto_increment,
	name      text,
	addressid int default null
	)`,
}

var createTablesPostgresql = []string{
	`DROP TABLE IF EXISTS pfx_addresses`,
	`DROP TABLE IF EXISTS pfx_persons`,

	`CREATE TABLE pfx_addresses (
	id        serial primary key,
	xlines    text,
	postcode  text
	)`,

	`CREATE TABLE pfx_persons (
	id        serial primary key,
	name      text,
	addressid integer default null
	)`,
}

const address1 = `INSERT INTO pfx_addresses (xlines, postcode) VALUES ('Laurel Cottage', 'FX1 1AA')`
const address2 = `INSERT INTO pfx_addresses (xlines, postcode) VALUES ('2 Nutmeg Lane', 'FX1 2BB')`
const address3 = `INSERT INTO pfx_addresses (xlines, postcode) VALUES ('Corner Shop', 'FX1 3CC')`
const address4 = `INSERT INTO pfx_addresses (xlines, postcode) VALUES ('4 The Oaks', 'FX1 5EE')`

const person1a = `INSERT INTO pfx_persons (name, addressid) VALUES ('John Brown', %d)`
const person1b = `INSERT INTO pfx_persons (name, addressid) VALUES ('Mary Brown', %d)`
const person2a = `INSERT INTO pfx_persons (name, addressid) VALUES ('Anne Bollin', %d)`
