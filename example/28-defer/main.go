// _Defer_ is used to ensure that a function call is
// performed later in a program's execution, usually for
// purposes of cleanup. `defer` is often used where e.g.
// `ensure` and `finally` would be used in other languages.
package main

import "solod.dev/so/mem"

type Point struct {
	x, y int
}

// Suppose we wanted to allocate an object on the heap,
// use it, and then deallocate it when we're done.
// Here's how we could do that with `defer`.
func main() {
	// Immediately after allocating a Point object on the heap,
	// we defer the deallocation of that object. This will be executed
	// at the end of the enclosing function (`main`), after we're done
	// using the object.
	p := mem.Alloc[Point](nil) // p is allocated on the heap, not on the stack
	defer mem.Free(nil, p)

	p.x = 11
	p.y = 22
	println(p.x, p.y)
}
