#include "main.h"

// -- Implementation --

int main(void) {
    {
        // Return `const char*` from C.
        so_const_char* cstr = get_cstring("Hello, C!");
        so_String str = c_String(so_const_char, (cstr));
        if (so_string_ne(str, so_str("Hello, C!"))) {
            so_panic(so_cstr(so_string_add(so_str("unexpected string: "), str)));
        }
    }
    {
        // Use header included via so:include.c
        if (!isalpha(U'a')) {
            so_panic("isalpha failed");
        }
    }
    {
        // Typed C expression.
        double nan = NAN;
        if (nan == nan) {
            so_panic("nan == nan");
        }
        double x = sqrt(49);
        if (x != 7) {
            so_panic("x != 7");
        }
    }
    {
        // Raw C block.
        so_int b = 0;
        int a = 7;
        b = a * a;
        b = sqrt(b);
        if (b != 7) {
            so_panic("b != 7");
        }
    }
    return 0;
}
