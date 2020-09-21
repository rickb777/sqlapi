package sqlapi

import (
	"io"
	"log"
	"strings"
	"sync/atomic"
	"time"
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

type toggleLogger struct {
	lgr     StdLog
	enabled int32
}

func NewLogger(lgr StdLog) Logger {
	if lgr == nil {
		return &toggleLogger{}
	}
	return &toggleLogger{lgr: lgr, enabled: 1}
}

func (lgr *toggleLogger) Log(msg string, v ...interface{}) {
	if lgr.loggingEnabled() {
		lgr.lgr.Printf(msg, v...)
	}
}

// LogT emits a log event, supporting an elapsed-time calculation and providing an easier
// way to supply data parameters
func (lgr *toggleLogger) LogT(msg string, startTime *time.Time, data ...interface{}) {
	if lgr.loggingEnabled() {
		if startTime != nil {
			took := time.Now().Sub(*startTime)
			data = append(data, took)
			msg += " took %s"
		}
		lgr.lgr.Printf(msg, data...)
	}
}

// LogIfError writes error info to the logger, if both the logger and the error are non-nil.
// It returns the error.
func (lgr *toggleLogger) LogIfError(err error) error {
	if err != nil {
		lgr.LogT("Error: %v", nil, err)
	}
	return err
}

// LogError writes error info to the logger, if the logger is not nil. It returns the error.
func (lgr *toggleLogger) LogError(err error) error {
	lgr.LogT("Error: %v", nil, err)
	return err
}

func (lgr *toggleLogger) SetOutput(w io.Writer) {
	if lgr.lgr != nil {
		if o, ok := lgr.lgr.(outable); ok {
			o.SetOutput(w)
		}
	}
}

func (lgr *toggleLogger) TraceLogging(on bool) {
	if on && lgr.lgr != nil {
		atomic.StoreInt32(&lgr.enabled, 1)
	} else {
		atomic.StoreInt32(&lgr.enabled, 0)
	}
}

func (lgr *toggleLogger) loggingEnabled() bool {
	return atomic.LoadInt32(&lgr.enabled) != 0
}

// LogQuery writes query info to the logger, if it is not nil.
func (lgr *toggleLogger) LogQuery(query string, args ...interface{}) {
	if lgr.loggingEnabled() {
		query = strings.TrimSpace(query)
		if len(args) > 0 {
			ss := make([]interface{}, len(args))
			for i, v := range args {
				ss[i] = derefArg(v)
			}
			lgr.lgr.Printf("%s %v\n", query, ss)
		} else {
			lgr.lgr.Printf("%s\n", query)
		}
	}
}

func derefArg(arg interface{}) interface{} {
	switch v := arg.(type) {
	case *int:
		return *v
	case *int8:
		return *v
	case *int16:
		return *v
	case *int32:
		return *v
	case *int64:
		return *v
	case *uint:
		return *v
	case *uint8:
		return *v
	case *uint16:
		return *v
	case *uint32:
		return *v
	case *uint64:
		return *v
	case *float32:
		return *v
	case *float64:
		return *v
	case *bool:
		return *v
	case *string:
		return *v
	}
	return arg
}
