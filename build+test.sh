#!/bin/bash -e
cd "$(dirname $0)"
PATH=$HOME/go/bin:$PATH
unset GOPATH
export GO111MODULE=on

function announce
{
  echo
  echo $@
}

function v
{
  announce $@
  $@
}

if ! type -p goveralls; then
  v go install github.com/mattn/goveralls
fi

if ! type -p shadow; then
  v go install golang.org/x/tools/go/analysis/passes/shadow/cmd/shadow
fi

if ! type -p goreturns; then
  v go install github.com/sqs/goreturns
fi

go mod download

# delete artefacts from previous build (if any)
mkdir -p reports
rm -f reports/*.out reports/*.html */*.txt demo/*_sql.go

### Collection Types ###
# these generated files hardly ever need to change (see github.com/rickb777/runtemplate to do so)
if [[ $1 != "travis" ]]; then
  [ -f schema/type_set.go ]  || runtemplate -tpl simple/set.tpl  -output schema/type_set.go Type=Type Comparable:true Ordered:false Numeric:false
  [ support/functions.tpl -ot support/functions_gen.go ] || rm -vf support/functions_gen.go
  [ support/functions.tpl -ot pgxapi/support/functions_gen.go ] || rm -vf pgxapi/support/functions_gen.go
  [ -f support/functions_gen.go ] || ./support/functions.sh
  [ -f pgxapi/support/functions_gen.go ] || ./pgxapi/support/functions.sh
fi

### Build Phase 1 ###

./version.sh

v goreturns -l -w *.go */*.go

v go vet ./...

v shadow ./...

v go install ./...

v ./test.sh sqlite

### Build Phase 2 ###
export PGDATABASE='test'
export PGUSER='testuser'
export PGPASSWORD='TestPasswd.9.9.9'

rm -f reports/*

# sqlapi test coverage
echo .
go test . -covermode=count -coverprofile=reports/sqlapi.out .
go tool cover -func=reports/sqlapi.out
#[ -z "$COVERALLS_TOKEN" ] || goveralls -coverprofile=reports/sqlapi.out -service=travis-ci -repotoken $COVERALLS_TOKEN || echo "Push to coveralls failed"

for d in constraint require schema support types; do
  announce ./$d
  go test -covermode=count -coverprofile=reports/$d.out ./$d
  go tool cover -func=reports/$d.out
  #[ -z "$COVERALLS_TOKEN" ] || goveralls -coverprofile=reports/$d.out -service=travis-ci -repotoken $COVERALLS_TOKEN || echo "Push to coveralls failed"
done

# pgxapi sub-package test coverage
for d in constraint support; do
  announce ./pgxapi/$d
  go test -covermode=count -coverprofile=reports/pgxapi-$d.out ./pgxapi/$d
  go tool cover -func=reports/pgxapi-$d.out
  #[ -z "$COVERALLS_TOKEN" ] || goveralls -coverprofile=reports/pgxapi-$d.out -service=travis-ci -repotoken $COVERALLS_TOKEN || echo "Push to coveralls failed"
done

./pgxapi/pgtest.sh $1

echo
echo "Now start MySQL and PostgreSQL, then run './test.sh all'"
