#include "main.h"

// -- Implementation --

int main(void) {
    {
        // Is.
        if (!unicode_IsDigit(U'0')) {
            so_panic("IsDigit failed");
        }
        if (!unicode_IsLetter(U'a')) {
            so_panic("IsLetter failed");
        }
        if (!unicode_IsLower(U'a')) {
            so_panic("IsLower failed");
        }
        if (!unicode_IsSpace(U' ')) {
            so_panic("IsSpace failed");
        }
        if (!unicode_IsTitle(U'ᾭ')) {
            so_panic("IsTitle failed");
        }
        if (!unicode_IsUpper(U'A')) {
            so_panic("IsUpper failed");
        }
    }
    {
        // To.
        if (unicode_ToLower(U'A') != U'a') {
            so_panic("ToLower failed");
        }
        if (unicode_ToTitle(U'a') != U'A') {
            so_panic("ToTitle failed");
        }
        if (unicode_ToUpper(U'a') != U'A') {
            so_panic("ToUpper failed");
        }
        if (unicode_To(unicode_UpperCase, U'a') != U'A') {
            so_panic("To failed");
        }
    }
    return 0;
}
