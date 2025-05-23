# sqlapi

[![GoDoc](https://img.shields.io/badge/api-Godoc-blue.svg)](https://pkg.go.dev/github.com/rickb777/sqlapi)
[![Go Report Card](https://goreportcard.com/badge/github.com/rickb777/sqlapi)](https://goreportcard.com/report/github.com/rickb777/sqlapi)
[![Issues](https://img.shields.io/github/issues/rickb777/sqlapi.svg)](https://github.com/rickb777/sqlapi/issues)

**sqlgen** generates SQL statements and database helper functions from your Go structs. It can be used in
place of a simple ORM or hand-written SQL. **sqlapi** (this package) supports the generated code.

Currently, support is included for **MySQL**, **PostgreSQL** and **SQLite**. Other dialects can be added relatively easy - send a Pull Request!

## Features

### package constraint

* Representations for inter-table constraints.

### package require

* Predicates allowing easier detection of unexpected results from SELECTS, e.g. when the result set size is not exactly one.

### package dialect

* SQL dialects for SQLite, MySQL, PostgreSQL and its pgx variant. This provides some conditional SQL generation and
 also 

## Install

Install with this command:

```
go get github.com/rickb777/sqlapi
```
