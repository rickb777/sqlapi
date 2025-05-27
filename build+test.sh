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

go mod download

# delete artefacts from previous build (if any)
mkdir -p reports
rm -f reports/*.out reports/*.html */*.txt demo/*_sql.go

### Collection Types ###
# these generated files hardly ever need to change (see github.com/rickb777/runtemplate to do so)
[ -f schema/type_set.go ]  || runtemplate -tpl simple/set.tpl  -output schema/type_set.go Type=Type Comparable:true Ordered:false Numeric:false
[ support/functions.tpl -ot support/functions_gen.go ] || rm -vf support/functions_gen.go
[ support/functions.tpl -ot pgxapi/support/functions_gen.go ] || rm -vf pgxapi/support/functions_gen.go
[ -f support/functions_gen.go ] || ./support/functions.sh
[ -f pgxapi/support/functions_gen.go ] || ./pgxapi/support/functions.sh

### Build Phase 1 ###

v gofmt -l -w *.go */*.go

v go vet ./...

v go install ./...

v ./test.sh -v sqlite

### Check Docker is available

if ! type -p docker; then
  echo
  echo "*** Docker is not installed. The remaining tests will be skipped, which is inconclusive."
  echo
  exit 0
fi

### Build Phase 2 ###
echo
echo "========== Build Phase 2 =========="

rm -f reports/*

# sqlapi test coverage
echo .
go test -covermode=count -coverprofile=reports/sqlapi.out .
go tool cover -func=reports/sqlapi.out

for d in constraint require schema support types; do
  announce ./$d
  go test -covermode=count -coverprofile=reports/$d.out ./$d
  go tool cover -func=reports/$d.out
done

# pgxapi sub-package test coverage
for d in constraint support; do
  announce ./pgxapi/$d
  go test -covermode=count -coverprofile=reports/pgxapi-$d.out ./pgxapi/$d
  go tool cover -func=reports/pgxapi-$d.out
done

echo "./pgxapi/pgtest.sh $1"
./pgxapi/pgtest.sh $1

echo
echo "Now start MySQL and PostgreSQL, then run './test.sh all'"
