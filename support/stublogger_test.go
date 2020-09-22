package support

import (
	"fmt"
	"io"
)

type StubLogger struct {
	Logged []string
}

func (r *StubLogger) Printf(format string, v ...interface{}) {
	r.Logged = append(r.Logged, fmt.Sprintf(format, v...))
}

//-------------------------------------------------------------------------------------------------

type stubLogger struct {
	logged []string
}

func (r *stubLogger) Printf(format string, v ...interface{}) {
	r.logged = append(r.logged, fmt.Sprintf(format, v...))
}

func (r *stubLogger) SetOutput(w io.Writer) {}
