package test

import (
	"fmt"

	"github.com/rickb777/sqlapi"
)

// StubLogger provides a non-functioning sqlapi.StdLog for testing purposes. It
// captures all the lines logged.
type StubLogger struct {
	Logged []string
}

var _ sqlapi.StdLog = new(StubLogger)

func (r *StubLogger) Printf(format string, v ...interface{}) {
	r.Logged = append(r.Logged, fmt.Sprintf(format, v...))
}
