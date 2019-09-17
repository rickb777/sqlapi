#!/bin/bash -e
cd $(dirname $0)

function announce
{
  echo
  echo $@
}

PATH=$HOME/go/bin:$PATH

unset GOPATH

go mod download

# delete artefacts from previous build (if any)
mkdir -p reports
rm -f reports/*.out reports/*.html */*.txt demo/*_sql.go

### Collection Types ###
# these generated files hardly ever need to change (see github.com/rickb777/runtemplate to do so)
[ -f schema/type_set.go ]  || runtemplate -tpl simple/set.tpl  -output schema/type_set.go Type=Type Comparable:true Ordered:false Numeric:false
[ -f support/functions_gen.go ] || ./support/functions.sh
[ -f pgxapi/support/functions_gen.go ] || ./pgxapi/support/functions.sh
#[ -f database_gen.go ] || ./database.sh

### Build Phase 1 ###

./version.sh

export PGDATABASE='test'
export PGUSER='testuser'
export PGPASSWORD='TestPasswd.9.9.9'

gofmt -l -w *.go */*.go
go vet ./...
go install ./...
./test.sh all

### Build Phase 2 ###

for d in constraint require schema types; do
  announce ./$d
  go test $1 -covermode=count -coverprofile=reports/$d.out ./$d
  go tool cover -html=reports/$d.out -o reports/$d.html
  [ -z "$COVERALLS_TOKEN" ] || goveralls -coverprofile=reports/$d.out -service=travis-ci -repotoken $COVERALLS_TOKEN
done

echo .
go test . -covermode=count -coverprofile=reports/dot.out .
go tool cover -func=reports/dot.out
[ -z "$COVERALLS_TOKEN" ] || goveralls -coverprofile=$d.out -service=travis-ci -repotoken $COVERALLS_TOKEN

#git checkout util/version.go
