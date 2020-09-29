#!/bin/bash -e
# Setup
# -----
# This script can run some of the tests against real PostgreSQL. Therefore,
# it is necessary to create testuser accounts in each one beforehand.
# These all assume the existence of either
#   postgres:postgres@/postgres
#   testuser:TestPasswd.9.9.9@/test

cd "$(dirname $0)"

PATH=$HOME/go/bin:$PATH

export PGHOST=localhost

if [[ $1 = "-v" ]]; then
  V=-v
  shift
fi

if [[ $1 = "travis" ]]; then
  export PGDATABASE='postgres'
  export PGUSER='postgres'
  export PGPASSWORD=''
elif [[ -z $PGUSER ]]; then
  export PGDATABASE='test'
  export PGUSER='testuser'
  export PGPASSWORD='TestPasswd.9.9.9'
fi

echo
echo "PGX (no quotes)...."
go clean -testcache ||:
PGQUOTE=none go test $V ./...

echo
echo "PGX (ANSI)...."
go clean -testcache ||:
PGQUOTE=ansi go test $V ./...

