// Copyright 2013 The Go Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package bufio

import (
	"solod.dev/so/bytes"
	"solod.dev/so/errors"
	"solod.dev/so/io"
	"solod.dev/so/mem"
	"solod.dev/so/unicode/utf8"
)

const maxInt = int(uint64(^uint(0)) >> 1)

// SplitResult holds the return values from a [SplitFunc].
type SplitResult struct {
	Advance  int
	Token    []byte
	HasToken bool
	Err      error
}

// SplitFunc is the signature of the split function used to tokenize the
// input. The arguments are an initial substring of the remaining unprocessed
// data and a flag, atEOF, that reports whether the [Reader] has no more data
// to give. The return value is a [SplitResult] containing the number of bytes
// to advance the input, the next token to return to the user (if any),
// and an error (if any).
//
// Scanning stops if the function returns an error, in which case some of
// the input may be discarded. If that error is [ErrFinalToken], scanning
// stops with no error. A token delivered with [ErrFinalToken] where
// HasToken is true will be the last token, and a result with HasToken
// false and [ErrFinalToken] immediately stops the scanning.
//
// Otherwise, the [Scanner] advances the input. If HasToken is true,
// the [Scanner] returns the token to the user. If HasToken is false, the
// Scanner reads more data and continues scanning; if there is no more
// data -- if atEOF was true -- the [Scanner] returns. If the data does not
// yet hold a complete token, for instance if it has no newline while
// scanning lines, a [SplitFunc] can return an empty [SplitResult] to signal the
// [Scanner] to read more data into the slice and try again with a
// longer slice starting at the same point in the input.
//
// The function is never called with an empty data slice unless atEOF
// is true. If atEOF is true, however, data may be non-empty and,
// as always, holds unprocessed text.
type SplitFunc func(data []byte, atEOF bool) SplitResult

// Scanner provides a convenient interface for reading data such as
// a file of newline-delimited lines of text. Successive calls to
// the [Scanner.Scan] method will step through the 'tokens' of a file, skipping
// the bytes between the tokens. The specification of a token is
// defined by a split function of type [SplitFunc]; the default split
// function breaks the input into lines with line termination stripped. [Scanner.Split]
// functions are defined in this package for scanning a file into
// lines, bytes, UTF-8-encoded runes, and space-delimited words. The
// client may instead provide a custom split function.
//
// Scanning stops unrecoverably at EOF, the first I/O error, or a token too
// large to fit in the [Scanner.Buffer]. When a scan stops, the reader may have
// advanced arbitrarily far past the last token. Programs that need more
// control over error handling or large tokens, or must run sequential scans
// on a reader, should use [bufio.Reader] instead.
type Scanner struct {
	a            mem.Allocator // Memory allocator for buffer growth.
	r            io.Reader     // The reader provided by the client.
	split        SplitFunc     // The function to split the tokens.
	maxTokenSize int           // Maximum size of a token; modified by tests.
	token        []byte        // Last token returned by split.
	buf          []byte        // Buffer used as argument to split.
	start        int           // First non-processed byte in buf.
	end          int           // End of data in buf.
	err          error         // Sticky error.
	empties      int           // Count of successive empty tokens.
	scanCalled   bool          // Scan has been called; buffer is in use.
	done         bool          // Scan has finished.
	ownsBuf      bool          // Whether the scanner allocated the buffer.
}

// Errors returned by Scanner.
var (
	ErrTooLong         = errors.New("bufio.Scanner: token too long")
	ErrNegativeAdvance = errors.New("bufio.Scanner: SplitFunc returns negative advance count")
	ErrAdvanceTooFar   = errors.New("bufio.Scanner: SplitFunc returns advance count beyond input")
	ErrBadReadCount    = errors.New("bufio.Scanner: Read returned impossible count")
)

const (
	// MaxScanTokenSize is the maximum size used to buffer a token
	// unless the user provides an explicit buffer with [Scanner.Buffer].
	// The actual maximum token size may be smaller as the buffer
	// may need to include, for instance, a newline.
	MaxScanTokenSize = 64 * 1024

	startBufSize = 4096 // Size of initial allocation for buffer.
)

// NewScanner returns a new [Scanner] to read from r.
// The split function defaults to [ScanLines].
//
// The caller is responsible for freeing the scanner resources
// with [Scanner.Free] when done using it.
func NewScanner(a mem.Allocator, r io.Reader) Scanner {
	return Scanner{
		a:            a,
		r:            r,
		split:        ScanLines,
		maxTokenSize: MaxScanTokenSize,
	}
}

// Free releases the internal scanner buffer.
// It is safe to call Free on a scanner that used a user-provided buffer
// via [Scanner.Buffer]; in that case Free is a no-op.
func (s *Scanner) Free() {
	if !s.ownsBuf {
		return
	}
	mem.FreeSlice(s.a, s.buf)
	s.buf = nil
	s.ownsBuf = false
}

// Err returns the first non-EOF error that was encountered by the [Scanner].
func (s *Scanner) Err() error {
	if s.err == io.EOF {
		return nil
	}
	return s.err
}

// Bytes returns the most recent token generated by a call to [Scanner.Scan].
// The underlying array may point to data that will be overwritten
// by a subsequent call to Scan. It does no allocation.
func (s *Scanner) Bytes() []byte {
	return s.token
}

// Text returns the most recent token generated by a call to [Scanner.Scan]
// as a string. The returned string is a zero-copy view into the buffer and
// is invalidated by the next call to [Scanner.Scan].
func (s *Scanner) Text() string {
	return string(s.token)
}

// ErrFinalToken is a special sentinel error value. It is intended to be
// returned by a Split function to indicate that the scanning should stop
// with no error. If the token being delivered with this error has HasToken
// true, the token is the last token.
//
// The value is useful to stop processing early or when it is necessary to
// deliver a final empty token. One could achieve the same behavior
// with a custom error value but providing one here is tidier.
var ErrFinalToken = errors.New("final token")

// Scan advances the [Scanner] to the next token, which will then be
// available through the [Scanner.Bytes] or [Scanner.Text] method. It returns false when
// there are no more tokens, either by reaching the end of the input or an error.
// After Scan returns false, the [Scanner.Err] method will return any error that
// occurred during scanning, except that if it was [io.EOF], [Scanner.Err]
// will return nil.
// Scan panics if the split function returns too many empty
// tokens without advancing the input. This is a common error mode for
// scanners.
func (s *Scanner) Scan() bool {
	if s.done {
		return false
	}
	s.scanCalled = true
	// Loop until we have a token.
	for {
		// See if we can get a token with what we already have.
		// If we've run out of data but have an error, give the split function
		// a chance to recover any remaining, possibly empty token.
		if s.end > s.start || s.err != nil {
			res := s.split(s.buf[s.start:s.end], s.err != nil)
			if res.Err != nil {
				if res.Err == ErrFinalToken {
					s.token = res.Token
					s.done = true
					return res.HasToken
				}
				s.setErr(res.Err)
				return false
			}
			if !s.advance(res.Advance) {
				return false
			}
			s.token = res.Token
			if res.HasToken {
				if s.err == nil || res.Advance > 0 {
					s.empties = 0
				} else {
					// Returning tokens not advancing input at EOF.
					s.empties++
					if s.empties > maxConsecutiveEmptyReads {
						panic("bufio.Scanner: too many empty tokens without progressing")
					}
				}
				return true
			}
		}
		// We cannot generate a token with what we are holding.
		// If we've already hit EOF or an I/O error, we are done.
		if s.err != nil {
			// Shut it down.
			s.start = 0
			s.end = 0
			return false
		}
		// Must read more data.
		// First, shift data to beginning of buffer if there's lots of empty space
		// or space is needed.
		if s.start > 0 && (s.end == len(s.buf) || s.start > len(s.buf)/2) {
			copy(s.buf, s.buf[s.start:s.end])
			s.end -= s.start
			s.start = 0
		}
		// Is the buffer full? If so, resize.
		if s.end == len(s.buf) {
			// Guarantee no overflow in the multiplication below.
			if len(s.buf) >= s.maxTokenSize || len(s.buf) > maxInt/2 {
				s.setErr(ErrTooLong)
				return false
			}
			newSize := len(s.buf) * 2
			if newSize == 0 {
				newSize = startBufSize
			}
			newSize = min(newSize, s.maxTokenSize)
			newBuf := mem.AllocSlice[byte](s.a, newSize, newSize)
			if s.end > s.start {
				// Only copy if s.buf is not nil to avoid UB in C.
				copy(newBuf, s.buf[s.start:s.end])
			}
			if s.ownsBuf {
				mem.FreeSlice(s.a, s.buf)
			}
			s.buf = newBuf
			s.ownsBuf = true
			s.end -= s.start
			s.start = 0
		}
		// Finally we can read some input. Make sure we don't get stuck with
		// a misbehaving Reader. Officially we don't need to do this, but let's
		// be extra careful: Scanner is for safe, simple jobs.
		for loop := 0; ; {
			n, err := s.r.Read(s.buf[s.end:len(s.buf)])
			if n < 0 || len(s.buf)-s.end < n {
				s.setErr(ErrBadReadCount)
				break
			}
			s.end += n
			if err != nil {
				s.setErr(err)
				break
			}
			if n > 0 {
				s.empties = 0
				break
			}
			loop++
			if loop > maxConsecutiveEmptyReads {
				s.setErr(io.ErrNoProgress)
				break
			}
		}
	}
}

// advance consumes n bytes of the buffer. It reports whether the advance was legal.
func (s *Scanner) advance(n int) bool {
	if n < 0 {
		s.setErr(ErrNegativeAdvance)
		return false
	}
	if n > s.end-s.start {
		s.setErr(ErrAdvanceTooFar)
		return false
	}
	s.start += n
	return true
}

// setErr records the first error encountered.
func (s *Scanner) setErr(err error) {
	if s.err == nil || s.err == io.EOF {
		s.err = err
	}
}

// Buffer controls memory allocation by the Scanner.
// It sets the initial buffer to use when scanning
// and the maximum size of buffer that may be allocated during scanning.
// The contents of the buffer are ignored.
//
// The maximum token size must be less than the larger of max and cap(buf).
// If max <= cap(buf), [Scanner.Scan] will use this buffer only and do no allocation.
//
// By default, [Scanner.Scan] uses an internal buffer and sets the
// maximum token size to [MaxScanTokenSize].
//
// Buffer panics if it is called after scanning has started.
func (s *Scanner) Buffer(buf []byte, max int) {
	if s.scanCalled {
		panic("bufio.Scanner: Buffer called after Scan")
	}
	if s.ownsBuf {
		mem.FreeSlice(s.a, s.buf)
		s.ownsBuf = false
	}
	s.buf = buf[0:cap(buf)]
	s.maxTokenSize = max
}

// Split sets the split function for the [Scanner].
// The default split function is [ScanLines].
//
// Split panics if it is called after scanning has started.
func (s *Scanner) Split(split SplitFunc) {
	if s.scanCalled {
		panic("bufio.Scanner: Split called after Scan")
	}
	s.split = split
}

// Split functions

// ScanBytes is a split function for a [Scanner] that returns each byte as a token.
func ScanBytes(data []byte, atEOF bool) SplitResult {
	if atEOF && len(data) == 0 {
		return SplitResult{}
	}
	return SplitResult{Advance: 1, Token: data[0:1], HasToken: true}
}

// errorRune is the UTF-8 encoding of U+FFFD (replacement character).
var errorRune = [3]byte{0xef, 0xbf, 0xbd}

// ScanRunes is a split function for a [Scanner] that returns each
// UTF-8-encoded rune as a token. The sequence of runes returned is
// equivalent to that from a range loop over the input as a string, which
// means that erroneous UTF-8 encodings translate to U+FFFD = "\xef\xbf\xbd".
// Because of the Scan interface, this makes it impossible for the client to
// distinguish correctly encoded replacement runes from encoding errors.
func ScanRunes(data []byte, atEOF bool) SplitResult {
	if atEOF && len(data) == 0 {
		return SplitResult{}
	}

	// Fast path 1: ASCII.
	if data[0] < utf8.RuneSelf {
		return SplitResult{Advance: 1, Token: data[0:1], HasToken: true}
	}

	// Fast path 2: Correct UTF-8 decode without error.
	_, width := utf8.DecodeRune(data)
	if width > 1 {
		// It's a valid encoding. Width cannot be one for a correctly encoded
		// non-ASCII rune.
		return SplitResult{Advance: width, Token: data[0:width], HasToken: true}
	}

	// We know it's an error: we have width==1 and implicitly r==utf8.RuneError.
	// Is the error because there wasn't a full rune to be decoded?
	// FullRune distinguishes correctly between erroneous and incomplete encodings.
	if !atEOF && !utf8.FullRune(data) {
		// Incomplete; get more bytes.
		return SplitResult{}
	}

	// We have a real UTF-8 encoding error. Return a properly encoded error rune
	// but advance only one byte. This matches the behavior of a range loop over
	// an incorrectly encoded string.
	return SplitResult{Advance: 1, Token: errorRune[:], HasToken: true}
}

// dropCR drops a terminal \r from the data.
func dropCR(data []byte) []byte {
	if len(data) > 0 && data[len(data)-1] == '\r' {
		return data[0 : len(data)-1]
	}
	return data
}

// ScanLines is a split function for a [Scanner] that returns each line of
// text, stripped of any trailing end-of-line marker. The returned line may
// be empty. The end-of-line marker is one optional carriage return followed
// by one mandatory newline. In regular expression notation, it is `\r?\n`.
// The last non-empty line of input will be returned even if it has no
// newline.
func ScanLines(data []byte, atEOF bool) SplitResult {
	if atEOF && len(data) == 0 {
		return SplitResult{}
	}
	if i := bytes.IndexByte(data, '\n'); i >= 0 {
		// We have a full newline-terminated line.
		return SplitResult{Advance: i + 1, Token: dropCR(data[0:i]), HasToken: true}
	}
	// If we're at EOF, we have a final, non-terminated line. Return it.
	if atEOF {
		return SplitResult{Advance: len(data), Token: dropCR(data), HasToken: true}
	}
	// Request more data.
	return SplitResult{}
}

// isSpace reports whether the character is a Unicode white space character.
// We avoid dependency on the unicode package, but check validity of the implementation
// in the tests.
func isSpace(r rune) bool {
	if r <= 0x00FF {
		// Obvious ASCII ones: \t through \r plus space. Plus two Latin-1 oddballs.
		switch r {
		case ' ', '\t', '\n', '\v', '\f', '\r':
			return true
		case 0x0085, 0x00A0:
			return true
		}
		return false
	}
	// High-valued ones.
	if 0x2000 <= r && r <= 0x200a {
		return true
	}
	switch r {
	case 0x1680, 0x2028, 0x2029, 0x202f, 0x205f, 0x3000:
		return true
	}
	return false
}

// ScanWords is a split function for a [Scanner] that returns each
// space-separated word of text, with surrounding spaces deleted. It will
// never return an empty string. The definition of space is set by
// unicode.IsSpace.
func ScanWords(data []byte, atEOF bool) SplitResult {
	// Skip leading spaces.
	start := 0
	for width := 0; start < len(data); start += width {
		var r rune
		r, width = utf8.DecodeRune(data[start:])
		if !isSpace(r) {
			break
		}
	}
	// Scan until space, marking end of word.
	width := 0
	for i := start; i < len(data); i += width {
		var r rune
		r, width = utf8.DecodeRune(data[i:])
		if isSpace(r) {
			return SplitResult{Advance: i + width, Token: data[start:i], HasToken: true}
		}
	}
	// If we're at EOF, we have a final, non-empty, non-terminated word. Return it.
	if atEOF && len(data) > start {
		return SplitResult{Advance: len(data), Token: data[start:], HasToken: true}
	}
	// Request more data.
	return SplitResult{Advance: start}
}
