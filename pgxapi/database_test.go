package pgxapi

import (
	"bytes"
	"errors"
	"fmt"
	. "github.com/onsi/gomega"
	"github.com/mercury-holidays/sqlapi/dialect"
	"github.com/mercury-holidays/sqlapi/pgxapi/logadapter"
	"log"
	"testing"
)

func TestLoggingOnOff(t *testing.T) {
	g := NewGomegaWithT(t)

	buf := &bytes.Buffer{}
	logger := logadapter.NewLogger(log.New(buf, "X.", 0))
	sh := &shim{lgr: &toggleLogger{lgr: logger, enabled: 1}}

	db := NewDatabase(sh, dialect.Sqlite, nil)
	db.Logger().LogError(errors.New("one"))
	db.Logger().LogError(errors.New("two"))

	s := buf.String()
	g.Expect(s).To(Equal("X.Error [error:one]\nX.Error [error:two]\n"))
}

func TestLoggingError(t *testing.T) {
	g := NewGomegaWithT(t)

	buf := &bytes.Buffer{}
	logger := logadapter.NewLogger(log.New(buf, "X.", 0))
	sh := &shim{lgr: &toggleLogger{lgr: logger, enabled: 1}}

	db := NewDatabase(sh, dialect.Sqlite, nil)
	db.Logger().LogIfError(nil)
	db.Logger().LogIfError(fmt.Errorf("four"))
	db.Logger().LogIfError(nil)

	s := buf.String()
	g.Expect(s).To(Equal("X.Error [error:four]\n"))
}
