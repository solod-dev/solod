#include "main.h"

// -- Variables and constants --

// -- Implementation --

int main(void) {
    {
        // Function with return.
        so_int x = identity(so_int, 42);
        if (x != (so_int)(42)) {
            so_panic("Function with return failed");
        }
    }
    {
        // Function w/o return.
        so_int y = 0;
        setPtr(so_int, &y, 42);
        if (y != 42) {
            so_panic("Function w/o return failed");
        }
    }
    {
        // Nested calls with variable shadowing.
        so_int z = a(so_int, 42);
        if (z != 45) {
            so_panic("Nested calls failed");
        }
    }
    {
        // Generic method.
        main_Box b = {0};
        main_Box_set(so_int, &b, 42);
        if (b.val != 42) {
            so_panic("Generic method failed");
        }
    }
    so_println("%s", "lang/macro ok");
}
