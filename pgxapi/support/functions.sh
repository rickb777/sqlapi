#!/bin/bash -e
cd "$(dirname $0)"
# Generates additional functions from templates.

cat <<PREAMBLE > functions_gen.go
package support

import (
	"context"
	"database/sql"
	"github.com/rickb777/sqlapi/pgxapi"
	"github.com/rickb777/sqlapi/require"
	"github.com/rickb777/where"
)
PREAMBLE

D=../../support
runtemplate -tpl $D/functions.tpl -o - SqlApi:pgxapi Type=string  NT:String  >> functions_gen.go
runtemplate -tpl $D/functions.tpl -o - SqlApi:pgxapi Type=int     NT:Int64   >> functions_gen.go
runtemplate -tpl $D/functions.tpl -o - SqlApi:pgxapi Type=int64   NT:Int64   >> functions_gen.go
runtemplate -tpl $D/functions.tpl -o - SqlApi:pgxapi Type=int32   NT:Int64   >> functions_gen.go
runtemplate -tpl $D/functions.tpl -o - SqlApi:pgxapi Type=int16   NT:Int64   >> functions_gen.go
runtemplate -tpl $D/functions.tpl -o - SqlApi:pgxapi Type=int8    NT:Int64   >> functions_gen.go
runtemplate -tpl $D/functions.tpl -o - SqlApi:pgxapi Type=uint    NT:Int64   >> functions_gen.go
runtemplate -tpl $D/functions.tpl -o - SqlApi:pgxapi Type=uint64  NT:Int64   >> functions_gen.go
runtemplate -tpl $D/functions.tpl -o - SqlApi:pgxapi Type=uint32  NT:Int64   >> functions_gen.go
runtemplate -tpl $D/functions.tpl -o - SqlApi:pgxapi Type=uint16  NT:Int64   >> functions_gen.go
runtemplate -tpl $D/functions.tpl -o - SqlApi:pgxapi Type=uint8   NT:Int64   >> functions_gen.go
runtemplate -tpl $D/functions.tpl -o - SqlApi:pgxapi Type=float64 NT:Float64 >> functions_gen.go
runtemplate -tpl $D/functions.tpl -o - SqlApi:pgxapi Type=float32 NT:Float64 >> functions_gen.go
