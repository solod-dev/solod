#include "main.h"

// -- Types --

typedef struct point point;

typedef struct point {
    so_int x;
    so_int y;
} point;

// -- Variables and constants --
static const int64_t ptrSize = ((int64_t)4 << ((uint64_t)(~(uintptr_t)(0)) >> 63));

// -- Implementation --

int main(void) {
    {
        // Sizeof.
        so_int x = 42;
        uintptr_t size = unsafe_Sizeof(x);
        if (size != ptrSize) {
            so_panic("invalid sizeof(int)");
        }
        point p = (point){1, 2};
        size = unsafe_Sizeof(p);
        if (size != 2 * ptrSize) {
            so_panic("invalid sizeof(point)");
        }
    }
    {
        // Alignof.
        so_int x = 42;
        uintptr_t align = unsafe_Alignof(x);
        if (align != ptrSize) {
            so_panic("invalid alignof(int)");
        }
        point p = (point){1, 2};
        align = unsafe_Alignof(p);
        if (align != ptrSize) {
            so_panic("invalid alignof(point)");
        }
    }
    // {
    // 	// Offsetof is not supported.
    // 	var p = point{1, 2}
    // 	offsetX := unsafe.Offsetof(p.x)
    // 	if offsetX != 0 {
    // 		panic("want offsetX == 0")
    // 	}
    // 	offsetY := unsafe.Offsetof(p.y)
    // 	if offsetY != 8 {
    // 		panic("want offsetY == 8")
    // 	}
    // }
    {
        // String.
        so_Slice b = so_string_bytes(so_str("hello"));
        so_String s = unsafe_String(&so_at(so_byte, b, 0), so_len(b));
        if (so_string_ne(s, so_str("hello"))) {
            so_panic("want s == 'hello'");
        }
    }
    {
        // StringData.
        so_String s = so_str("hello");
        so_byte* b = unsafe_StringData(s);
        if (*b != 'h') {
            so_panic("want *b == 'h'");
        }
    }
    {
        // Slice.
        so_int a[5] = {1, 2, 3, 4, 5};
        so_Slice slice = unsafe_Slice(&a[0], 5);
        if (so_len(slice) != 5) {
            so_panic("want len(slice) == 5");
        }
        if (so_at(so_int, slice, 0) != 1 || so_at(so_int, slice, 4) != 5) {
            so_panic("want slice[0] == 1 and slice[4] == 5");
        }
    }
    {
        // SliceData.
        so_Slice s = (so_Slice){(so_int[5]){1, 2, 3, 4, 5}, 5, 5};
        so_int* p = unsafe_SliceData(s);
        if (*p != 1) {
            so_panic("want *p == 1");
        }
    }
    {
        // Pointer.
        so_int x = 42;
        void* p = (void*)(&x);
        if (*(so_int*)(p) != 42) {
            so_panic("want *(int*)p == 42");
        }
    }
}
