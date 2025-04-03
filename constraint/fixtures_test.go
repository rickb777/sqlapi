package constraint_test

import (
	"context"
	"fmt"
	"testing"

	"github.com/rickb777/expect"
	"github.com/rickb777/sqlapi"
	"github.com/rickb777/sqlapi/driver"
	"github.com/rickb777/where/dialect"
)

func insertFixtures(t *testing.T, d sqlapi.Execer) (aid1, aid2, aid3, aid4 int64) {
	for _, s := range createTablesSql(d.Dialect()) {
		_, err := d.Exec(context.Background(), s)
		expect.Error(err).ToBeNil(t)
	}

	aid1 = insertOne(t, d, address1)
	aid2 = insertOne(t, d, address2)
	aid3 = insertOne(t, d, address3)
	aid4 = insertOne(t, d, address4)

	insertOne(t, d, fmt.Sprintf(person1a, aid1))
	insertOne(t, d, fmt.Sprintf(person1b, aid1))
	insertOne(t, d, fmt.Sprintf(person2a, aid2))

	return aid1, aid2, aid3, aid4
}

func insertOne(t *testing.T, d sqlapi.Execer, query string) int64 {
	id, err := d.Insert(context.Background(), "id", query)
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
	`DROP TABLE IF EXISTS constraint_addresses`,
	`DROP TABLE IF EXISTS constraint_persons`,

	`CREATE TABLE constraint_addresses (
	id        integer primary key autoincrement,
	xlines    text,
	postcode  text
	)`,

	`CREATE TABLE constraint_persons (
	id        integer primary key autoincrement,
	name      text,
	addressid integer default null
	)`,
}

var createTablesMysql = []string{
	`DROP TABLE IF EXISTS constraint_addresses`,
	`DROP TABLE IF EXISTS constraint_persons`,

	`CREATE TABLE constraint_addresses (
	id        int primary key auto_increment,
	xlines    text,
	postcode  text
	)`,

	`CREATE TABLE constraint_persons (
	id        int primary key auto_increment,
	name      text,
	addressid int default null
	)`,
}

var createTablesPostgresql = []string{
	`DROP TABLE IF EXISTS constraint_addresses`,
	`DROP TABLE IF EXISTS constraint_persons`,

	`CREATE TABLE constraint_addresses (
	id        serial primary key,
	xlines    text,
	postcode  text
	)`,

	`CREATE TABLE constraint_persons (
	id        serial primary key,
	name      text,
	addressid integer default null
	)`,
}

const address1 = `INSERT INTO constraint_addresses (xlines, postcode) VALUES ('Laurel Cottage', 'FX1 1AA')`
const address2 = `INSERT INTO constraint_addresses (xlines, postcode) VALUES ('2 Nutmeg Lane', 'FX1 2BB')`
const address3 = `INSERT INTO constraint_addresses (xlines, postcode) VALUES ('Corner Shop', 'FX1 3CC')`
const address4 = `INSERT INTO constraint_addresses (xlines, postcode) VALUES ('4 The Oaks', 'FX1 5EE')`

const person1a = `INSERT INTO constraint_persons (name, addressid) VALUES ('John Brown', %d)`
const person1b = `INSERT INTO constraint_persons (name, addressid) VALUES ('Mary Brown', %d)`
const person2a = `INSERT INTO constraint_persons (name, addressid) VALUES ('Anne Bollin', %d)`
