#include "main.h"

// -- Implementation --

int main(void) {
    {
        // Print.
        so_Result _res1 = fmt_Print("hello", "world");
        so_int n = _res1.val.as_int;
        so_Error err = _res1.err;
        if (err != NULL) {
            so_panic("Print failed");
        }
        if (n != 11) {
            so_panic("Print: wrong count");
        }
    }
    {
        // Println.
        so_Result _res2 = fmt_Println("hello", "world");
        so_int n = _res2.val.as_int;
        so_Error err = _res2.err;
        if (err != NULL) {
            so_panic("Println failed");
        }
        if (n != 12) {
            so_panic("Println: wrong count");
        }
    }
    {
        // Printf.
        so_String s = so_str("world");
        so_int d = 42;
        so_Result _res3 = fmt_Printf("s = %s, d = %d\n", so_cstr(s), d);
        so_int n = _res3.val.as_int;
        so_Error err = _res3.err;
        if (err != NULL) {
            so_panic("Printf failed");
        }
        if (n != 18) {
            so_panic("Printf: wrong count");
        }
    }
    {
        // Fprintf.
        strings_Builder sb = {0};
        so_String s = so_str("world");
        so_Result _res4 = fmt_Fprintf((io_Writer){.self = &sb, .Write = strings_Builder_Write}, "hello %s", so_cstr(s));
        so_int n = _res4.val.as_int;
        so_Error err = _res4.err;
        if (err != NULL) {
            so_panic("Fprintf failed");
        }
        if (n != 11) {
            so_panic("Fprintf: wrong count");
        }
        if (so_string_ne(strings_Builder_String(&sb), so_str("hello world"))) {
            so_panic("Fprintf: wrong output");
        }
        strings_Builder_Free(&sb);
    }
    {
        // Sscanf.
        int32_t a = 0;
        int32_t b = 0;
        so_Result _res5 = fmt_Sscanf("42 7", "%d %d", &a, &b);
        so_int n = _res5.val.as_int;
        so_Error err = _res5.err;
        if (err != NULL) {
            so_panic("Sscanf failed");
        }
        if (n != 2) {
            so_panic("Sscanf: wrong count");
        }
        if (a != 42 || b != 7) {
            so_panic("Sscanf: wrong values");
        }
    }
    {
        // Fscanf.
        strings_Reader r = strings_NewReader(so_str("100 200"));
        int32_t a = 0;
        int32_t b = 0;
        so_Result _res6 = fmt_Fscanf((io_Reader){.self = &r, .Read = strings_Reader_Read}, "%d %d", &a, &b);
        so_int n = _res6.val.as_int;
        so_Error err = _res6.err;
        if (err != NULL) {
            so_panic("Fscanf failed");
        }
        if (n != 2) {
            so_panic("Fscanf: wrong count");
        }
        if (a != 100 || b != 200) {
            so_panic("Fscanf: wrong values");
        }
    }
}
