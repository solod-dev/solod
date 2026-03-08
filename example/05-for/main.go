// `for` is So's only looping construct. Here are
// some basic types of `for` loops.
package main

func main() {
	// The most basic type, with a single condition.
	i := 1
	for i <= 3 {
		println(i)
		i = i + 1
	}

	// A classic initial/condition/after `for` loop.
	for j := 0; j < 3; j++ {
		println(j)
	}

	// Another way of accomplishing the basic "do this
	// N times" iteration is `range` over an integer.
	for i := range 3 {
		println("range", i)
	}

	// `for` without a condition will loop repeatedly
	// until you `break` out of the loop or `return` from
	// the enclosing function.
	for {
		println("loop")
		break
	}

	// You can also `continue` to the next iteration of
	// the loop.
	for n := range 6 {
		if n%2 == 0 {
			continue
		}
		println(n)
	}
}
