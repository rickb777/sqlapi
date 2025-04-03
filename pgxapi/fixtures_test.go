package pgxapi

import (
	"context"
	"fmt"
	"testing"

	"github.com/rickb777/expect"
	"github.com/rickb777/sqlapi/driver"
	"github.com/rickb777/where/dialect"
)

func insertFixtures(t *testing.T, e Execer) (aid1, aid2, aid3, aid4 int64) {
	for _, s := range createTablesSql(e.Dialect()) {
		_, err := e.Exec(context.Background(), s)
		expect.Error(err).ToBeNil(t)
	}

	aid1 = insertOne(t, e, address1)
	aid2 = insertOne(t, e, address2)
	aid3 = insertOne(t, e, address3)
	aid4 = insertOne(t, e, address4)

	insertOne(t, e, fmt.Sprintf(person1a, aid1))
	insertOne(t, e, fmt.Sprintf(person1b, aid1))
	insertOne(t, e, fmt.Sprintf(person2a, aid2))

	return aid1, aid2, aid3, aid4
}

func insertOne(t *testing.T, e Execer, query string) int64 {
	id, err := e.Insert(context.Background(), "id", query)
	expect.Error(err).ToBeNil(t)
	return id
}

func createTablesSql(di driver.Dialect) []string {
	switch di.Index() {
	case dialect.Sqlite:
		return createTablesSqlite
	case dialect.Mysql:
		return createTablesMysql
	case dialect.Postgres:
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
