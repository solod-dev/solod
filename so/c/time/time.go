// Package time wraps the C <time.h> header.
// It offers calendar time, broken-down time, and time formatting.
package time

import _ "embed"

//so:embed time.h
var time_h string

// TimeT represents a calendar time value (seconds since epoch -
// 1970-01-01 00:00:00 UTC). Does not count leap seconds.
//
//so:extern
type TimeT int64

// Tm represents a broken-down time with individual date/time components.
//
//so:extern
type Tm struct {
	Sec   int // seconds after the minute [0, 60]
	Min   int // minutes after the hour [0, 59]
	Hour  int // hours since midnight [0, 23]
	Mday  int // day of the month [1, 31]
	Mon   int // months since January [0, 11]
	Year  int // years since 1900
	Wday  int // days since Sunday [0, 6]
	Yday  int // days since January 1 [0, 365]
	Isdst int // daylight saving time flag
}

// ClocksPerSec is the number of CPU clock ticks per second.
//
//so:extern
var ClocksPerSec int

// Time returns the current calendar time.
// If timer is not nil, the return value is also stored in *timer.
// Returns TimeT(-1) on error.
//
//so:extern
func Time(timer *TimeT) TimeT { _ = timer; return 0 }

// Clock returns the approximate processor time used by the program
// since it started, in clock ticks. To convert result value to seconds,
// divide it by [ClocksPerSec].
//
// If the processor time is not available, returns -1.
//
//so:extern
func Clock() int { return 0 }

// Difftime returns the difference between two calendar times
// (end - start) in seconds. If end is before start, the result is negative.
//
//so:extern
func Difftime(end TimeT, start TimeT) float64 { _, _ = end, start; return 0 }

// Gmtime converts a calendar time to a broken-down time
// expressed as UTC.
//
//so:extern
func Gmtime(timer *TimeT) Tm { _ = timer; return Tm{} }

// Mktime converts a broken-down time to a calendar time value.
// It also normalizes the fields of timeptr (e.g. Mon=13 becomes
// the next year, etc.).
//
// Returns -1 if the calendar time cannot be represented.
//
//so:extern
func Mktime(timeptr *Tm) TimeT { _ = timeptr; return 0 }

// Strftime formats the broken-down time timeptr into buf
// according to the format string, writing at most maxsize bytes.
// Returns the number of bytes written (not counting the terminating
// null character), or 0 if the result would exceed maxsize bytes.
//
//so:extern
func Strftime(buf *byte, maxsize int, format string, timeptr *Tm) int {
	_, _, _, _ = buf, maxsize, format, timeptr
	return 0
}
