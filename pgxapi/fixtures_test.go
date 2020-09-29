package pgxapi

import (
	"context"
	"fmt"
	"testing"

	"github.com/rickb777/sqlapi/dialect"

	"github.com/onsi/gomega"
)

func insertFixtures(t *testing.T, e Execer) (aid1, aid2, aid3, aid4 int64) {
	g := gomega.NewGomegaWithT(t)

	for _, s := range createTablesSql(e.Dialect()) {
		_, err := e.Exec(context.Background(), s)
		g.Expect(err).To(gomega.BeNil())
	}

	aid1 = insertOne(g, e, address1)
	aid2 = insertOne(g, e, address2)
	aid3 = insertOne(g, e, address3)
	aid4 = insertOne(g, e, address4)

	insertOne(g, e, fmt.Sprintf(person1a, aid1))
	insertOne(g, e, fmt.Sprintf(person1b, aid1))
	insertOne(g, e, fmt.Sprintf(person2a, aid2))

	return aid1, aid2, aid3, aid4
}

func insertOne(g *gomega.GomegaWithT, e Execer, query string) int64 {
	id, err := e.Insert(context.Background(), "id", query)
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
const address4 = `INSERT INTO pfx_addresses (xlines, postcode) VALUES ('4 The Oaks', 'FX1 4DD')`

const person1a = `INSERT INTO pfx_persons (name, addressid) VALUES ('John Brown', %d)`
const person1b = `INSERT INTO pfx_persons (name, addressid) VALUES ('Mary Brown', %d)`
const person2a = `INSERT INTO pfx_persons (name, addressid) VALUES ('Anne Bollin', %d)`
