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
[ -f schema/type_set.go ]  || runtemplate -tpl simple/set.tpl  -output schema/type_set.go Type=Type Comparable:true Ordered:false Numeric:false
[ -f support/functions_gen.go ] || ./support/functions.sh
[ -f pgxapi/support/functions_gen.go ] || ./pgxapi/support/functions.sh

### Build Phase 1 ###

./version.sh

export PGDATABASE='test'
export PGUSER='testuser'
export PGPASSWORD='TestPasswd.9.9.9'

v goreturns -l -w *.go */*.go

v go vet ./...

v shadow ./...

v go install ./...

v ./test.sh $1

### Build Phase 2 ###

for d in constraint require schema types; do
  announce ./$d
  go test ./$d -covermode=count -coverprofile=reports/$d.out ./$d
  go tool cover -html=reports/$d.out -o reports/$d.html
  [ -z "$COVERALLS_TOKEN" ] || goveralls -coverprofile=reports/$d.out -service=travis-ci -repotoken $COVERALLS_TOKEN || echo "Push to coveralls failed"
done

echo .
go test . -covermode=count -coverprofile=reports/dot.out .
go tool cover -func=reports/dot.out
[ -z "$COVERALLS_TOKEN" ] || goveralls -coverprofile=reports/dot.out -service=travis-ci -repotoken $COVERALLS_TOKEN || echo "Push to coveralls failed"

echo
echo "Now start MySQL and PostgreSQL, then run './test.sh all'"
