// So provides the `so/mem` package for manual memory
// management. It supports heap allocation for single
// values and slices, and custom allocators for more
// control.
package main

import "solod.dev/so/mem"

type Point struct {
	x, y int
}

func main() {
	// `mem.Alloc` allocates a single value on the heap.
	// Always pair it with `mem.Free` to avoid leaks.
	// The first argument is an optional allocator;
	// passing `nil` uses the default allocator.
	p := mem.Alloc[Point](nil)
	defer mem.Free(nil, p)

	p.x = 10
	p.y = 20
	println("point:", p.x, p.y)

	// `mem.AllocSlice` allocates a slice on the heap
	// with a given length and capacity.
	nums := mem.AllocSlice[int](nil, 5, 5)
	defer mem.FreeSlice(nil, nums)

	for i := range nums {
		nums[i] = i * 10
	}
	println("slice:", nums[0], nums[1], nums[2], nums[3], nums[4])

	// `mem.TryAlloc` returns an error if allocation fails
	// instead of panicking.
	q, err := mem.TryAlloc[Point](mem.System)
	if err != nil {
		panic(err)
	}
	defer mem.Free(mem.System, q)
	q.x = 30
	q.y = 40
	println("allocated:", q.x, q.y)
}
