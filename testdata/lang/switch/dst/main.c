#include "main.h"

// -- Implementation --

int main(void) {
    {
        // Empty switch statement.
        if (false) {
        }
    }
    {
        // Switch on int with cases and default.
        so_int i = 2;
        if (i == (1)) {
            so_panic("unexpected i == 1");
        } else if (i == (2)) {
            so_println("%s", "i == 2");
        } else if (i == (3)) {
            so_panic("unexpected i == 3");
        } else {
            so_panic("unexpected default");
        }
    }
    {
        // Tagless switch (bool conditions).
        so_int x = 10;
        if (x > 100) {
            so_panic("unexpected x > 100");
        } else if (x > 0) {
            so_println("%s", "x > 0");
        } else {
            so_panic("unexpected default");
        }
    }
    {
        // Multiple values per case.
        so_int y = 3;
        if (y == (1) || y == (2) || y == (3)) {
            so_println("%s", "y == 3");
        } else if (y == (4) || y == (5) || y == (6)) {
            so_panic("unexpected y == 4, 5, 6");
        } else {
            so_panic("unexpected default");
        }
    }
    {
        // Switch with init statement.
        {
            so_int n = 42;
            if (n == (42)) {
                so_println("%s", "n == 42");
            } else {
                so_panic("unexpected default");
            }
        }
    }
    {
        // Switch on string.
        so_String s = so_str("hello");
        if (so_string_eq(s, so_str("hello"))) {
            so_println("%s", "s == hello");
        } else if (so_string_eq(s, so_str("bye"))) {
            so_panic("unexpected s == bye");
        } else {
            so_panic("unexpected default");
        }
    }
    {
        // Cases without default.
        so_int z = 5;
        if (z == (1)) {
            so_panic("unexpected z == 1");
        } else if (z == (5)) {
            so_println("%s", "z == 5");
        }
    }
    return 0;
}
