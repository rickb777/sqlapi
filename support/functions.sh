#!/bin/bash -e
cd "$(dirname $0)"
# Generates additional functions from templates.

cat <<PREAMBLE > functions_gen.go
package support

import (
	"database/sql"
	"github.com/rickb777/sqlapi"
	"github.com/rickb777/sqlapi/require"
	"github.com/rickb777/where"
)
PREAMBLE

runtemplate -tpl functions.tpl -o - SqlApi:sqlapi Type=string  NT:String  >> functions_gen.go
runtemplate -tpl functions.tpl -o - SqlApi:sqlapi Type=int     NT:Int64   >> functions_gen.go
runtemplate -tpl functions.tpl -o - SqlApi:sqlapi Type=int64   NT:Int64   >> functions_gen.go
runtemplate -tpl functions.tpl -o - SqlApi:sqlapi Type=int32   NT:Int64   >> functions_gen.go
runtemplate -tpl functions.tpl -o - SqlApi:sqlapi Type=int16   NT:Int64   >> functions_gen.go
runtemplate -tpl functions.tpl -o - SqlApi:sqlapi Type=int8    NT:Int64   >> functions_gen.go
runtemplate -tpl functions.tpl -o - SqlApi:sqlapi Type=uint    NT:Int64   >> functions_gen.go
runtemplate -tpl functions.tpl -o - SqlApi:sqlapi Type=uint64  NT:Int64   >> functions_gen.go
runtemplate -tpl functions.tpl -o - SqlApi:sqlapi Type=uint32  NT:Int64   >> functions_gen.go
runtemplate -tpl functions.tpl -o - SqlApi:sqlapi Type=uint16  NT:Int64   >> functions_gen.go
runtemplate -tpl functions.tpl -o - SqlApi:sqlapi Type=uint8   NT:Int64   >> functions_gen.go
runtemplate -tpl functions.tpl -o - SqlApi:sqlapi Type=float64 NT:Float64 >> functions_gen.go
runtemplate -tpl functions.tpl -o - SqlApi:sqlapi Type=float32 NT:Float64 >> functions_gen.go
