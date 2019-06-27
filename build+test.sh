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
[ -f schema/type_set.go ]  || runtemplate -tpl simple/set.tpl  -output schema/type_set.go  Type=Type   Comparable:true Ordered:false Numeric:false
[ -f util/int64_set.go ]   || runtemplate -tpl simple/set.tpl  -output util/int64_set.go   Type=int64  Comparable:true Ordered:true  Numeric:true
[ -f util/string_list.go ] || runtemplate -tpl simple/list.tpl -output util/string_list.go Type=string Comparable:true Ordered:true  Numeric:false
[ -f util/string_set.go ]  || runtemplate -tpl simple/set.tpl  -output util/string_set.go  Type=string Comparable:true Ordered:true  Numeric:false

if [ ! -f util/string_any_map.go ]; then
  runtemplate -v -tpl simple/map.tpl  -output util/string_any_map.xx Key=string Type=any
  sed 's#any#interface\{\}#g#' < util/string_any_map.xx > util/string_any_map.go
  rm util/string_any_map.xx
fi

[ -f support/functions_gen.go ] || ./support/functions.sh
#[ -f database_gen.go ] || ./database.sh

### Build Phase 1 ###

./version.sh

gofmt -l -w *.go */*.go
go vet ./...
go install ./...
go test ./...

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
