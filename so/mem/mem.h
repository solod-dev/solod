#include "so/builtin/builtin.h"

// SwapByte swaps n bytes between a and b.
// Panics if either a or b is nil.
//
// SwapByte temporarily allocates a buffer of size n
// on the stack, so it's not suitable for large n.
static inline void mem_SwapByte(void* a, void* b, so_int n) {
    assert(a != NULL && "mem: nil pointer");
    assert(b != NULL && "mem: nil pointer");
    assert(n >= 0 && "mem: negative size");
    if (n == 0) return;

    size_t size = (size_t)n;
    char tmp[size];
    memcpy(tmp, a, size);
    memcpy(a, b, size);
    memcpy(b, tmp, size);
}

#if !__STDC_HOSTED__

// Bump allocator over a static buffer for freestanding environments.
// Memory is never reclaimed: free is a no-op, realloc copies into a new bump.
// Suitable for short-lived programs that don't need much memory.
// The heap is off by default, enable with -DSO_HEAP_SIZE=N.

#ifndef SO_HEAP_SIZE
#define SO_HEAP_SIZE (0)  // in bytes
#endif

#if SO_HEAP_SIZE > 0

static char so_heap[SO_HEAP_SIZE];
static size_t so_heap_offset = 0;

static inline void* malloc(size_t size) {
    if (size == 0) return NULL;
    // Align to 16 bytes.
    so_heap_offset = (so_heap_offset + 15) & ~(size_t)15;
    if (so_heap_offset + size > SO_HEAP_SIZE) return NULL;
    void* ptr = &so_heap[so_heap_offset];
    so_heap_offset += size;
    return ptr;
}

static inline void* calloc(size_t num, size_t size) {
    if (num != 0 && size > SIZE_MAX / num) return NULL;
    size_t total = num * size;
    void* ptr = malloc(total);
    if (ptr) memset(ptr, 0, total);
    return ptr;
}

static inline void* realloc(void* ptr, size_t new_size) {
    if (new_size == 0) return NULL;
    void* new_ptr = malloc(new_size);
    if (ptr && new_ptr) {
        // We don't track allocation sizes, so we copy new_size bytes.
        // When growing, this over-reads from the old allocation into
        // adjacent bump memory - harmless but yields garbage in the tail.
        memcpy(new_ptr, ptr, new_size);
    }
    return new_ptr;
}

#else

static inline void* malloc(size_t size) {
    (void)size;
    return NULL;
}

static inline void* calloc(size_t num, size_t size) {
    (void)num;
    (void)size;
    return NULL;
}

static inline void* realloc(void* ptr, size_t new_size) {
    (void)ptr;
    (void)new_size;
    return NULL;
}

#endif  // SO_HEAP_SIZE > 0

static inline void free(void* ptr) {
    (void)ptr;
}

#endif  // !__STDC_HOSTED__
