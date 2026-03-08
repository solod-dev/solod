// So provides the `so/mem` package for manual memory
// management. It supports heap allocation for single
// values and slices, and custom allocators for more
// control.
package main

import "github.com/nalgeon/solod/so/mem"

type Point struct {
	x, y int
}

func main() {
	// `mem.New` allocates a single value on the heap.
	// Always pair it with `mem.Free` to avoid leaks.
	p := mem.New[Point]()
	defer mem.Free(p)

	p.x = 10
	p.y = 20
	println("point:", p.x, p.y)

	// `mem.NewSlice` allocates a slice on the heap
	// with a given length and capacity.
	nums := mem.NewSlice[int](5, 5)
	defer mem.FreeSlice(nums)

	for i := range nums {
		nums[i] = i * 10
	}
	println("slice:", nums[0], nums[1], nums[2], nums[3], nums[4])

	// `mem.Alloc` uses a custom allocator and returns
	// an error if allocation fails.
	q, err := mem.Alloc[Point](mem.System)
	if err != nil {
		panic(err)
	}
	defer mem.Dealloc(mem.System, q)

	q.x = 30
	q.y = 40
	println("allocated:", q.x, q.y)
}
