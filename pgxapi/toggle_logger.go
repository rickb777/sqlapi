package pgxapi

import (
	"github.com/jackc/pgx"
	"sync/atomic"
	"time"
)

type toggleLogger struct {
	lgr     pgx.Logger
	enabled int32
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
