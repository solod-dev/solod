#include "main.h"

// -- Implementation --

int main(void) {
    {
        // Append within capacity.
        so_Slice s = mem_AllocSlice(so_int, (mem_Allocator){0}, 0, 8);
        s = slices_Append(so_int, (mem_Allocator){0}, s, 10, 20, 30);
        if (so_len(s) != 3 || so_at(so_int, s, 0) != 10 || so_at(so_int, s, 1) != 20 || so_at(so_int, s, 2) != 30) {
            so_panic("Append: unexpected value");
        }
        mem_FreeSlice(so_int, (mem_Allocator){0}, s);
    }
    {
        // Append that triggers growth.
        so_Slice s = mem_AllocSlice(so_int, (mem_Allocator){0}, 0, 2);
        s = slices_Append(so_int, (mem_Allocator){0}, s, 1, 2);
        s = slices_Append(so_int, (mem_Allocator){0}, s, 3, 4, 5);
        if (so_len(s) != 5 || so_at(so_int, s, 0) != 1 || so_at(so_int, s, 4) != 5) {
            so_panic("Append grow: unexpected value");
        }
        mem_FreeSlice(so_int, (mem_Allocator){0}, s);
    }
    {
        // Extend from another slice.
        so_Slice s = mem_AllocSlice(so_int, (mem_Allocator){0}, 0, 8);
        so_Slice other = (so_Slice){(so_int[3]){100, 200, 300}, 3, 3};
        s = slices_Extend(so_int, (mem_Allocator){0}, s, other);
        if (so_len(s) != 3 || so_at(so_int, s, 0) != 100 || so_at(so_int, s, 2) != 300) {
            so_panic("Extend: unexpected value");
        }
        mem_FreeSlice(so_int, (mem_Allocator){0}, s);
    }
    {
        // Clone a slice.
        so_Slice s1 = (so_Slice){(so_int[3]){11, 22, 33}, 3, 3};
        so_Slice s2 = slices_Clone(so_int, (mem_Allocator){0}, s1);
        so_at(so_int, s2, 0) = 99;
        if (so_at(so_int, s1, 0) != 11 || so_at(so_int, s2, 0) != 99) {
            so_panic("Clone: unexpected value");
        }
        mem_FreeSlice(so_int, (mem_Allocator){0}, s2);
    }
}
