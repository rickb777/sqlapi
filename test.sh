#!/bin/bash -e
# Setup
# -----
# This script can run some of the tests against real databases. Therefore,
# it is necessary to create testuser accounts in each one beforehand.
# These all assume the existence of
#   testuser:TestPasswd.9.9.9@/test
#
# e.g.
# create user testuser@localhost identified by 'TestPasswd.9.9.9';
# create database test;
# grant all on test.* to testuser@localhost;

cd "$(dirname $0)"

PATH=$PWD/..:$HOME/go/bin:$PATH

unset GOPATH GO_DRIVER GO_DSN GO_QUOTER
export PGHOST=localhost

if [[ -z $DBUSER ]]; then
  DBUSER=testuser
  DBPASS=TestPasswd.9.9.9
fi

DBS=$*
if [[ -z $1 ]]; then
  DBS="sqlite"
elif [[ $1 = "all" ]]; then
  DBS="sqlite mysql postgres pgx"
fi

PACKAGES=". ./constraint ./dialect"

for db in $DBS; do
  echo
  go clean -testcache ||:

  case $db in
    mysql)
      echo
      echo "MySQL...."
      GO_DRIVER=mysql GO_DSN=$DBUSER:$DBPASS@/test go test -v $PACKAGES
      ;;

    postgres)
      echo
      echo "PostgreSQL (no quotes)...."
      GO_DRIVER=postgres GO_DSN="postgres://$DBUSER:$DBPASS@/test" GO_QUOTER=none go test -v $PACKAGES
      echo
      echo "PostgreSQL (ANSI)...."
      GO_DRIVER=postgres GO_DSN="postgres://$DBUSER:$DBPASS@/test" GO_QUOTER=ansi go test -v $PACKAGES
      ;;

    pgx)
      echo
      echo "PGX (no quotes)...."
      GO_DRIVER=pgx GO_DSN="postgres://$DBUSER:$DBPASS@/test" GO_QUOTER=none go test -v $PACKAGES
      echo
      echo "PGX (ANSI)...."
      GO_DRIVER=pgx GO_DSN="postgres://$DBUSER:$DBPASS@/test" GO_QUOTER=ansi go test -v $PACKAGES
      echo
      echo "PGXAPI (ANSI)...."
      PGUSER=$DBUSER PGPASSWORD=$DBPASS PGDATABASE=test GO_QUOTER=ansi go test -v ./pgxapi/...
      ;;

    sqlite)
      unset GO_DRIVER GO_DSN
      echo
      echo "SQLite3 (no quotes)..."
      GO_QUOTER=none go test -v $PACKAGES
      echo
      echo "SQLite3 (ANSI)..."
      GO_QUOTER=ansi go test -v $PACKAGES
      ;;

    *)
      echo "$db: unrecognised; must be sqlite, mysql, or postgres. Use 'all' for all of these."
      exit 1
      ;;
  esac
done
