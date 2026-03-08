// In So, an _array_ is a numbered sequence of elements of a
// specific length. In typical So code, _slices_ are
// much more common; arrays are useful in some special
// scenarios.
package main

import "github.com/nalgeon/solod/so/c/stdio"

func main() {
	// Here we create an array `a` that will hold exactly
	// 5 `int`s. The type of elements and length are both
	// part of the array's type. By default an array is
	// zero-valued, which for `int`s means `0`s.
	var a [3]int32
	printArray("emp:", a)

	// We can set a value at an index using the
	// `array[index] = value` syntax, and get a value with
	// `array[index]`.
	a[2] = 100
	printArray("set:", a)
	println("get:", a[2])

	// The builtin `len` returns the length of an array.
	println("len:", len(a))

	// Use this syntax to declare and initialize an array
	// in one line.
	b := [3]int32{1, 2, 3}
	printArray("dcl:", b)

	// You can also have the compiler count the number of
	// elements for you with `...`
	b = [...]int32{1, 2, 3}
	printArray("dcl:", b)

	// If you specify the index with `:`, the elements in
	// between will be zeroed.
	b = [...]int32{2: 300}
	printArray("idx:", b)

	// Array types are one-dimensional, but you can
	// compose types to build multi-dimensional data
	// structures.
	var twoD [2][3]int32
	for i := range 2 {
		for j := range 3 {
			twoD[i][j] = int32(i + j)
		}
	}
	stdio.Printf("2d: [%d...%d]\n", twoD[0][0], twoD[1][2])

	// You can create and initialize multi-dimensional
	// arrays at once too.
	twoD = [2][3]int32{
		{1, 2, 3},
		{1, 2, 3},
	}
	stdio.Printf("2d: [%d...%d]\n", twoD[0][0], twoD[1][2])
}

func printArray(msg string, arr [3]int32) {
	stdio.Printf("%s [%d %d %d]\n", msg, arr[0], arr[1], arr[2])
}
