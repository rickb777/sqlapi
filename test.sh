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

PATH=$HOME/go/bin:$PATH

unset GOPATH DB_DRIVER GO_DSN PGQUOTE
export PGHOST=localhost

#
# accommodate different ways of running this script, including Travis
#
if [[ -n $PGUSER ]]; then
  DBUSER=$PGUSER
  DBPASS=$PASSWORD
fi

if [[ -z $DBUSER ]]; then
  export DBUSER=testuser
  export DBPASS=TestPasswd.9.9.9
fi

if [[ -z $PGUSER ]]; then
  export PGUSER=$DBUSER
  export PGPASSWORD=$DBPASS
  export PGDATABASE=test
fi

if [[ $1 = "-v" ]]; then
  V=-v
  shift
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
  case $db in
    mysql)
      echo "MySQL...."
      go clean -testcache ||:
      DB_DRIVER=mysql DB_URL=$DBUSER:$DBPASS@/test go test $V $PACKAGES
      ;;

    postgres)
      echo "PostgreSQL (no quotes)...."
      go clean -testcache ||:
      DB_DRIVER=postgres DB_URL="postgres://$DBUSER:$DBPASS@/test" PGQUOTE=none go test $V $PACKAGES
      echo
      echo "PostgreSQL (ANSI)...."
      go clean -testcache ||:
      DB_DRIVER=postgres DB_URL="postgres://$DBUSER:$DBPASS@/test" PGQUOTE=ansi go test $V $PACKAGES
      ;;

    pgx)
      echo "PGX (no quotes)...."
      go clean -testcache ||:
      DB_DRIVER=pgx DB_URL="postgres://$DBUSER:$DBPASS@/test" PGQUOTE=none go test $V $PACKAGES
      echo
      echo "PGX (ANSI)...."
      go clean -testcache ||:
      DB_DRIVER=pgx DB_URL="postgres://$DBUSER:$DBPASS@/test" PGQUOTE=ansi go test $V $PACKAGES
      echo
      echo "PGXAPI (ANSI)...."
      go clean -testcache ||:
      PGQUOTE=ansi go test $V ./pgxapi/...
      ;;

    sqlite)
      unset DB_DRIVER DB_URL
      echo "SQLite3 (no quotes)..."
      go clean -testcache ||:
      DB_QUOTE=none go test $V $PACKAGES
      echo
      echo "SQLite3 (ANSI)..."
      go clean -testcache ||:
      DB_QUOTE=ansi go test $V $PACKAGES
      ;;

    *)
      echo "$db: unrecognised; must be sqlite, mysql, or postgres. Use 'all' for all of these."
      exit 1
      ;;
  esac
done
