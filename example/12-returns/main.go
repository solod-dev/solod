// So has a limited support for _multiple return values_.
// This feature can only be used to return both result
// and error values from a function.
package main

// The `(int, error)` in this function signature shows that
// the function returns an `int` and an `error`.
func vals() (int, error) {
	return 3, nil
}

func main() {
	// Here we use the 2 different return values from the
	// call with _multiple assignment_.
	v, err := vals()
	println("v =", v)
	println("err =", err)

	// If you only want a subset of the returned values,
	// use the blank identifier `_`.
	_, err = vals()
	println("err =", err)
}
