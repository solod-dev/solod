// Copyright 2009 The Go Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

// Package hex implements hexadecimal encoding and decoding.
// Based on the [encoding/hex] package.
//
// [encoding/hex]: https://github.com/golang/go/blob/go1.26.2/src/encoding/hex
package hex

import (
	"solod.dev/so/errors"
	"solod.dev/so/io"
	"solod.dev/so/mem"
	"solod.dev/so/strings"
)

const (
	hextable        = "0123456789abcdef"
	reverseHexTable = "" +
		"\xff\xff\xff\xff\xff\xff\xff\xff\xff\xff\xff\xff\xff\xff\xff\xff" +
		"\xff\xff\xff\xff\xff\xff\xff\xff\xff\xff\xff\xff\xff\xff\xff\xff" +
		"\xff\xff\xff\xff\xff\xff\xff\xff\xff\xff\xff\xff\xff\xff\xff\xff" +
		"\x00\x01\x02\x03\x04\x05\x06\x07\x08\x09\xff\xff\xff\xff\xff\xff" +
		"\xff\x0a\x0b\x0c\x0d\x0e\x0f\xff\xff\xff\xff\xff\xff\xff\xff\xff" +
		"\xff\xff\xff\xff\xff\xff\xff\xff\xff\xff\xff\xff\xff\xff\xff\xff" +
		"\xff\x0a\x0b\x0c\x0d\x0e\x0f\xff\xff\xff\xff\xff\xff\xff\xff\xff" +
		"\xff\xff\xff\xff\xff\xff\xff\xff\xff\xff\xff\xff\xff\xff\xff\xff" +
		"\xff\xff\xff\xff\xff\xff\xff\xff\xff\xff\xff\xff\xff\xff\xff\xff" +
		"\xff\xff\xff\xff\xff\xff\xff\xff\xff\xff\xff\xff\xff\xff\xff\xff" +
		"\xff\xff\xff\xff\xff\xff\xff\xff\xff\xff\xff\xff\xff\xff\xff\xff" +
		"\xff\xff\xff\xff\xff\xff\xff\xff\xff\xff\xff\xff\xff\xff\xff\xff" +
		"\xff\xff\xff\xff\xff\xff\xff\xff\xff\xff\xff\xff\xff\xff\xff\xff" +
		"\xff\xff\xff\xff\xff\xff\xff\xff\xff\xff\xff\xff\xff\xff\xff\xff" +
		"\xff\xff\xff\xff\xff\xff\xff\xff\xff\xff\xff\xff\xff\xff\xff\xff" +
		"\xff\xff\xff\xff\xff\xff\xff\xff\xff\xff\xff\xff\xff\xff\xff\xff"
)

// ErrLength reports an attempt to decode an odd-length input
// using [Decode] or [DecodeString].
// The stream-based Decoder returns [io.ErrUnexpectedEOF] instead of ErrLength.
var ErrLength = errors.New("encoding/hex: odd length hex string")

// ErrInvalidByte reports an invalid byte in the input.
var ErrInvalidByte = errors.New("encoding/hex: invalid byte")

// ErrDumperClosed reports an attempt to write to a closed dumper.
var ErrDumperClosed = errors.New("encoding/hex: dumper closed")

// EncodedLen returns the length of an encoding of n source bytes.
// Specifically, it returns n * 2.
func EncodedLen(n int) int { return n * 2 }

// Encode encodes src into [EncodedLen](len(src))
// bytes of dst. As a convenience, it returns the number
// of bytes written to dst, but this value is always [EncodedLen](len(src)).
// Encode implements hexadecimal encoding.
func Encode(dst, src []byte) int {
	j := 0
	if src == nil { // for gcc analyzer
		return 0
	}
	if dst == nil && len(src) != 0 { // for gcc analyzer
		panic("encoding/hex: nil dst")
	}
	for _, v := range src {
		dst[j] = hextable[v>>4]
		dst[j+1] = hextable[v&0x0f]
		j += 2
	}
	return len(src) * 2
}

// AppendEncode appends the hexadecimally encoded src to dst
// and returns the extended buffer.
func AppendEncode(a mem.Allocator, dst, src []byte) []byte {
	n := EncodedLen(len(src))
	dst = mem.ReallocSlice(a, dst, len(dst), len(dst)+n)
	Encode(dst[len(dst):][:n], src)
	return dst[:len(dst)+n]
}

// DecodedLen returns the length of a decoding of x source bytes.
// Specifically, it returns x / 2.
func DecodedLen(x int) int { return x / 2 }

// Decode decodes src into [DecodedLen](len(src)) bytes,
// returning the actual number of bytes written to dst.
//
// Decode expects that src contains only hexadecimal
// characters and that src has even length.
// If the input is malformed, Decode returns the number
// of bytes decoded before the error.
func Decode(dst, src []byte) (int, error) {
	i, j := 0, 1
	if src == nil { // for gcc analyzer
		return 0, nil
	}
	if dst == nil && len(src) > 1 { // for gcc analyzer
		panic("encoding/hex: nil dst")
	}
	for ; j < len(src); j += 2 {
		p := src[j-1]
		q := src[j]

		a := reverseHexTable[p]
		b := reverseHexTable[q]
		if a > 0x0f {
			return i, ErrInvalidByte
		}
		if b > 0x0f {
			return i, ErrInvalidByte
		}
		dst[i] = (a << 4) | b
		i++
	}
	if len(src)%2 == 1 {
		// Check for invalid char before reporting bad length,
		// since the invalid char (if present) is an earlier problem.
		if reverseHexTable[src[j-1]] > 0x0f {
			return i, ErrInvalidByte
		}
		return i, ErrLength
	}
	return i, nil
}

// AppendDecode appends the hexadecimally decoded src to dst
// and returns the extended buffer.
// If the input is malformed, it returns the partially decoded src and an error.
func AppendDecode(a mem.Allocator, dst, src []byte) ([]byte, error) {
	n := DecodedLen(len(src))
	dst = mem.ReallocSlice(a, dst, len(dst), len(dst)+n)
	n, err := Decode(dst[len(dst):][:n], src)
	return dst[:len(dst)+n], err
}

// EncodeToString returns the hexadecimal encoding of src.
// The returned string is allocated; the caller owns it.
func EncodeToString(a mem.Allocator, src []byte) string {
	n := EncodedLen(len(src))
	dst := mem.AllocSlice[byte](a, n, n)
	Encode(dst, src)
	return string(dst)
}

// DecodeString returns the bytes represented by the hexadecimal string s.
//
// DecodeString expects that src contains only hexadecimal
// characters and that src has even length.
// If the input is malformed, DecodeString returns
// the bytes decoded before the error.
//
// The returned slice is allocated; the caller owns it.
func DecodeString(a mem.Allocator, s string) ([]byte, error) {
	n := DecodedLen(len(s))
	dst := mem.AllocSlice[byte](a, n, n)
	n, err := Decode(dst, []byte(s))
	return dst[:n], err
}

// Dump returns a string that contains a hex dump of the given data. The format
// of the hex dump matches the output of `hexdump -C` on the command line.
// The returned string is allocated; the caller owns it.
func Dump(a mem.Allocator, data []byte) string {
	if len(data) == 0 {
		return ""
	}

	buf := strings.NewBuilder(a)
	// Dumper will write 79 bytes per complete 16 byte chunk, and at least
	// 64 bytes for whatever remains. Round the allocation up, since only a
	// maximum of 15 bytes will be wasted.
	buf.Grow((1 + ((len(data) - 1) / 16)) * 79)

	dumper := NewDumper(&buf)
	dumper.Write(data)
	dumper.Close()
	return buf.String()
}

// bufferSize is the number of hexadecimal characters to buffer in encoder and decoder.
const bufferSize = 1024

// Encoder implements an [io.Writer] that writes hexadecimal characters
// to an underlying [io.Writer].
type Encoder struct {
	w   io.Writer
	err error
	out [bufferSize]byte // output buffer
}

// NewEncoder returns a new [Encoder] that writes to w.
func NewEncoder(w io.Writer) Encoder {
	return Encoder{w: w}
}

func (e *Encoder) Write(p []byte) (int, error) {
	var n int
	for len(p) > 0 && e.err == nil {
		chunkSize := bufferSize / 2
		if len(p) < chunkSize {
			chunkSize = len(p)
		}

		var written int
		encoded := Encode(e.out[:], p[:chunkSize])
		written, e.err = e.w.Write(e.out[:encoded])
		n += written / 2
		p = p[chunkSize:]
	}
	return n, e.err
}

// Decoder implements an [io.Reader] that reads hexadecimal characters
// from an underlying [io.Reader] and decodes them.
type Decoder struct {
	r   io.Reader
	err error
	in  []byte           // input buffer (encoded form)
	arr [bufferSize]byte // backing array for in
}

// NewDecoder returns a new [Decoder] that reads from r.
// NewDecoder expects that r contain only an even number of hexadecimal characters.
func NewDecoder(r io.Reader) Decoder {
	return Decoder{r: r}
}

func (d *Decoder) Read(p []byte) (int, error) {
	// Fill internal buffer with sufficient bytes to decode
	if len(d.in) < 2 && d.err == nil {
		var numCopy, numRead int
		numCopy = copy(d.arr[:], d.in) // Copies either 0 or 1 bytes
		numRead, d.err = d.r.Read(d.arr[numCopy:])
		d.in = d.arr[:numCopy+numRead]
		if d.err == io.EOF && len(d.in)%2 != 0 {
			if a := reverseHexTable[d.in[len(d.in)-1]]; a > 0x0f {
				d.err = ErrInvalidByte
			} else {
				d.err = io.ErrUnexpectedEOF
			}
		}
	}

	// Decode internal buffer into output buffer
	if numAvail := len(d.in) / 2; len(p) > numAvail {
		p = p[:numAvail]
	}
	numDec, err := Decode(p, d.in[:len(p)*2])
	d.in = d.in[2*numDec:]
	if err != nil {
		d.in, d.err = nil, err // Decode error; discard input remainder
	}

	if len(d.in) < 2 {
		return numDec, d.err // Only expose errors when buffer fully consumed
	}
	return numDec, nil
}

type Dumper struct {
	w          io.Writer
	rightChars [18]byte
	buf        [14]byte
	used       int  // number of bytes in the current line
	n          uint // number of bytes, total
	closed     bool
}

// NewDumper returns a [io.WriteCloser] that writes a hex dump of all written data to
// w. The format of the dump matches the output of `hexdump -C` on the command
// line.
func NewDumper(w io.Writer) Dumper {
	return Dumper{w: w}
}

func (h *Dumper) Write(data []byte) (int, error) {
	if h.closed {
		return 0, ErrDumperClosed
	}

	// Output lines look like:
	// 00000010  2e 2f 30 31 32 33 34 35  36 37 38 39 3a 3b 3c 3d  |./0123456789:;<=|
	// ^ offset                          ^ extra space              ^ ASCII of line.
	var n int
	for i := range data {
		if h.used == 0 {
			// At the beginning of a line we print the current
			// offset in hex.
			h.buf[0] = byte(h.n >> 24)
			h.buf[1] = byte(h.n >> 16)
			h.buf[2] = byte(h.n >> 8)
			h.buf[3] = byte(h.n)
			Encode(h.buf[4:], h.buf[:4])
			h.buf[12] = ' '
			h.buf[13] = ' '
			_, err := h.w.Write(h.buf[4:])
			if err != nil {
				return n, err
			}
		}
		Encode(h.buf[:], data[i:i+1])
		h.buf[2] = ' '
		l := 3
		if h.used == 7 {
			// There's an additional space after the 8th byte.
			h.buf[3] = ' '
			l = 4
		} else if h.used == 15 {
			// At the end of the line there's an extra space and
			// the bar for the right column.
			h.buf[3] = ' '
			h.buf[4] = '|'
			l = 5
		}
		_, err := h.w.Write(h.buf[:l])
		if err != nil {
			return n, err
		}
		n++
		h.rightChars[h.used] = toChar(data[i])
		h.used++
		h.n++
		if h.used == 16 {
			h.rightChars[16] = '|'
			h.rightChars[17] = '\n'
			_, err := h.w.Write(h.rightChars[:])
			if err != nil {
				return n, err
			}
			h.used = 0
		}
	}
	return n, nil
}

func (h *Dumper) Close() error {
	// See the comments in Write() for the details of this format.
	if h.closed {
		return nil
	}
	h.closed = true
	if h.used == 0 {
		return nil
	}
	h.buf[0] = ' '
	h.buf[1] = ' '
	h.buf[2] = ' '
	h.buf[3] = ' '
	h.buf[4] = '|'
	nBytes := h.used
	for h.used < 16 {
		l := 3
		if h.used == 7 {
			l = 4
		} else if h.used == 15 {
			l = 5
		}
		_, err := h.w.Write(h.buf[:l])
		if err != nil {
			return err
		}
		h.used++
	}
	h.rightChars[nBytes] = '|'
	h.rightChars[nBytes+1] = '\n'
	_, err := h.w.Write(h.rightChars[:nBytes+2])
	return err
}

func toChar(b byte) byte {
	if b < 32 || b > 126 {
		return '.'
	}
	return b
}
