package test

import (
	"context"
	"fmt"
	"strings"

	"github.com/jackc/pgx/v5/tracelog"
	"github.com/rickb777/collection"
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

var _ tracelog.Logger = new(StubLogger)

func (r *StubLogger) Log(_ context.Context, level tracelog.LogLevel, msg string, data map[string]interface{}) {
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
