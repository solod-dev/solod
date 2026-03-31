package time

import (
	"solod.dev/so/c"
	"solod.dev/so/errors"
)

// Commonly used layouts for Format and Parse.
const (
	RFC3339     = "%Y-%m-%dT%H:%M:%S%z"
	RFC3339Nano = "%Y-%m-%dT%H:%M:%S.%f%z"
	DateTime    = "%Y-%m-%d %H:%M:%S"
	DateOnly    = "%Y-%m-%d"
	TimeOnly    = "%H:%M:%S"
)

// ErrParse is returned by Parse when the input cannot be parsed.
var ErrParse = errors.New("time: cannot parse")

// Format formats the time per layout (strftime verbs like "%Y-%m-%d"),
// writing into buf. Returns the formatted string (a view into buf).
// buf length must be large enough for the formatted output plus a null terminator.
func (t Time) Format(buf []byte, layout string, offset Offset) string {
	date := t.Date(offset)
	clock := t.Clock(offset)
	var tm time_tm
	tm.tm_year = date.Year - 1900
	tm.tm_mon = int(date.Month) - 1
	tm.tm_mday = date.Day
	tm.tm_hour = clock.Hour
	tm.tm_min = clock.Minute
	tm.tm_sec = clock.Second
	tm.tm_wday = int(t.Weekday())
	tm.tm_yday = t.YearDay() - 1
	tm.tm_isdst = 0
	n := strftime(c.CharPtr(&buf[0]), uintptr(len(buf)), layout, &tm)
	return string(buf[:n])
}

// String formats the time as ISO 8601 "2006-01-02T15:04:05Z",
// writing into buf. Returns the formatted string (a view into buf).
// buf must have a length of at least 21 bytes.
func (t Time) String(buf []byte) string {
	return t.Format(buf, "%Y-%m-%dT%H:%M:%SZ", UTC)
}

// Parse parses value per layout (strptime verbs) and returns the Time.
// offset specifies what timezone the input value is in.
func Parse(layout string, value string, offset Offset) (Time, error) {
	var tm time_tm
	end := strptime(value, layout, &tm)
	if end == nil {
		return Time{}, ErrParse
	}
	return Date(int(tm.tm_year)+1900, Month(int(tm.tm_mon)+1), int(tm.tm_mday),
		int(tm.tm_hour), int(tm.tm_min), int(tm.tm_sec), 0, offset), nil
}
