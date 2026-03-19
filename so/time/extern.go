package time

import _ "embed"

//so:embed time.h
var time_h string

//so:embed time.c
var time_c string

//so:extern
type time_tm struct {
	tm_sec   int
	tm_min   int
	tm_hour  int
	tm_mday  int
	tm_mon   int
	tm_year  int
	tm_wday  int
	tm_yday  int
	tm_isdst int
}

//so:extern
func strftime(buf *byte, count uintptr, format string, tm *time_tm) uintptr {
	_, _, _, _ = buf, count, format, tm
	return 0
}

//so:extern
func strptime(value string, format string, tm *time_tm) any {
	_, _, _ = value, format, tm
	return nil
}

// wall returns the current wall clock time.
//
//so:extern
func time_wall() (int64, int32) { return 0, 0 }

// mono returns the current monotonic time in nanoseconds.
//
//so:extern
func time_mono() int64 { return 0 }

// Monotonic times are reported as offsets from monoStart.
// We initialize monoStart to time_mono() - 1 so that on systems where
// monotonic time resolution is fairly low (e.g. Windows 2008
// which appears to have a default resolution of 15ms),
// we avoid ever reporting a monotonic time of 0.
// (Callers may want to use 0 as "time not set".)
//
//so:extern
var time_monoStart int64 = time_mono() - 1
