# sqlapi

[![GoDoc](https://img.shields.io/badge/api-Godoc-blue.svg)](https://pkg.go.dev/github.com/rickb777/sqlapi)
[![Build Status](https://travis-ci.org/rickb777/sqlapi.svg?branch=master)](https://travis-ci.org/rickb777/sqlapi)
[![Code Coverage](https://img.shields.io/coveralls/rickb777/sqlapi.svg)](https://coveralls.io/r/rickb777/sqlapi)
[![Go Report Card](https://goreportcard.com/badge/github.com/rickb777/sqlapi)](https://goreportcard.com/report/github.com/rickb777/sqlapi)
[![Issues](https://img.shields.io/github/issues/rickb777/sqlapi.svg)](https://github.com/rickb777/sqlapi/issues)

**sqlgen** generates SQL statements and database helper functions from your Go structs. It can be used in
place of a simple ORM or hand-written SQL. **sqlapi** (this package) supports the generated code.

See the [demo](https://github.com/rickb777/sqlgen2/tree/master/demo) directory for examples. Look in the
generated files `*_sql.go` and the hand-crafted files (`hook.go`, `issue.go`, `user.go`).

Currently, support is included for **MySQL**, **PostgreSQL** and **SQLite**. Other dialects can be added relatively easy - send a Pull Request!

## Features

### package constraint

* Representations for inter-table constraints.

### package require

* Predicates allowing easier detection of unexpected results from SELECTS, e.g. when the result set size is not exactly one.

### package dialect

* SQL dialects for SQLite, MySQL, PostgreSQL and its pgx variant. This provides some conditional SQL generation and
 also 

### package where

* Fluent construction of WHERE and HAVING clauses: this is now [github.com/rickb777/sqlapi](https://github
.com/rickb777/where). This also provides control over identifier quoting, e.g. ANSI SQL, back-ticks, etc.
 
## Install

Install with this command:

```
go get github.com/rickb777/sqlapi
```

Please continue reading about [sqlgen2](https://github.com/rickb777/sqlgen2/tree/master/README.md).
