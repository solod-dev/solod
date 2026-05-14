#include "main.h"

// -- Implementation --

int main(void) {
    {
        // Read.
        so_Slice buf = so_make_slice(so_byte, 16, 16);
        so_R_int_err _res1 = crand_Read(buf);
        so_int n = _res1.val;
        so_Error err = _res1.err;
        if (err.self != NULL) {
            so_panic("failed to read random data");
        }
        if (n != so_len(buf)) {
            so_panic("short read of random data");
        }
    }
    {
        // Read empty slice.
        so_Slice buf = so_make_slice(so_byte, 0, 0);
        so_R_int_err _res2 = crand_Read(buf);
        so_int n = _res2.val;
        so_Error err = _res2.err;
        if (err.self != NULL) {
            so_panic("failed to read random data");
        }
        if (n != 0) {
            so_panic("non-zero read of empty slice");
        }
    }
    {
        // Reader.
        so_Slice buf = so_make_slice(so_byte, 16, 16);
        so_R_int_err _res3 = crand_Reader.Read(crand_Reader.self, buf);
        so_int n = _res3.val;
        so_Error err = _res3.err;
        if (err.self != NULL) {
            so_panic("failed to read random data");
        }
        if (n != so_len(buf)) {
            so_panic("short read of random data");
        }
    }
    {
        // Text.
        so_Slice buf = so_make_slice(so_byte, 26, 26);
        so_String s = crand_Text(buf);
        if (so_len(s) != 26) {
            so_panic("unexpected length of random text");
        }
        so_println("%.*s", s.len, s.ptr);
    }
    return 0;
}
