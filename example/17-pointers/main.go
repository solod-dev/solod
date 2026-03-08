// So supports _pointers_, allowing you to pass references
// to values and records within your program.
package main

// We'll show how pointers work in contrast to values with
// 2 functions: `zeroval` and `zeroptr`. `zeroval` has an
// `int` parameter, so arguments will be passed to it by
// value. `zeroval` will get a copy of `ival` distinct
// from the one in the calling function.
func zeroval(ival int) {
	ival = 0
}

// `zeroptr` in contrast has an `*int` parameter, meaning
// that it takes an `int` pointer. The `*iptr` code in the
// function body then _dereferences_ the pointer from its
// memory address to the current value at that address.
// Assigning a value to a dereferenced pointer changes the
// value at the referenced address.
func zeroptr(iptr *int) {
	*iptr = 0
}

func main() {
	i := 1
	println("initial:", i)

	zeroval(i)
	println("zeroval:", i)

	// The `&i` syntax gives the memory address of `i`,
	// i.e. a pointer to `i`.
	zeroptr(&i)
	println("zeroptr:", i)

	// Pointers can be printed too.
	println("pointer:", &i)

	// any translates to (void*) in C,
	// and can hold any pointer value.
	var n byte = 15
	var a any = n // void* a = &n
	println("any pointer:", a)

	// any can be converted back to a pointer of the original type,
	// or to a pointer of a different type (at your own risk).
	b := a.(*byte) // uint32_t* b = (uint32_t*)a
	println("any as byte:", *b)
	i32 := a.(*int32) // int32_t* i32 = (int32_t*)a
	println("any as int32:", *i32)
}
