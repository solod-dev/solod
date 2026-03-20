// A `panic` typically means something went unexpectedly
// wrong. Mostly we use it to fail fast on errors that
// shouldn't occur during normal operation, or that we
// aren't prepared to handle gracefully.
package main

import "solod.dev/so/errors"

var Err42 = errors.New("got 42")

func work(n int) (int, error) {
	if n == 42 {
		return 0, Err42
	}
	return n, nil
}

func main() {
	// A common use of panic is to abort if a function
	// returns an error value that we don't know how to
	// (or want to) handle. Here's an example of
	// panicking if we get an unexpected error.
	n, err := work(42)
	if err != nil {
		panic(err)
	}
	println(n)
}
