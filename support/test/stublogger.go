package test

import (
	"fmt"
	"github.com/jackc/pgx"
	"github.com/rickb777/collection"
	"strings"
)

// TestingLogger interface defines the subset of testing.TB methods used by this adapter.
type TestingLogger interface {
	Log(args ...interface{})
}

// StubLogger provides a testingadapter.TestingLogger that captures logged information
// and optionally plays it through a child logger too.
type StubLogger struct {
	Testing TestingLogger
	Logged  []string
}

var _ pgx.Logger = new(StubLogger)

func (r *StubLogger) Log(level pgx.LogLevel, msg string, data map[string]interface{}) {
	m := collection.StringAnyMap(data)
	args := m.OrderedSlice(m.Keys().Sorted()).MkString4("[", ", ", "]", "=")
	buf := &strings.Builder{}
	fmt.Fprintf(buf, "%-5s %s %s", level, msg, args)
	s := buf.String()
	r.Logged = append(r.Logged, s)
	if r.Testing != nil {
		r.Testing.Log(s)
	}
}
