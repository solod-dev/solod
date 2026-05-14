#include "main.h"

// -- Implementation --

int main(void) {
    flag_FlagSet flags = flag_NewFlagSet(so_str("example"), flag_ContinueOnError);
    bool b = false;
    flag_FlagSet_BoolVar(&flags, &b, so_str("b"), false, so_str("a boolean flag"));
    so_int n = 0;
    flag_FlagSet_IntVar(&flags, &n, so_str("n"), 0, so_str("an int flag"));
    double f = 0;
    flag_FlagSet_Float64Var(&flags, &f, so_str("f"), 0.0, so_str("a float flag"));
    so_String s = so_str("");
    flag_FlagSet_StringVar(&flags, &s, so_str("s"), so_str("default"), so_str("a string flag"));
    so_Error err = flag_FlagSet_Parse(&flags, (so_Slice){(so_String[7]){so_str("-b"), so_str("-n"), so_str("42"), so_str("-f"), so_str("3.14"), so_str("-s"), so_str("hello")}, 7, 7});
    if (err.self != NULL) {
        so_panic(so_error_cstr(err));
    }
    if (!b) {
        so_panic("b != true");
    }
    if (n != 42) {
        so_panic("n != 42");
    }
    if (f != 3.14) {
        so_panic("f != 3.14");
    }
    if (so_string_ne(s, so_str("hello"))) {
        so_panic("s != hello");
    }
    return 0;
}
