#include "so/builtin/builtin.h"

#define so_const_char const char

#define c_Alignof(T) ((so_int)alignof(T))

#define c_Alloca(T, n) ((T*)so_alloca(sizeof(T) * (size_t)(n)))

static inline void c_Assert(bool cond, const char* msg) {
    assert((cond) && msg);
}

static inline so_Slice c_Bytes(void* ptr, so_int n) {
    return ptr ? (so_Slice){ptr, n, n} : (so_Slice){0};
}

static inline char* c_CharPtr(void* ptr) {
    return (char*)ptr;
}

#define c_PtrAdd(T, ptr, offset) ((ptr) + (size_t)(offset))

#define c_PtrAs(T, ptr) ((T*)(ptr))

#define c_PtrAt(T, ptr, index) (&(ptr)[(index)])

#define c_Sizeof(T) ((so_int)sizeof(T))

#define c_Slice(T, ptr, len, cap) \
    (ptr ? (so_Slice){(ptr), (len), (cap)} : (so_Slice){0})

#define c_String(T, ptr) ({                                            \
    const char* _ptr = (const char*)(ptr);                             \
    (_ptr ? (so_String){_ptr, (so_int)strlen(_ptr)} : (so_String){0}); \
})

#define c_Zero(T) ((T){0})
