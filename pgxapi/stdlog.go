package pgxapi

import (
	"context"
	"fmt"
	"io"
	"log"

	"github.com/jackc/pgx/v5/tracelog"
)

type StdLog interface {
	Println(v ...interface{})
}

// outable interface for loggers that allow setting the output writer
type outable interface {
	SetOutput(out io.Writer)
}

var _ StdLog = new(log.Logger)
var _ outable = new(log.Logger)

type stdLogAdapter struct {
	std StdLog
}

func (s stdLogAdapter) Log(_ context.Context, level tracelog.LogLevel, msg string, data map[string]interface{}) {
	ss := make([]interface{}, 2, 2+len(data))
	ss[0] = fmt.Sprintf("%-5s", level)
	ss[1] = msg
	for k, v := range data {
		ss = append(ss, fmt.Sprintf("%s=%v", k, v))
	}
	s.std.Println(ss...)
}

func (s stdLogAdapter) SetOutput(w io.Writer) {
	if o, ok := s.std.(outable); ok {
		o.SetOutput(w)
	}
}
