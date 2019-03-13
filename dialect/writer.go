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

//-------------------------------------------------------------------------------------------------

//func NewBuffer(d Dialect) *Buffer {
//	return &Buffer{
//		buf:      &bytes.Buffer{},
//		quoter:   d.Quoter(),
//		numbered: d.HasNumberedPlaceholders(),
//	}
//}
//
//type Buffer struct {
//	buf      *bytes.Buffer
//	quoter   Quoter
//	numbered bool
//}
//
//func (b *Buffer) Len() int {
//	return b.buf.Len()
//}
//
//func (b *Buffer) Write(bb []byte) (n int, err error) {
//	return b.buf.Write(bb)
//}
//
//func (b *Buffer) WriteByte(bb byte) error {
//	return b.buf.WriteByte(bb)
//}
//
//func (b *Buffer) WriteRune(r rune) (n int, err error) {
//	return b.buf.WriteRune(r)
//}
//
//func (b *Buffer) WriteString(s string) (n int, err error) {
//	return b.buf.WriteString(s)
//}
//
//// Þ appends a string to the buffer.
//func (b *Buffer) Þ(s string) *Buffer {
//	b.buf.WriteString(s)
//	return b
//}
//
//// Append appends a string to the buffer.
//func (b *Buffer) Append(s string) *Buffer {
//	b.buf.WriteString(s)
//	return b
//}
//
//// Quote appends a quoted identifier to the buffer.
//func (b *Buffer) Quote(id string) *Buffer {
//	b.quoter.QuoteW(b.buf, id)
//	return b
//}
//
//func (b *Buffer) String() string {
//	return b.buf.String()
//}
//
//// ReplacePlaceholders converts a string containing '?' placeholders to
//// the form used by PostgreSQL.
//func (b *Buffer) ReplacePlaceholders() *Buffer {
//	if !b.numbered {
//		return b
//	}
//
//	b2 := &Buffer{
//		buf:      &bytes.Buffer{},
//		quoter:   b.quoter,
//		numbered: b.numbered,
//	}
//	idx := 1
//	sql := b.buf.String()
//	for _, r := range sql {
//		if r == '?' {
//			b2.buf.WriteByte('$')
//			b2.buf.WriteString(strconv.Itoa(idx))
//			idx++
//		} else {
//			b2.buf.WriteRune(r)
//		}
//	}
//	return b2
//}
