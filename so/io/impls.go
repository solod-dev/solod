// Copyright 2009 The Go Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package io

const maxint64 = int64(^uint64(0) >> 1)

// A DiscardWriter provides Write methods
// that succeed without doing anything.
type DiscardWriter struct{}

// Discard is a [Writer] on which all Write calls
// succeed without doing anything.
var Discard Writer = &DiscardWriter{}

func (*DiscardWriter) Write(p []byte) (int, error) {
	return len(p), nil
}

func (*DiscardWriter) WriteString(s string) (int, error) {
	return len(s), nil
}

// LimitReader returns a LimitedReader that reads from r
// but stops with EOF after n bytes.
func LimitReader(r Reader, n int64) LimitedReader { return LimitedReader{r, n} }

// A LimitedReader reads from R but limits the amount of
// data returned to just N bytes. Each call to Read
// updates N to reflect the new amount remaining.
// Read returns EOF when N <= 0 or when the underlying R returns EOF.
type LimitedReader struct {
	R Reader // underlying reader
	N int64  // max bytes remaining
}

func (l *LimitedReader) Read(p []byte) (int, error) {
	if l.N <= 0 {
		return 0, EOF
	}
	if int64(len(p)) > l.N {
		p = p[0:l.N]
	}
	n, err := l.R.Read(p)
	l.N -= int64(n)
	return n, err
}

// A NopCloser is a [ReadCloser] with a no-op Close method wrapping
// the provided [Reader] r.
type NopCloser struct {
	r Reader
}

// NewNopCloser returns a [NopCloser] wrapping r.
func NewNopCloser(r Reader) NopCloser {
	return NopCloser{r}
}

func (n *NopCloser) Read(p []byte) (int, error) {
	return n.r.Read(p)
}

func (*NopCloser) Close() error { return nil }

// NewSectionReader returns a [SectionReader] that reads from r
// starting at offset off and stops with EOF after n bytes.
func NewSectionReader(r ReaderAt, off int64, n int64) SectionReader {
	var remaining int64
	if off <= maxint64-n {
		remaining = n + off
	} else {
		// Overflow, with no way to return error.
		// Assume we can read up to an offset of 1<<63 - 1.
		remaining = maxint64
	}
	return SectionReader{r, off, off, remaining, n}
}

// SectionReader implements Read, Seek, and ReadAt on a section
// of an underlying [ReaderAt].
type SectionReader struct {
	r     ReaderAt // constant after creation
	base  int64    // constant after creation
	off   int64
	limit int64 // constant after creation
	n     int64 // constant after creation
}

func (s *SectionReader) Read(p []byte) (int, error) {
	if s.off >= s.limit {
		return 0, EOF
	}
	if max := s.limit - s.off; int64(len(p)) > max {
		p = p[0:max]
	}
	n, err := s.r.ReadAt(p, s.off)
	s.off += int64(n)
	return n, err
}

func (s *SectionReader) Seek(offset int64, whence int) (int64, error) {
	if whence == SeekStart {
		offset += s.base
	} else if whence == SeekCurrent {
		offset += s.off
	} else if whence == SeekEnd {
		offset += s.limit
	} else {
		return 0, ErrWhence
	}

	if offset < s.base {
		return 0, ErrOffset
	}
	s.off = offset
	return offset - s.base, nil
}

func (s *SectionReader) ReadAt(p []byte, off int64) (int, error) {
	if off < 0 || off >= s.Size() {
		return 0, EOF
	}
	off += s.base
	if max := s.limit - off; int64(len(p)) > max {
		p = p[0:max]
		n, err := s.r.ReadAt(p, off)
		if err == nil {
			err = EOF
		}
		return n, err
	}
	return s.r.ReadAt(p, off)
}

// Size returns the size of the section in bytes.
func (s *SectionReader) Size() int64 { return s.limit - s.base }

// ReaderAtOffset represents the underlying [ReaderAt] and offsets for a section.
type ReaderAtOffset struct {
	R   ReaderAt
	Off int64
	N   int64
}

// Outer returns the underlying [ReaderAt] and offsets for the section.
//
// The returned values are the same that were passed to [NewSectionReader]
// when the [SectionReader] was created.
func (s *SectionReader) Outer() ReaderAtOffset {
	return ReaderAtOffset{s.r, s.base, s.n}
}
