package test

import (
	"context"
	"fmt"
	"strings"

	"github.com/bobg/go-generics/v3/maps"
	"github.com/bobg/go-generics/v3/slices"
	"github.com/jackc/pgx/v5/tracelog"
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

func (r *StubLogger) Log(ctx context.Context, level tracelog.LogLevel, msg string, data map[string]any) {
	m := maps.Clone(data)
	skeys := maps.Keys(m)
	slices.Sort(skeys)
	buf := &strings.Builder{}
	fmt.Fprintf(buf, "%-5s %s [", level, msg)
	comma := ""
	for _, k := range skeys {
		v := data[k]
		fmt.Fprintf(buf, "%s%v=%v", comma, k, v)
		comma = ", "
	}
	buf.WriteString("]")
	s := buf.String()
	r.Logged = append(r.Logged, s)
	if r.Testing != nil {
		r.Testing.Log(s)
	}

}
