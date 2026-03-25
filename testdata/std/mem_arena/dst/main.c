#include "main.h"

// -- Implementation --

int main(void) {
    so_Slice buf = so_make_slice(so_byte, 1024, 1024);
    mem_Arena arena = mem_NewArena(buf);
    mem_Allocator a = (mem_Allocator){.self = &arena, .Alloc = mem_Arena_Alloc, .Free = mem_Arena_Free, .Realloc = mem_Arena_Realloc};
    // Allocate a Point.
    so_Result _res1 = mem_TryAlloc(main_Point, a);
    main_Point* p = _res1.val.as_ptr;
    so_Error err = _res1.err;
    if (err != NULL) {
        so_panic("initial allocation failed");
    }
    p->x = 11;
    p->y = 22;
    if (p->x != 11 || p->y != 22) {
        so_panic("unexpected p.x or p.y");
    }
    fmt_Println("alloc ok");
    // Free is a no-op.
    mem_Free(main_Point, a, p);
    // Reset and reallocate.
    mem_Arena_Reset(&arena);
    so_Result _res2 = mem_TryAlloc(main_Point, a);
    main_Point* p2 = _res2.val.as_ptr;
    err = _res2.err;
    if (err != NULL) {
        so_panic("allocation after reset failed");
    }
    // Memory should be zeroed.
    if (p2->x != 0 || p2->y != 0) {
        so_panic("memory not zeroed after reset");
    }
    p2->x = 33;
    p2->y = 44;
    fmt_Printf("reset ok: %d %d\n", p2->x, p2->y);
}
