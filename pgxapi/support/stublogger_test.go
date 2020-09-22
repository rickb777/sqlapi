package support

import (
	"fmt"
	"github.com/jackc/pgx"
)

type StubLogger struct {
	Logged []string
}

func (r *StubLogger) Log(level pgx.LogLevel, msg string, data map[string]interface{}) {
	r.Logged = append(r.Logged, fmt.Sprintf("%s %v", msg, data))
}

//-------------------------------------------------------------------------------------------------

type stubLogger struct {
	logged []string
}

func (r *stubLogger) Log(level pgx.LogLevel, msg string, data map[string]interface{}) {
	r.logged = append(r.logged, fmt.Sprintf("%s %v", msg, data))
}
