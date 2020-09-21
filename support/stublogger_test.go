package support

import (
	"fmt"
)

type StubLogger struct {
	Logged []string
}

func (r *StubLogger) Printf(format string, v ...interface{}) {
	r.Logged = append(r.Logged, fmt.Sprintf(format, v...))
}
