package time

import "solod.dev/so/c"

//so:embed time.h
var time_h string

//so:extern
type time_tm struct {
	tm_sec   c.Int
	tm_min   c.Int
	tm_hour  c.Int
	tm_mday  c.Int
	tm_mon   c.Int
	tm_year  c.Int
	tm_wday  c.Int
	tm_yday  c.Int
	tm_isdst c.Int
}

// size_t strftime(char* str, size_t count, const char* format, const struct tm* tp);
//
//so:extern
func strftime(buf *c.Char, count uintptr, format string, tm *time_tm) uintptr {
	_, _, _, _ = buf, count, format, tm
	return 0
}

// char* strptime(const char* s, const char* format, struct tm* tm);
//
//so:extern
func strptime(value string, format string, tm *time_tm) *c.Char {
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

// time_sleep pauses the calling thread for at least ns nanoseconds.
//
//so:extern
func time_sleep(ns int64) { _ = ns }

// Monotonic times are reported as offsets from monoStart.
// We initialize monoStart to time_mono() - 1 so that on systems where
// monotonic time resolution is fairly low (e.g. Windows 2008
// which appears to have a default resolution of 15ms),
// we avoid ever reporting a monotonic time of 0.
// (Callers may want to use 0 as "time not set".)
var monoStart int64

func init() {
	monoStart = time_mono() - 1
}
