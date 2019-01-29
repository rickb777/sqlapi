package sqlapi

import (
	"bytes"
	"fmt"
	"github.com/rickb777/sqlapi/schema"
	"log"
	"testing"
)

func TestLoggingOnOff(t *testing.T) {
	buf := &bytes.Buffer{}
	logger := log.New(buf, "X.", 0)

	db := NewDatabase(nil, schema.Sqlite, logger, nil)
	db.LogQuery("one")
	db.TraceLogging(false)
	db.LogQuery("two")
	db.TraceLogging(true)
	db.LogQuery("three")

	s := buf.String()
	if s != "X.one\nX.three\n" {
		t.Errorf("Got %q\n", s)
	}
}

func TestLoggingError(t *testing.T) {
	buf := &bytes.Buffer{}
	logger := log.New(buf, "X.", 0)

	db := NewDatabase(nil, schema.Sqlite, logger, nil)
	db.LogError(fmt.Errorf("one"))
	db.TraceLogging(false)
	db.LogError(fmt.Errorf("two"))
	db.TraceLogging(true)
	db.LogError(fmt.Errorf("three"))
	db.LogIfError(nil)
	db.LogIfError(fmt.Errorf("four"))

	s := buf.String()
	if s != "X.Error: one\nX.Error: three\nX.Error: four\n" {
		t.Errorf("Got %q\n", s)
	}
}

func TestSchemaSupport(t *testing.T) {
	for _, s := range schema.AllDialects {
		showTables(s)
	}
}
