package dialect

import (
	"io"
)

// StringWriter is an interface that wraps the WriteString method.
// Note that bytes.Buffer happens to implement this interface.
type StringWriter interface {
	io.Writer
	WriteString(s string) (n int, err error)
}

type swAdapter struct {
	w io.Writer
}

// Write writes bytes to its writer.
func (w swAdapter) Write(b []byte) (n int, err error) {
	return w.w.Write(b)
}

// WriteString writes a string to its writer.
func (w swAdapter) WriteString(s string) (n int, err error) {
	return w.w.Write([]byte(s))
}

// Adapt wraps an io.Writer as a StringWriter.
func Adapt(w io.Writer) StringWriter {
	if sw, ok := w.(StringWriter); ok {
		return sw
	}
	return swAdapter{w}
}
