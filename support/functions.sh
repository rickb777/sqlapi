#!/bin/bash -e
cd $(dirname $0)
# Generates additional functions from templates.

cat <<PREAMBLE > functions_gen.go
package support

import (
	"database/sql"
	"github.com/rickb777/sqlapi"
	"github.com/rickb777/sqlapi/require"
	"github.com/rickb777/sqlapi/where"
)
PREAMBLE

runtemplate -tpl functions.tpl -o - Type=string  NT:String  >> functions_gen.go
runtemplate -tpl functions.tpl -o - Type=int     NT:Int64   >> functions_gen.go
runtemplate -tpl functions.tpl -o - Type=int64   NT:Int64   >> functions_gen.go
runtemplate -tpl functions.tpl -o - Type=int32   NT:Int64   >> functions_gen.go
runtemplate -tpl functions.tpl -o - Type=int16   NT:Int64   >> functions_gen.go
runtemplate -tpl functions.tpl -o - Type=int8    NT:Int64   >> functions_gen.go
runtemplate -tpl functions.tpl -o - Type=uint    NT:Int64   >> functions_gen.go
runtemplate -tpl functions.tpl -o - Type=uint64  NT:Int64   >> functions_gen.go
runtemplate -tpl functions.tpl -o - Type=uint32  NT:Int64   >> functions_gen.go
runtemplate -tpl functions.tpl -o - Type=uint16  NT:Int64   >> functions_gen.go
runtemplate -tpl functions.tpl -o - Type=uint8   NT:Int64   >> functions_gen.go
runtemplate -tpl functions.tpl -o - Type=float64 NT:Float64 >> functions_gen.go
runtemplate -tpl functions.tpl -o - Type=float32 NT:Float64 >> functions_gen.go
