// _Slices_ are an important data type in Go, giving
// a more powerful interface to sequences than arrays.
package main

import "solod.dev/so/c/stdio"

func main() {
	// Unlike arrays, slices are typed only by the
	// elements they contain (not the number of elements).
	// An uninitialized slice equals to nil and has
	// length 0.
	var s []string
	printSlice("uninit:", s)
	println("len:", len(s), "cap:", cap(s))

	// To create a slice with non-zero length, use
	// the builtin `make`. Here we make a slice of
	// `string`s of length `3` (initially zero-valued).
	// By default a new slice's capacity is equal to its
	// length; if we know the slice is going to grow,
	// we must pass a capacity as an additional parameter
	// to `make` (`6` in this case).
	s = make([]string, 3, 6)
	printSlice("emp:", s)
	println("len:", len(s), "cap:", cap(s))

	// We can set and get just like with arrays.
	s[0] = "a"
	s[1] = "b"
	s[2] = "c"
	printSlice("set:", s)
	println("get:", s[2])

	// `len` returns the length of the slice as expected.
	println("len:", len(s))

	// In addition to these basic operations, slices
	// support several more that make them richer than
	// arrays. One is the builtin `append`, which
	// returns a slice containing one or more new values.
	// Note that we need to accept a return value from
	// `append` as we may get a new slice value.
	s = append(s, "d")
	s = append(s, "e", "f")
	printSlice("apd:", s)

	// By default, slices in So are stack-allocated, and
	// cannot grow beyond their initial capacity.
	// To allocate a slice on the heap and grow it beyond
	// its initial capacity, use the `so/mem` package.

	// Slices can also be `copy`'d. Here we create an
	// empty slice `c` of the same length as `s` and copy
	// into `c` from `s`.
	c := make([]string, len(s))
	copy(c, s)
	printSlice("cpy:", c)

	// Slices support a "slice" operator with the syntax
	// `slice[low:high]`. For example, this gets a slice
	// of the elements `s[2]`, `s[3]`, and `s[4]`.
	l := s[2:5]
	printSlice("sl1:", l)

	// This slices up to (but excluding) `s[5]`.
	l = s[:5]
	printSlice("sl2:", l)

	// And this slices up from (and including) `s[2]`.
	l = s[2:]
	printSlice("sl3:", l)

	// We can declare and initialize a variable for slice
	// in a single line as well.
	t := []string{"g", "h", "i"}
	printSlice("dcl:", t)
}

func printSlice(msg string, s []string) {
	if len(s) == 0 {
		stdio.Printf("%s []\n", msg)
		return
	}
	stdio.Printf("%s", msg)
	stdio.Printf(" [")
	for i, v := range s {
		stdio.Printf("%s", v)
		if i < len(s)-1 {
			stdio.Printf(" ")
		}
	}
	stdio.Printf("]\n")
}
