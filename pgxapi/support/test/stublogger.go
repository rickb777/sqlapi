package test

import (
	"fmt"
	"github.com/jackc/pgx"
)

// StubLogger provides a non-functioning pgx.Logger for testing purposes. It
// captures all the lines logged.
type StubLogger struct {
	Logged []string
}

var _ pgx.Logger = new(StubLogger)

func (r *StubLogger) Log(level pgx.LogLevel, msg string, data map[string]interface{}) {
	r.Logged = append(r.Logged, fmt.Sprintf("%s %v", msg, data))
}
