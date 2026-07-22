package mem

import (
	"unsafe"

	"solod.dev/so/c"
)

// Array is a fixed-length collection of equally-sized values stored inline in
// a single allocation. The element size is fixed at construction, and Store,
// Load and At operate through untyped (void*) pointers, so Array works with
// any element type without being generic. Use it as a building block for
// typed containers.
type Array struct {
	alloc Allocator
	vals  []byte
	vsize int
	count int
}

// NewArray allocates storage for count values of vsize bytes each.
// Both vsize and count must be greater than 0.
// Call [Array.Free] exactly once when done.
func NewArray(alloc Allocator, vsize int, count int) Array {
	c.Assert(vsize > 0, "mem.NewArray: vsize must be greater than 0")
	c.Assert(count > 0, "mem.NewArray: count must be greater than 0")
	vals := AllocSlice[byte](alloc, count*vsize, count*vsize)
	return Array{
		alloc: alloc,
		vals:  vals,
		vsize: vsize,
		count: count,
	}
}

// Load copies the value at index i into dst.
// dst must point to storage of at least vsize bytes.
func (a *Array) Load(i int, dst any) {
	Copy(dst, a.At(i), a.vsize)
}

// Store copies the value pointed to by v into slot i.
// v must point to storage of at least vsize bytes.
func (a *Array) Store(i int, v any) {
	Copy(a.At(i), v, a.vsize)
}

// At returns a pointer to the value at index i. The pointer stays
// valid until the slot is overwritten or [Array.Free] is called.
func (a *Array) At(i int) any {
	c.Assert(i >= 0 && i < a.count, "mem.Array.At: index out of range")
	vptr := unsafe.SliceData(a.vals)
	return c.PtrAdd(vptr, i*a.vsize)
}

// Len returns the number of values.
func (a *Array) Len() int {
	return a.count
}

// Free releases the memory allocated for the values.
// After calling Free, the Array is unusable.
func (a *Array) Free() {
	FreeSlice(a.alloc, a.vals)
	a.vals = nil
	a.count = 0
}
