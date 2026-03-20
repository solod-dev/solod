// The `so/time` package provides support for
// times and durations; here are some examples.
package main

import "solod.dev/so/time"

func main() {
	// Pre-allocate a buffer for string formatting.
	buf := make([]byte, 64)

	// We'll start by getting the current time.
	now := time.Now()
	println("now", now.String(buf))

	// You can build a Time struct by providing the
	// year, month, day, etc. Times are always in UTC.
	// You can provide a UTC offset in seconds, and the
	// time will be adjusted accordingly.
	then := time.Date(2026, 3, 17, 20, 34, 58, 651387237, time.UTC)
	println("then", then.String(buf))

	// You can extract the various components of the time
	// value as expected.
	println("year", then.Year())
	println("month", then.Month())
	println("day", then.Day())
	println("hour", then.Hour())
	println("minute", then.Minute())
	println("second", then.Second())
	println("nano", then.Nanosecond())

	// The Monday-Sunday `Weekday` is also available.
	// Monday is 0, Tuesday is 1, and so on.
	println("weekday", then.Weekday())

	// These methods compare two times, testing if the
	// first occurs before, after, or at the same time
	// as the second, respectively.
	println("then before now", then.Before(now))
	println("then after now", then.After(now))
	println("then equal now", then.Equal(now))

	// The `Sub` methods returns a `Duration` representing
	// the interval between two times.
	diff := now.Sub(then)
	println("diff", diff.String(buf))

	// We can compute the length of the duration in
	// various units.
	println("hours", diff.Hours())
	println("minutes", diff.Minutes())
	println("seconds", diff.Seconds())
	println("nanos", diff.Nanoseconds())

	// You can use `Add` to advance a time by a given
	// duration, or with a `-` to move backwards by a
	// duration.
	println("then + diff", then.Add(diff).String(buf))
	println("then - diff", then.Add(-diff).String(buf))
}
