#include "main.h"

// -- Implementation --

int main(void) {
    {
        // IndexRabinKarp.
        so_Slice b = so_string_bytes(so_str("go is fun"));
        so_int idx = bytealg_IndexRabinKarp(b, so_string_bytes(so_str("is")));
        if (idx != 3) {
            so_panic("IndexRabinKarp failed");
        }
    }
    {
        // LastIndexRabinKarp.
        so_Slice b = so_string_bytes(so_str("hello"));
        so_int idx = bytealg_LastIndexRabinKarp(b, so_string_bytes(so_str("l")));
        if (idx != 3) {
            so_panic("LastIndexRabinKarp failed");
        }
    }
    {
        // Compare.
        so_Slice b = so_string_bytes(so_str("abc"));
        if (bytealg_Compare(b, so_string_bytes(so_str("abb"))) <= 0) {
            so_panic("Compare failed");
        }
        if (bytealg_Compare(b, so_string_bytes(so_str("abd"))) >= 0) {
            so_panic("Compare failed");
        }
        if (bytealg_Compare(b, so_string_bytes(so_str("abc"))) != 0) {
            so_panic("Compare failed");
        }
    }
    {
        // Count and CountString.
        so_Slice b = so_string_bytes(so_str("hello world"));
        so_int n = bytealg_Count(b, 'o');
        if (n != 2) {
            so_panic("Count failed");
        }
        so_String s = so_str("hello world");
        n = bytealg_CountString(s, 'o');
        if (n != 2) {
            so_panic("CountString failed");
        }
    }
    {
        // Equal.
        so_Slice a = so_string_bytes(so_str("hello"));
        so_Slice b = so_string_bytes(so_str("hello"));
        if (!bytealg_Equal(a, b)) {
            so_panic("Equal failed");
        }
        so_Slice c = so_string_bytes(so_str("world"));
        if (bytealg_Equal(a, c)) {
            so_panic("Equal failed");
        }
    }
    {
        // IndexByte and IndexByteString.
        so_Slice b = so_string_bytes(so_str("hello"));
        so_int idx = bytealg_IndexByte(b, 'l');
        if (idx != 2) {
            so_panic("IndexByte failed");
        }
        so_String s = so_str("hello");
        idx = bytealg_IndexByteString(s, 'l');
        if (idx != 2) {
            so_panic("IndexByteString failed");
        }
    }
    {
        // LastIndexByte and LastIndexByteString.
        so_Slice b = so_string_bytes(so_str("hello"));
        so_int idx = bytealg_LastIndexByte(b, 'l');
        if (idx != 3) {
            so_panic("LastIndexByte failed");
        }
        so_String s = so_str("hello");
        idx = bytealg_LastIndexByteString(s, 'l');
        if (idx != 3) {
            so_panic("LastIndexByteString failed");
        }
    }
    return 0;
}
