#include "so/builtin/builtin.h"

#if !__STDC_HOSTED__
static inline void* memchr(const void* s, int c, size_t n) {
    const unsigned char* p = s;
    unsigned char target = (unsigned char)c;
    while (n--) {
        if (*p == target) return (void*)p;
        p++;
    }
    return NULL;
}
#endif

static inline so_int bytealg_Compare(so_Slice a, so_Slice b) {
    so_int n = a.len;
    if (b.len < n) n = b.len;
    int cmp = memcmp(a.ptr, b.ptr, (size_t)n);
    if (cmp != 0) return cmp;
    if (a.len < b.len) return -1;
    if (a.len > b.len) return +1;
    return 0;
}

static inline so_int bytealg_IndexByte(so_Slice b, so_byte c) {
    void* at = memchr(b.ptr, (int)c, (size_t)b.len);
    if (at == NULL) return -1;
    return (so_int)((char*)at - (char*)b.ptr);
}
