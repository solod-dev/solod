#include "main.h"

// -- Implementation --

int main(void) {
    {
        // Encode.
        so_Slice src = so_string_bytes(so_str("Hello Gopher!"));
        so_Slice dst = slices_Make(so_byte, (mem_System), (hex_EncodedLen(so_len(src))));
        hex_Encode(dst, src);
        if (so_string_ne(so_bytes_string(dst), so_str("48656c6c6f20476f7068657221"))) {
            so_panic("unexpected Encode result");
        }
        mem_FreeSlice(so_byte, (mem_System), (dst));
    }
    {
        // EncodeToString.
        so_Slice src = so_string_bytes(so_str("Hello Gopher!"));
        so_String encoded = hex_EncodeToString(mem_System, src);
        if (so_string_ne(encoded, so_str("48656c6c6f20476f7068657221"))) {
            so_panic("unexpected EncodeToString result");
        }
        mem_FreeString(mem_System, encoded);
    }
    {
        // Decode.
        so_Slice src = so_string_bytes(so_str("48656c6c6f20476f7068657221"));
        so_Slice dst = slices_Make(so_byte, (mem_System), (hex_DecodedLen(so_len(src))));
        so_R_int_err _res1 = hex_Decode(dst, src);
        so_int n = _res1.val;
        so_Error err = _res1.err;
        if (err.self != NULL) {
            so_panic(so_error_cstr(err));
        }
        if (so_string_ne(so_bytes_string(so_slice(so_byte, dst, 0, n)), so_str("Hello Gopher!"))) {
            so_panic("unexpected Decode result");
        }
        mem_FreeSlice(so_byte, (mem_System), (dst));
    }
    {
        // DecodeString.
        const so_String s = so_str("48656c6c6f20476f7068657221");
        so_R_slice_err _res2 = hex_DecodeString(mem_System, s);
        so_Slice decoded = _res2.val;
        so_Error err = _res2.err;
        if (err.self != NULL) {
            so_panic(so_error_cstr(err));
        }
        if (so_string_ne(so_bytes_string(decoded), so_str("Hello Gopher!"))) {
            so_panic("unexpected DecodeString result");
        }
        mem_FreeSlice(so_byte, (mem_System), (decoded));
    }
    {
        // Dump.
        so_Slice content = so_string_bytes(so_str("Go is an open source programming language."));
        so_String dmp = hex_Dump(mem_System, content);
        fmt_Printf("%s", so_cstr(dmp));
        mem_FreeString(mem_System, dmp);
    }
    return 0;
}
