package sqlapi

import (
	"bytes"
	"fmt"
	. "github.com/onsi/gomega"
	"github.com/rickb777/sqlapi/dialect"
	"log"
	"testing"
)

func TestLoggingOnOff(t *testing.T) {
	g := NewGomegaWithT(t)

	buf := &bytes.Buffer{}
	logger := log.New(buf, "X.", 0)

	db := NewDatabase(nil, dialect.Sqlite, logger, nil)
	db.LogQuery("one")
	db.TraceLogging(false)
	db.LogQuery("two")
	db.TraceLogging(true)
	db.LogQuery("three")

	s := buf.String()
	g.Expect(s).To(Equal("X.one\nX.three\n"))
}

func TestLoggingError(t *testing.T) {
	g := NewGomegaWithT(t)

	buf := &bytes.Buffer{}
	logger := log.New(buf, "X.", 0)

	db := NewDatabase(nil, dialect.Sqlite, logger, nil)
	db.LogError(fmt.Errorf("one"))
	db.TraceLogging(false)
	db.LogError(fmt.Errorf("two"))
	db.TraceLogging(true)
	db.LogError(fmt.Errorf("three"))
	db.LogIfError(nil)
	db.LogIfError(fmt.Errorf("four"))

	s := buf.String()
	g.Expect(s).To(Equal("X.Error: one\nX.Error: three\nX.Error: four\n"))
}
