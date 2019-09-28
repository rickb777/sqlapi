package pgxapi

import (
	"sync/atomic"
	"time"

	"github.com/jackc/pgx"
)

type toggleLogger struct {
	lgr     pgx.Logger
	enabled int32
}

func NewLogger(lgr pgx.Logger) Logger {
	if lgr == nil {
		return &toggleLogger{}
	}
	return &toggleLogger{lgr: lgr, enabled: 1}
}

func (lgr *toggleLogger) Log(level pgx.LogLevel, msg string, data map[string]interface{}) {
	if lgr.loggingEnabled() {
		lgr.lgr.Log(level, msg, data)
	}
}

// Log emits a log event, supporting an elapsed-time calculation and providing an easier
// way to supply data parameters as name,value pairs.
func (lgr *toggleLogger) LogT(level pgx.LogLevel, msg string, startTime *time.Time, data ...interface{}) {
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
		lgr.lgr.Log(level, msg, m)
	}
}

func (lgr *toggleLogger) LogQuery(query string, args ...interface{}) {
	// no-op: pgx handles this
}

// LogIfError writes error info to the logger, if both the logger and the error are non-nil.
// It returns the error.
func (lgr *toggleLogger) LogIfError(err error) error {
	if err != nil {
		lgr.LogT(pgx.LogLevelError, "Error", nil, "error", err)
	}
	return err
}

// LogError writes error info to the logger, if the logger is not nil. It returns the error.
func (lgr *toggleLogger) LogError(err error) error {
	lgr.LogT(pgx.LogLevelError, "Error", nil, "error", err)
	return err
}

func (lgr *toggleLogger) TraceLogging(on bool) {
	if on {
		atomic.StoreInt32(&lgr.enabled, 1)
	} else {
		atomic.StoreInt32(&lgr.enabled, 0)
	}
}

func (lgr *toggleLogger) loggingEnabled() bool {
	return atomic.LoadInt32(&lgr.enabled) != 0
}

//// LogQuery writes query info to the logger, if it is not nil.
//func (database *database) LogQuery(query string, args ...interface{}) {
//	if database.loggingEnabled() {
//		query = strings.TrimSpace(query)
//		if len(args) > 0 {
//			ss := make([]interface{}, len(args))
//			for i, v := range args {
//				ss[i] = derefArg(v)
//			}
//			database.logger.Printf("%s %v\n", query, ss)
//		} else {
//			database.logger.Println(query)
//		}
//	}
//}

//func derefArg(arg interface{}) interface{} {
//	switch v := arg.(type) {
//	case *int:
//		return *v
//	case *int8:
//		return *v
//	case *int16:
//		return *v
//	case *int32:
//		return *v
//	case *int64:
//		return *v
//	case *uint:
//		return *v
//	case *uint8:
//		return *v
//	case *uint16:
//		return *v
//	case *uint32:
//		return *v
//	case *uint64:
//		return *v
//	case *float32:
//		return *v
//	case *float64:
//		return *v
//	case *bool:
//		return *v
//	case *string:
//		return *v
//	}
//	return arg
//}
