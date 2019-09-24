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

unset GOPATH GO_DRIVER GO_DSN GO_QUOTER

# these names are as used by PostgreSQL, complementing the DSN
export PGHOST='localhost'
export PGDATABASE='test'
export PGUSER='testuser'
export PGPASSWORD='TestPasswd.9.9.9'

function startMysqlDocker
{
  if [ ! -f .test-mysql-$$ ]; then
    touch .test-mysql-$$
    echo "Using docker to provide MySQL means"
    echo " 1. if you have one, you must stop your local MySQL instance"
    echo " 2. docker pull mysql:5.7"
    echo
    echo "starting docker test-postgres..."
    docker run --name test-mysql -e MYSQL_DATABASE=$PGDATABASE -e MYSQL_USER=$PGUSER -e MYSQL_PASSWORD="$PGPASSWORD" -e MYSQL_ROOT_PASSWORD=mysql -p 127.0.0.1:3306:3306/tcp -d mysql:5.7
    sleep 1
  fi
}

function startPostgresDocker
{
  if [ ! -f .test-postgres-$$ ]; then
    touch .test-postgres-$$
    echo "Using docker to provide PostgreSQL means"
    echo " 1. if you have one, you must stop your local PostgreSQL instance"
    echo " 2. docker pull postgres:11-alpine"
    echo
    echo "starting docker test-postgres..."
    docker run --name test-postgres -e POSTGRES_DB=$PGDATABASE -e POSTGRES_USER=$PGUSER -e POSTGRES_PASSWORD="$PGPASSWORD" -p 127.0.0.1:5432:5432/tcp -d postgres:11-alpine
    sleep 1
  fi
}

function stopDockers
{
  if [ -f .test-mysql-$$ ]; then
    rm .test-mysql-$$
    docker rm -f test-mysql #>/dev/null 2>&1
  fi

  if [ -f .test-postgres-$$ ]; then
    rm .test-postgres-$$
    docker rm -f test-postgres >/dev/null 2>&1
  fi
}

trap stopDockers EXIT

if [ "$1" = "-v" ]; then
  V=-v
  shift
fi

DBS=$*
if [ "$1" = "all" ]; then
  DBS="sqlite mysql postgres pgx"
fi

PACKAGES=". ./constraint ./dialect"

for db in $DBS; do
  echo
  go clean -testcache ||:

  case $db in
    mysql)
      startMysqlDocker
      echo
      echo "MySQL...."
      GO_DRIVER=mysql GO_DSN="testuser:TestPasswd.9.9.9@/test" go test $V ./constraint
      ;;

    postgres)
      startPostgresDocker
      echo
      echo "PostgreSQL (no quotes)...."
      GO_DRIVER=postgres GO_DSN="postgres://$PGUSER:$PGPASSWORD@/$PGDATABASE?sslmode=disable" GO_QUOTER=none go test $V $PACKAGES
      echo
      echo "PostgreSQL (ANSI)...."
      GO_DRIVER=postgres GO_DSN="postgres://$PGUSER:$PGPASSWORD@/$PGDATABASE?sslmode=disable" GO_QUOTER=ansi go test $V $PACKAGES
      ;;

    pgx)
      startPostgresDocker
      echo
      echo "PGX (no quotes)...."
      GO_DRIVER=pgx GO_DSN="postgres://$PGUSER:$PGPASSWORD@/$PGDATABASE?sslmode=disable" GO_QUOTER=none go test $V $PACKAGES
      echo
      echo "PGX (ANSI)...."
      GO_DRIVER=pgx GO_DSN="postgres://$PGUSER:$PGPASSWORD@/$PGDATABASE?sslmode=disable" GO_QUOTER=ansi go test $V $PACKAGES
      echo
      echo "PGXAPI (ANSI)...."
      GO_QUOTER=ansi go test ./pgxapi/...
      ;;

    sqlite)
      unset GO_DRIVER GO_DSN
      echo
      echo "SQLite3 (no quotes)..."
      GO_QUOTER=none go test $V $PACKAGES
      echo
      echo "SQLite3 (ANSI)..."
      GO_QUOTER=ansi go test $V $PACKAGES
      ;;

    *)
      echo "$db: unrecognised; must be sqlite, mysql, or postgres. Use 'all' for all of these."
      exit 1
      ;;
  esac
done
