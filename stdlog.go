package sqlapi

import (
	"io"
	"log"
)

type StdLog interface {
	Printf(format string, v ...interface{})
}

// outable interface for loggers that allow setting the output writer
type outable interface {
	SetOutput(out io.Writer)
}

var _ StdLog = new(log.Logger)
var _ outable = new(log.Logger)
