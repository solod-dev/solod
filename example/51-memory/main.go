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

	// `mem.Arena` uses a pre-allocated buffer for fast, linear allocation.
	// Individual `Free` calls are no-ops - call `Reset` to reclaim
	// all memory at once. If the backing buffer is heap-allocated,
	// free it with `mem.FreeSlice` when done.
	//
	// Arenas are ideal for short-lived allocations with simple lifetimes,
	// such as during parsing or temporary buffers.
	buf := make([]byte, 1024)
	arena := mem.NewArena(buf)
	var a mem.Allocator = &arena

	// Allocate 10 points (16B x 10 = 160B) from the arena.
	points := make([]*Point, 10)
	for i := range points {
		pt := mem.Alloc[Point](a)
		pt.x = i + 1
		pt.y = (i + 1) * 2
		points[i] = pt
	}
	println("allocated", len(points), "points in arena")
	println("points[0]:", points[0].x, points[0].y)
	println("points[9]:", points[9].x, points[9].y)

	// Reset reclaims all arena memory for reuse.
	arena.Reset()

	// If the buffer was heap-allocated, we'd need to free it here:
	// mem.FreeSlice(nil, buf)
}
