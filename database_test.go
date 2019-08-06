package sqlapi

import (
	"bytes"
	"fmt"
	. "github.com/onsi/gomega"
	"github.com/mercury-holidays/sqlapi/dialect"
	"log"
	"testing"
)

func TestLoggingOnOff(t *testing.T) {
	g := NewGomegaWithT(t)

	buf := &bytes.Buffer{}
	logger := log.New(buf, "X.", 0)

	db := NewDatabase(nil, dialect.Sqlite, logger, nil)
	db.Logger().LogQuery("one")
	db.Logger().TraceLogging(false)
	db.Logger().LogQuery("two")
	db.Logger().TraceLogging(true)
	db.Logger().LogQuery("three")

	s := buf.String()
	g.Expect(s).To(Equal("X.one\nX.three\n"))
}

func TestLoggingError(t *testing.T) {
	g := NewGomegaWithT(t)

	buf := &bytes.Buffer{}
	logger := log.New(buf, "X.", 0)

	db := NewDatabase(nil, dialect.Sqlite, logger, nil)
	db.Logger().LogError(fmt.Errorf("one"))
	db.Logger().TraceLogging(false)
	db.Logger().LogError(fmt.Errorf("two"))
	db.Logger().TraceLogging(true)
	db.Logger().LogError(fmt.Errorf("three"))
	db.Logger().LogIfError(nil)
	db.Logger().LogIfError(fmt.Errorf("four"))

	s := buf.String()
	g.Expect(s).To(Equal("X.Error: one\nX.Error: three\nX.Error: four\n"))
}
