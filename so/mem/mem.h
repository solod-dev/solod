// Alloc allocates a single value of type T using allocator a.
// Returns a pointer to the allocated memory or panics on failure.
// If the allocator is nil, uses the system allocator.
#define mem_Alloc(T, a) ({                     \
    so_Result _mem_res = mem_TryAlloc(T, (a)); \
    if (_mem_res.err != NULL)                  \
        so_panic(_mem_res.err->msg);           \
    _mem_res.val.as_ptr;                       \
})

// TryAlloc allocates memory for a single value of type T using allocator a.
// Returns a pointer to the allocated memory or an error if allocation fails.
// If the allocator is nil, uses the system allocator.
#define mem_TryAlloc(T, a) ({                            \
    mem_Allocator _a = (a);                              \
    if (!_a.self) _a = mem_System;                       \
    _a.Alloc(_a.self, sizeof(T), alignof(so_typeof(T))); \
})

// Free frees a value previously allocated with [Alloc] or [TryAlloc].
// If the allocator is nil, uses the system allocator.
#define mem_Free(T, a, ptr) ({                                 \
    mem_Allocator _a = (a);                                    \
    if (!_a.self) _a = mem_System;                             \
    _a.Free(_a.self, (ptr), sizeof(T), alignof(so_typeof(T))); \
})

// AllocSlice allocates a slice of type T with given length
// and capacity using allocator a.
// Returns a slice of the allocated memory or panics on failure.
// If the allocator is nil, uses the system allocator.
#define mem_AllocSlice(T, a, len, cap) ({                     \
    so_Result _res = mem_TryAllocSlice(T, (a), (len), (cap)); \
    if (_res.err != NULL)                                     \
        so_panic(_res.err->msg);                              \
    _res.val.as_slice;                                        \
})

// TryAllocSlice allocates a slice of type T with given length and capacity using allocator a.
// Returns a slice of the allocated memory or an error if allocation fails.
// If the allocator is nil, uses the system allocator.
#define mem_TryAllocSlice(T, a, slen, scap) ({                                    \
    mem_Allocator _a = (a);                                                       \
    if (!_a.self) _a = mem_System;                                                \
    if ((slen) > (scap)) so_panic("mem: length exceeds capacity");                \
    so_Result _mem_res = _a.Alloc(_a.self, sizeof(T) * (scap),                    \
                                  alignof(so_typeof(T)));                         \
    so_Slice _slice = {.ptr = _mem_res.val.as_ptr, .len = (slen), .cap = (scap)}; \
    so_Result _slice_res = {.val.as_slice = _slice, .err = _mem_res.err};         \
    _slice_res;                                                                   \
})

// FreeSlice frees a slice previously allocated with [AllocSlice] or [TryAllocSlice].
// If the allocator is nil, uses the system allocator.
#define mem_FreeSlice(T, a, s) ({                                        \
    mem_Allocator _a = (a);                                              \
    so_Slice _s = (s);                                                   \
    if (!_a.self) _a = mem_System;                                       \
    _a.Free(_a.self, _s.ptr, sizeof(T) * _s.cap, alignof(so_typeof(T))); \
})
