//go:build ignore
#include "testing.h"

void testing_T_Errorf(void* self, const char* format, ...) {
    char buf[fmt_BufSize];
    va_list args;
    va_start(args, format);
    so_String msg = fmt_vsprintf((fmt_Buffer){buf, fmt_BufSize}, format, args);
    va_end(args);
    testing_T_Error(self, msg);
}

void testing_T_Fatalf(void* self, const char* format, ...) {
    char buf[fmt_BufSize];
    va_list args;
    va_start(args, format);
    so_String msg = fmt_vsprintf((fmt_Buffer){buf, fmt_BufSize}, format, args);
    va_end(args);
    testing_T_Fatal(self, msg);
}
