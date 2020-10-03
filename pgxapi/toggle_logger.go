package pgxapi

import (
	"context"
	"fmt"
	"io"
	"reflect"
	"sync/atomic"
	"time"

	"github.com/jackc/pgx/v4"
)

type toggleLogger struct {
	lgr     pgx.Logger
	enabled int32
}

func NewStdLogger(lgr StdLog) Logger {
	return NewLogger(stdLogAdapter{lgr})
}

func NewLogger(lgr pgx.Logger) Logger {
	if lgr == nil {
		return &toggleLogger{}
	}
	// because StdLog is an interface, it might be not nil yet hold a nil pointer
	value := reflect.ValueOf(lgr)
	if value.Kind() == reflect.Ptr && value.IsNil() {
		return &toggleLogger{}
	}
	return &toggleLogger{lgr: lgr, enabled: 1}
}

func (lgr *toggleLogger) Log(ctx context.Context, level pgx.LogLevel, msg string, data map[string]interface{}) {
	if lgr.loggingEnabled() {
		lgr.lgr.Log(ctx, level, msg, data)
	}
}

// LogT emits a log event, supporting an elapsed-time calculation and providing an easier
// way to supply data parameters as name,value pairs.
func (lgr *toggleLogger) LogT(ctx context.Context, level pgx.LogLevel, msg string, startTime *time.Time, data ...interface{}) {
	if lgr.loggingEnabled() {
		m := make(map[string]interface{})
		if startTime != nil {
			took := time.Now().Sub(*startTime)
			m["took"] = took
		}
		for i := 1; i < len(data); i += 2 {
			k := data[i-1].(string)
			v := data[i]
			m[k] = v
		}
		lgr.lgr.Log(ctx, level, msg, m)
	}
}

// LogIfError writes error info to the logger, if both the logger and the error are non-nil.
// It returns the error.
func (lgr *toggleLogger) LogIfError(ctx context.Context, err error) error {
	if err != nil {
		lgr.LogT(ctx, pgx.LogLevelError, "Error", nil, "error", err)
	}
	return err
}

// LogError writes error info to the logger, if the logger is not nil. It returns the error.
func (lgr *toggleLogger) LogError(ctx context.Context, err error) error {
	lgr.LogT(ctx, pgx.LogLevelError, "Error", nil, "error", err)
	return err
}

func (lgr *toggleLogger) TraceLogging(on bool) {
	if lgr.lgr == nil {
		return
	}

	// because pgx.Logger is an interface, it might be not nil yet hold a nil pointer
	value := reflect.ValueOf(lgr.lgr)
	if value.Kind() == reflect.Ptr && value.IsNil() {
		return
	}

	if on {
		atomic.StoreInt32(&lgr.enabled, 1)
	} else {
		atomic.StoreInt32(&lgr.enabled, 0)
	}
}

func (lgr *toggleLogger) loggingEnabled() bool {
	return atomic.LoadInt32(&lgr.enabled) != 0
}

func (lgr *toggleLogger) LogQueryWithError(ctx context.Context, err error, query string, args ...interface{}) {
	if err != nil {
		m := make(map[string]interface{})
		for i, v := range args {
			k := fmt.Sprintf("$%d", i+1)
			m[k] = derefArg(v)
		}
		m["error"] = err
		lgr.Log(ctx, pgx.LogLevelError, query, m)
	}
	// else no-op: pgx handles this
}

func (lgr *toggleLogger) LogQuery(_ context.Context, _ string, _ ...interface{}) {
	// no-op: pgx handles this
}

func (lgr *toggleLogger) SetOutput(w io.Writer) {
	if lgr.lgr != nil {
		if o, ok := lgr.lgr.(outable); ok {
			o.SetOutput(w)
		}
	}
}

func derefArg(arg interface{}) interface{} {
	value := reflect.ValueOf(arg)
	if value.Kind() == reflect.Ptr {
		return value.Elem()
	}
	return arg
}
