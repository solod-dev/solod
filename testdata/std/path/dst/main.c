#include "main.h"

// -- Implementation --

int main(void) {
    {
        // Clean.
        so_String cleaned = path_Clean((mem_Allocator){0}, so_str("/opt/app/../config.json"));
        if (so_string_ne(cleaned, so_str("/opt/config.json"))) {
            so_panic(so_cstr(so_string_add(so_str("unexpected cleaned path: "), cleaned)));
        }
        mem_FreeString((mem_Allocator){0}, cleaned);
    }
    {
        // Split.
        so_R_str_str _res1 = path_Split(so_str("/opt/app/config.json"));
        so_String dir = _res1.val;
        so_String file = _res1.val2;
        if (so_string_ne(dir, so_str("/opt/app/"))) {
            so_panic(so_cstr(so_string_add(so_str("unexpected dir: "), dir)));
        }
        if (so_string_ne(file, so_str("config.json"))) {
            so_panic(so_cstr(so_string_add(so_str("unexpected file: "), file)));
        }
    }
    {
        // Join.
        so_String joined = path_Join((mem_Allocator){0}, (so_Slice){(so_String[3]){so_str("opt"), so_str("app"), so_str("config.json")}, 3, 3});
        if (so_string_ne(joined, so_str("opt/app/config.json"))) {
            so_panic(so_cstr(so_string_add(so_str("unexpected path: "), joined)));
        }
        mem_FreeString((mem_Allocator){0}, joined);
    }
    {
        // IsAbs.
        if (!path_IsAbs(so_str("/opt/app/config.json"))) {
            so_panic("want absolute");
        }
        if (path_IsAbs(so_str("opt/app/config.json"))) {
            so_panic("want not absolute");
        }
    }
    {
        // Dir.
        so_String dir = path_Dir((mem_Allocator){0}, so_str("/opt/app/config.json"));
        if (so_string_ne(dir, so_str("/opt/app"))) {
            so_panic(so_cstr(so_string_add(so_str("unexpected dir: "), dir)));
        }
        mem_FreeString((mem_Allocator){0}, dir);
    }
    {
        // Base.
        so_String base = path_Base(so_str("/opt/app/config.json"));
        if (so_string_ne(base, so_str("config.json"))) {
            so_panic(so_cstr(so_string_add(so_str("unexpected base: "), base)));
        }
    }
    {
        // Ext.
        so_String ext = path_Ext(so_str("/opt/app/config.json"));
        if (so_string_ne(ext, so_str(".json"))) {
            so_panic(so_cstr(so_string_add(so_str("unexpected ext: "), ext)));
        }
    }
    {
        // Match.
        so_R_bool_err _res2 = path_Match(so_str("/opt/*/*.js?n"), so_str("/opt/app/config.json"));
        bool ok = _res2.val;
        so_Error err = _res2.err;
        if (err.self != NULL) {
            so_panic(so_error_cstr(err));
        }
        if (!ok) {
            so_panic("want match");
        }
    }
    return 0;
}
