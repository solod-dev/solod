// So can call C standard library functions by using
// `so:include`, `so:extern`, `so:embed` and declaring
// functions without a body.
package main

import (
	_ "embed"

	"github.com/nalgeon/solod/so/c"
)

// The include directive tell So to #include
// the given headers into the generated C code.
//so:include <stdlib.h>
//so:include <time.h>

// Declare C functions to make them callable from So.
// String arguments auto-decay to C's `char*`.
func getenv(name string) *byte
func atoi(s string) int32

// Scalar types like int32 map directly to C's int.
func abs(n int32) int32

// `time_t` is a C type for representing time values.
// We declare it as extern so that So can access it,
// but doesn't codegen it, since the type is already
// defined in the C header.
//
//so:extern
type time_t int64

// `time` returns the current calendar time
// as seconds since the Unix epoch.
func time(timer *time_t) time_t

// `struct tm` in C has no typedef, so we add one
// via a small embedded header to bridge the gap.
//
//so:embed interop.h
var interop_h string

// `tm` is C's broken-down time structure.
//
//so:extern
type tm struct{}

// `localtime` converts a time value to local time.
func localtime(timer *time_t) *tm

// `strftime` formats a time value into a string.
func strftime(buf *byte, maxsize int, format string, timeptr *tm) int

// `difftime` computes the difference between two times.
func difftime(end, start time_t) float64

func main() {
	// Read an environment variable.
	// `getenv` returns a C string (*byte), so use
	// `c.String` to convert it to a So string.
	home := c.String(getenv("HOME"))
	println("home directory:", home)

	// Convert a string to a number.
	n := atoi("12345")
	println("atoi:", n)

	// Compute an absolute value.
	println("abs(-42):", abs(-42))

	// Get the current time.
	now := time(nil)
	println("seconds since epoch:", now)

	// Format the current time as a string.
	// `localtime` returns a *tm (C's broken-down time).
	// `strftime` writes the formatted time into a buffer.
	local := localtime(&now)
	var buf [64]byte
	strftime(c.CharPtr(&buf[0]), 64, "%Y-%m-%d %H:%M:%S", local)
	println("current time:", c.String(&buf[0]))

	// Compute the difference between two times.
	start := now
	end := start + 3600
	diff := difftime(end, start)
	println("seconds in an hour:", diff)
}
