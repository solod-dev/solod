package mem

import "unsafe"

//so:extern
var maxAllocSize = 1 << 20 // 1 MiB, for testing only

// for testing only
//
//so:extern
var allocs = map[uintptr]int{} // ptr -> size, for testing only

// void* calloc(size_t num, size_t size);
//
//so:extern
func calloc(count uintptr, size uintptr) any {
	total := count * size
	if total > uintptr(maxAllocSize) {
		return nil
	}
	s := make([]byte, total)
	ptr := &s[0]
	allocs[uintptr(unsafe.Pointer(ptr))] = int(total)
	return ptr
}

// void* malloc(size_t size);
//
//so:extern
func malloc(size uintptr) any {
	if size > uintptr(maxAllocSize) {
		return nil
	}
	s := make([]byte, size)
	ptr := &s[0]
	allocs[uintptr(unsafe.Pointer(ptr))] = int(size)
	return ptr
}

// void* realloc(void* ptr, size_t new_size);
//
//so:extern
func realloc(ptr any, newSize uintptr) any {
	if newSize > uintptr(maxAllocSize) {
		return nil
	}
	s := make([]byte, newSize)
	newPtr := &s[0]
	if ptr != nil {
		oldPtr := ptrVal(ptr)
		oldSize := allocs[uintptr(oldPtr)]
		copySize := min(oldSize, int(newSize))
		if copySize > 0 {
			copy(s, unsafe.Slice((*byte)(oldPtr), copySize))
		}
		delete(allocs, uintptr(oldPtr))
	}
	allocs[uintptr(unsafe.Pointer(newPtr))] = int(newSize)
	return newPtr
}

// void free(void *ptr);
//
//so:extern
func free(ptr any) {
	if ptr != nil {
		delete(allocs, uintptr(ptrVal(ptr)))
	}
}

// ptrVal extracts a raw pointer from an interface containing any pointer type.
// For testing only; in C, any pointers are void*.
//
//so:extern
func ptrVal(v any) unsafe.Pointer {
	type iface struct {
		_    unsafe.Pointer
		data unsafe.Pointer
	}
	return (*iface)(unsafe.Pointer(&v)).data
}
