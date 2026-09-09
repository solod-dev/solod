package bufio

import (
	"solod.dev/so/io"
	"solod.dev/so/mem"
	"solod.dev/so/unicode/utf8"
)

// Writer implements buffering for an [io.Writer] object.
// If an error occurs writing to a [Writer], no more data will be
// accepted and all subsequent writes, and [Writer.Flush], will return the error.
// After all data has been written, the client should call the
// [Writer.Flush] method to guarantee all data has been forwarded to
// the underlying [io.Writer].
//
// The caller is responsible for freeing the writer's resources
// with [Writer.Free] when done using it.
type Writer struct {
	a   mem.Allocator
	err error
	buf []byte
	n   int
	wr  io.Writer
}

// NewWriterSize returns a new [Writer] whose buffer has at least the specified
// size. If the argument io.Writer is already a [Writer] with large enough
// size, it returns the underlying [Writer].
func NewWriterSize(a mem.Allocator, w io.Writer, size int) Writer {
	// Is it already a Writer?
	b, ok := w.(*Writer)
	if ok && len(b.buf) >= size {
		return *b
	}
	if size <= 0 {
		size = DefaultBufSize
	}
	wr := Writer{a: a}
	wr.buf = mem.AllocSlice[byte](a, size, size)
	wr.wr = w
	return wr
}

// NewWriter returns a new [Writer] whose buffer has the default size.
// If the argument io.Writer is already a [Writer] with large enough buffer size,
// it returns the underlying [Writer].
func NewWriter(a mem.Allocator, w io.Writer) Writer {
	return NewWriterSize(a, w, DefaultBufSize)
}

// Size returns the size of the underlying buffer in bytes.
func (b *Writer) Size() int { return len(b.buf) }

// Reset discards any unflushed buffered data, clears any error, and
// resets b to write its output to w.
// Calling Reset on the zero value of [Writer] initializes the internal buffer
// to the default size.
// Calling w.Reset(w) (that is, resetting a [Writer] to itself) does nothing.
func (b *Writer) Reset(w io.Writer) {
	// Avoid no-op reset to self.
	if wb, ok := w.(*Writer); ok && wb == b {
		return
	}
	if b.buf == nil {
		b.buf = mem.AllocSlice[byte](b.a, DefaultBufSize, DefaultBufSize)
	}
	b.err = nil
	b.n = 0
	b.wr = w
}

// Free releases the internal buffer memory.
// The writer must not be used after calling Free.
func (b *Writer) Free() {
	mem.FreeSlice(b.a, b.buf)
	b.buf = nil
}

// Flush writes any buffered data to the underlying [io.Writer].
func (b *Writer) Flush() error {
	if b.err != nil {
		return b.err
	}
	if b.n == 0 {
		return nil
	}
	n, err := b.wr.Write(b.buf[0:b.n])
	if n < b.n && err == nil {
		err = io.ErrShortWrite
	}
	if err != nil {
		if n > 0 && n < b.n {
			copy(b.buf[0:b.n-n], b.buf[n:b.n])
		}
		b.n -= n
		b.err = err
		return err
	}
	b.n = 0
	return nil
}

// Available returns how many bytes are unused in the buffer.
func (b *Writer) Available() int { return len(b.buf) - b.n }

// AvailableBuffer returns an empty buffer with b.Available() capacity.
// This buffer is intended to be appended to and
// passed to an immediately succeeding [Writer.Write] call.
// The buffer is only valid until the next write operation on b.
func (b *Writer) AvailableBuffer() []byte {
	return b.buf[b.n:][:0]
}

// Buffered returns the number of bytes that have been written into the current buffer.
func (b *Writer) Buffered() int { return b.n }

// Write writes the contents of p into the buffer.
// It returns the number of bytes written.
// If nn < len(p), it also returns an error explaining
// why the write is short.
func (b *Writer) Write(p []byte) (int, error) {
	nn := 0
	for len(p) > b.Available() && b.err == nil {
		var n int
		if b.Buffered() == 0 {
			// Large write, empty buffer.
			// Write directly from p to avoid copy.
			n, b.err = b.wr.Write(p)
		} else {
			n = copy(b.buf[b.n:], p)
			b.n += n
			b.Flush()
		}
		nn += n
		p = p[n:]
	}
	if b.err != nil {
		return nn, b.err
	}
	n := copy(b.buf[b.n:], p)
	b.n += n
	nn += n
	return nn, nil
}

// WriteByte writes a single byte.
func (b *Writer) WriteByte(c byte) error {
	if b.err != nil {
		return b.err
	}
	if b.Available() <= 0 && b.Flush() != nil {
		return b.err
	}
	b.buf[b.n] = c
	b.n++
	return nil
}

// WriteRune writes a single Unicode code point, returning
// the number of bytes written and any error.
func (b *Writer) WriteRune(r rune) (int, error) {
	// Compare as uint32 to correctly handle negative runes.
	if uint32(r) < utf8.RuneSelf {
		err := b.WriteByte(byte(r))
		if err != nil {
			return 0, err
		}
		return 1, nil
	}
	if b.err != nil {
		return 0, b.err
	}
	n := b.Available()
	if n < utf8.UTFMax {
		if b.Flush(); b.err != nil {
			return 0, b.err
		}
		n = b.Available()
		if n < utf8.UTFMax {
			// Can only happen if buffer is silly small.
			return b.WriteString(string(r))
		}
	}
	size := utf8.EncodeRune(b.buf[b.n:], r)
	b.n += size
	return size, nil
}

// WriteString writes a string.
// It returns the number of bytes written.
// If the count is less than len(s), it also returns an error explaining
// why the write is short.
func (b *Writer) WriteString(s string) (int, error) {
	nn := 0
	for len(s) > b.Available() && b.err == nil {
		n := copy(b.buf[b.n:], s)
		b.n += n
		b.Flush()
		nn += n
		s = s[n:]
	}
	if b.err != nil {
		return nn, b.err
	}
	n := copy(b.buf[b.n:], s)
	b.n += n
	nn += n
	return nn, nil
}

// ReadFrom implements [io.ReaderFrom].
func (b *Writer) ReadFrom(r io.Reader) (int64, error) {
	if b.err != nil {
		return 0, b.err
	}
	n := int64(0)
	var m int
	var err error
	for {
		if b.Available() == 0 {
			if err1 := b.Flush(); err1 != nil {
				return n, err1
			}
		}
		nr := 0
		for nr < maxConsecutiveEmptyReads {
			m, err = r.Read(b.buf[b.n:])
			if m != 0 || err != nil {
				break
			}
			nr++
		}
		if nr == maxConsecutiveEmptyReads {
			return n, io.ErrNoProgress
		}
		b.n += m
		n += int64(m)
		if err != nil {
			break
		}
	}
	if err == io.EOF {
		// If we filled the buffer exactly, flush preemptively.
		if b.Available() == 0 {
			err = b.Flush()
		} else {
			err = nil
		}
	}
	return n, err
}

// ReadWriter stores pointers to a [Reader] and a [Writer].
// It implements [io.ReadWriter].
type ReadWriter struct {
	r *Reader
	w *Writer
}

// NewReadWriter allocates a new [ReadWriter] that dispatches to r and w.
func NewReadWriter(r *Reader, w *Writer) ReadWriter {
	return ReadWriter{r: r, w: w}
}

func (rw *ReadWriter) Read(p []byte) (int, error) {
	return rw.r.Read(p)
}

func (rw *ReadWriter) Write(p []byte) (int, error) {
	return rw.w.Write(p)
}
