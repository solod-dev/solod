/*
Package fmt implements formatted I/O with functions analogous
to C's printf and scanf. The format 'verbs' are the same as in C
(not the ones used in Go):

	%%	literal percent sign

	%d	integer, base 10, signed
	%u	integer, base 10, unsigned
	%o	integer, base 8, unsigned
	%x	integer, base 16, unsigned

	%f	floating-point, decimal notation
	%e	floating-point, decimal exponent notation
	%a	floating-point, hexadecimal exponent notation
	%g	floating-point, decimal or exponent notation as needed

	%c	single literal character
	%s	character string

	%p	pointer, base 16 notation, with leading 0x
*/
package fmt

import (
	_ "embed"

	"github.com/nalgeon/solod/so/errors"
	"github.com/nalgeon/solod/so/io"
)

//so:embed fmt.h
var fmt_h string

//so:embed fmt.c
var fmt_c string

// BufSize is the size of the internal formatting buffer in bytes.
//
//so:extern
const BufSize = 1024

//so:extern
var ErrPrint = errors.New("print failure")

//so:extern
var ErrScan = errors.New("scan failure")

//so:extern
var ErrSize = errors.New("buffer size exceeded")

// Print writes its arguments to standard output, separated by spaces.
// It returns the number of bytes written and any write error encountered.
//
//so:extern
func Print(a ...string) (int, error) {
	return 0, nil
}

// Println is like Print but adds a newline at the end.
//
//so:extern
func Println(a ...string) (int, error) {
	return 0, nil
}

// Printf formats according to a format specifier and writes to standard output.
// It returns the number of bytes written and any write error encountered.
//
//so:extern
func Printf(format string, a ...any) (int, error) {
	return 0, nil
}

// Fprintf formats according to a format specifier and writes to w.
// It returns the number of bytes written and any write error encountered.
// Returns [ErrSize] if the output size exceeds BufSize.
//
//so:extern
func Fprintf(w io.Writer, format string, a ...any) (n int, err error) {
	return 0, nil
}

// Scanf scans text read from standard input, storing successive
// space-separated values into successive arguments as determined by the format.
// It returns the number of items successfully scanned.
//
//so:extern
func Scanf(format string, a ...any) (n int, err error) {
	return 0, nil
}

// Sscanf scans the argument string, storing successive space-separated
// values into successive arguments as determined by the format.
// It returns the number of items successfully scanned.
//
//so:extern
func Sscanf(str string, format string, a ...any) (n int, err error) {
	return 0, nil
}

// Fscanf scans text read from r, storing successive space-separated
// values into successive arguments as determined by the format.
// It returns the number of items successfully scanned.
//
//so:extern
func Fscanf(r io.Reader, format string, a ...any) (int, error) {
	return 0, nil
}
