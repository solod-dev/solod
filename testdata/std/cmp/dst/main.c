#include "main.h"

// -- Implementation --

int main(void) {
    {
        // Compare numbers.
        so_int a = 11;
        so_int b = 22;
        if (cmp_Compare(so_int, (a), (b)) >= 0) {
            so_panic("Compare failed");
        }
        if (cmp_Compare(so_int, (a), (a)) != 0) {
            so_panic("Compare failed");
        }
    }
    {
        // Compare strings.
        so_String a = so_str("hello");
        so_String b = so_str("world");
        if (cmp_Compare(so_String, (a), (b)) >= 0) {
            so_panic("Compare failed");
        }
        if (cmp_Compare(so_String, (a), (a)) != 0) {
            so_panic("Compare failed");
        }
    }
    {
        // Equal numbers.
        so_int a = 11;
        so_int b = 22;
        if (cmp_Equal(so_int, (a), (b))) {
            so_panic("Equal failed");
        }
        if (!cmp_Equal(so_int, (a), (a))) {
            so_panic("Equal failed");
        }
    }
    {
        // Equal strings.
        so_String a = so_str("hello");
        so_String b = so_str("world");
        if (cmp_Equal(so_String, (a), (b))) {
            so_panic("Equal failed");
        }
        if (!cmp_Equal(so_String, (a), (a))) {
            so_panic("Equal failed");
        }
    }
    {
        // Less numbers.
        so_int a = 11;
        so_int b = 22;
        if (!cmp_Less(so_int, (a), (b))) {
            so_panic("Less failed");
        }
        if (cmp_Less(so_int, (b), (a))) {
            so_panic("Less failed");
        }
    }
    {
        // Less strings.
        so_String a = so_str("hello");
        so_String b = so_str("world");
        if (!cmp_Less(so_String, (a), (b))) {
            so_panic("Less failed");
        }
        if (cmp_Less(so_String, (b), (a))) {
            so_panic("Less failed");
        }
    }
}
