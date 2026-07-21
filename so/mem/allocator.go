package mem

// Allocator defines the interface for memory allocators.
// Whether allocated or reallocated memory is zeroed is allocator-specific.
type Allocator interface {
	// Alloc allocates a block of memory of the given size and alignment.
	Alloc(size int, align int) (any, error)
	// Realloc resizes a previously allocated block of memory.
	Realloc(ptr any, oldSize int, newSize int, align int) (any, error)
	// Free frees a previously allocated block of memory.
	Free(ptr any, size int, align int)
}

// A Stats records statistics about the memory allocator.
// The number of live objects is Mallocs - Frees.
type Stats struct {
	// Alloc is bytes of allocated heap objects.
	Alloc uint64

	// TotalAlloc is cumulative bytes allocated for heap objects.
	//
	// TotalAlloc increases as heap objects are allocated, but unlike Alloc,
	// it does not decrease when objects are freed.
	TotalAlloc uint64

	// Mallocs is the cumulative count of heap objects allocated.
	Mallocs uint64

	// Frees is the cumulative count of heap objects freed.
	Frees uint64
}
