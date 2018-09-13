# sqlapi

[![GoDoc](https://img.shields.io/badge/api-Godoc-blue.svg?style=flat-square)](https://godoc.org/github.com/rickb777/sqlapi)
[![Build Status](https://travis-ci.org/rickb777/sqlapi.svg?branch=master)](https://travis-ci.org/rickb777/sqlapi)
[![Code Coverage](https://img.shields.io/coveralls/rickb777/sqlapi.svg)](https://coveralls.io/r/rickb777/sqlapi)
[![Issues](https://img.shields.io/github/issues/rickb777/sqlapi.svg)](https://github.com/rickb777/sqlapi/issues)

**sqlgen** generates SQL statements and database helper functions from your Go structs. It can be used in
place of a simple ORM or hand-written SQL. **sqlapi** (this package) supports the generated code.

See the [demo](https://github.com/rickb777/sqlgen2/tree/master/demo) directory for examples. Look in the
generated files `*_sql.go` and the hand-crafted files (`hook.go`, `issue.go`, `user.go`).

## Features

* Auto-generates DAO-style table-support code for SQL databases.
* Sophisticated parsing of a Go `struct` that describes records in the table.
* Allows nesting of structs, fields that are structs or pointers to structs etc.
* Struct tags give fine control over the semantics of each field.
* Supports indexes and constraints.
* Supports foreign key relationships between tables.
* Helps you develop code for joins and views.
* Supports JSON-encoded columns, allowing a more no-SQL model when needed.
* Provides a builder-style API for constructing where-clauses and query constraints.
* Allows declarative requirements on the expected result of each query, enhancing error checking.
* Very flexible configuration.
* Fast and easy to use.

Currently, support is included for **MySQL**, **PostgreSQL** and **SQLite**. Other dialects can be added relatively easy - send a Pull Request!

## Install

Add to your imports some references to github.com/rickb777/sqlapi and then install with this command:

```
go mod tidy
```

If you're using an older version of Go, you'll need

```
go get github.com/rickb777/sqlapi
```

Please continue reading about [sqlgen2](https://github.com/rickb777/sqlgen2/tree/master/README.md).
