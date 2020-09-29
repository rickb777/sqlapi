// Package logadapter provides a logger that writes to a log.Logger log.
package logadapter

import (
	"context"
	"fmt"
	"github.com/jackc/pgx/v4"
	"log"
)

type Logger struct {
	l     *log.Logger
	level pgx.LogLevel
}

func NewLogger(l *log.Logger) *Logger {
	return &Logger{l: l}
}

func (l *Logger) Log(ctx context.Context, level pgx.LogLevel, msg string, data map[string]interface{}) {
	if level >= l.level {
		m := fmt.Sprintf("%v", data)
		l.l.Printf("%s %s", msg, m[3:])
	}
}
