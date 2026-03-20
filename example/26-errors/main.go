// In So it's idiomatic to communicate errors via an
// explicit, separate return value. This contrasts with
// the exceptions used in languages like Java, Python and
// Ruby and the overloaded single result / error value
// sometimes used in C. So's approach makes it easy to
// see which functions return errors and to handle them
// using the same language constructs employed for other,
// non-error tasks.
package main

import "solod.dev/so/errors"

// A sentinel error is a predeclared variable that
// is used to signify a specific error condition.
var ErrOutOfTea = errors.New("no more tea available")
var ErrPower = errors.New("can't boil water")

// Unlike Go, So only suppors sentinel errors.
// They must be defined at the package level using
// errors.New with a plain string message.
var Err42 = errors.New("can't work with 42")

// By convention, errors are the last return value and
// have type `error`, a built-in interface.
func f(arg int) (int, error) {
	if arg == 42 {
		return -1, Err42
	}

	// A `nil` value in the error position indicates that
	// there was no error.
	return arg + 3, nil
}

func makeTea(arg int) error {
	if arg == 2 {
		return ErrOutOfTea
	} else if arg == 4 {
		return ErrPower
	}
	return nil
}

func main() {
	s := []int{7, 42}
	for _, i := range s {
		// It's idiomatic to use an inline
		// error check in the `if` line.
		if r, e := f(i); e != nil {
			println("f failed:", e)
		} else {
			println("f worked:", r)
		}
	}

	for i := range 5 {
		if err := makeTea(i); err != nil {
			// Since So only supports sentinel errors,
			// they are compared with a simple pointer equality check.
			// No errors.Is or errors.As needed (or supported).
			if err == ErrOutOfTea {
				println("We should buy new tea!")
			} else if err == ErrPower {
				println("Now it is dark.")
			} else {
				println("unknown error:", err)
			}
			continue
		}

		println("Tea is ready!")
	}
}
