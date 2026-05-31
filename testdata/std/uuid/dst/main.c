#include "main.h"

// -- Implementation --

int main(void) {
    const so_String ustr = so_str("f81d4fae-7dec-11d0-a765-00a0c91e6bf6");
    {
        // NewV4 and NewV7.
        uuid_UUID u4 = uuid_NewV4();
        if (uuid_UUID_Version(u4) != 4) {
            so_panic("NewV4() version != 4");
        }
        uuid_UUID u7 = uuid_NewV7();
        if (uuid_UUID_Version(u7) != 7) {
            so_panic("NewV7() version != 7");
        }
    }
    {
        // String and Parse.
        uuid_UUID u1 = uuid_MustParse(ustr);
        so_Slice buf = so_make_slice(so_byte, uuid_UUIDLen, uuid_UUIDLen);
        so_String s = uuid_UUID_String(u1, buf);
        if (so_string_ne(s, ustr)) {
            so_panic("String() mismatch");
        }
        uuid_UUIDResult _res1 = uuid_Parse(s);
        uuid_UUID u2 = _res1.val;
        so_Error err = _res1.err;
        if (err.self != NULL) {
            so_panic(so_error_cstr(err));
        }
        if (so_mem_ne(&u1, &u2, sizeof(uuid_UUID))) {
            so_panic("Parse/String mismatch");
        }
    }
    {
        // Compare.
        uuid_UUID unil = uuid_Nil();
        uuid_UUID uid = uuid_MustParse(ustr);
        uuid_UUID umax = uuid_Max();
        if (uuid_UUID_Compare(uid, unil) <= 0) {
            so_panic("Compare: uid <= unil");
        }
        if (uuid_UUID_Compare(uid, umax) >= 0) {
            so_panic("Compare: uid >= umax");
        }
        if (uuid_UUID_Compare(uid, uid) != 0) {
            so_panic("Compare: uid != uid");
        }
    }
    return 0;
}
