package support

import (
	"fmt"
	"io"
	"regexp"

	"github.com/rickb777/collection"
	"github.com/rickb777/sqlapi"
	"github.com/rickb777/sqlapi/dialect"
)

type StubDatabase struct {
	execer StubExecer
	stdLog *stubLogger
}

func (*StubDatabase) DB() sqlapi.SqlDB {
	panic("implement me")
}

func (*StubDatabase) Dialect() dialect.Dialect {
	panic("implement me")
}

func (d *StubDatabase) Logger() sqlapi.Logger {
	return sqlapi.NewLogger(d.stdLog)
}

func (*StubDatabase) Wrapper() interface{} {
	panic("implement me")
}

func (*StubDatabase) ListTables(re *regexp.Regexp) (collection.StringList, error) {
	panic("implement me")
}

//-------------------------------------------------------------------------------------------------

type stubLogger struct {
	logged []string
}

func (r *stubLogger) Printf(format string, v ...interface{}) {
	r.logged = append(r.logged, fmt.Sprintf(format, v...))
}

func (r *stubLogger) SetOutput(w io.Writer) {}
