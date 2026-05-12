#include <stdarg.h>
#include "builtin.h"

#ifdef so_build_hosted
#include <stdio.h>
#endif

// Command-line arguments, populated by main().
so_Slice os_Args = {0};

// utf8_decode decodes one UTF-8 rune from string s at byte offset i.
// Stores the byte width in *w.
// Returns the decoded rune, or 0xFFFD for invalid UTF-8.
so_rune so_utf8_decode(so_String s, so_int i, so_int* w) {
    const uint8_t* p = (const uint8_t*)s.ptr + i;
    so_int remaining = s.len - i;
    if (remaining <= 0) {
        *w = 0;
        return 0xFFFD;
    }

    uint8_t b = p[0];
    if (b < 0x80) {
        *w = 1;
        return (so_rune)b;
    }
    if ((b & 0xE0) == 0xC0 && remaining >= 2 &&
        (p[1] & 0xC0) == 0x80) {
        *w = 2;
        return ((so_rune)(b & 0x1F) << 6) |
               ((so_rune)(p[1] & 0x3F));
    }
    if ((b & 0xF0) == 0xE0 && remaining >= 3 &&
        (p[1] & 0xC0) == 0x80 && (p[2] & 0xC0) == 0x80) {
        *w = 3;
        return ((so_rune)(b & 0x0F) << 12) |
               ((so_rune)(p[1] & 0x3F) << 6) |
               ((so_rune)(p[2] & 0x3F));
    }
    if ((b & 0xF8) == 0xF0 && remaining >= 4 &&
        (p[1] & 0xC0) == 0x80 &&
        (p[2] & 0xC0) == 0x80 &&
        (p[3] & 0xC0) == 0x80) {
        *w = 4;
        return ((so_rune)(b & 0x07) << 18) |
               ((so_rune)(p[1] & 0x3F) << 12) |
               ((so_rune)(p[2] & 0x3F) << 6) |
               ((so_rune)(p[3] & 0x3F));
    }

    *w = 1;
    return 0xFFFD;
}

// string_runes_impl decodes UTF-8 string bytes into a rune buffer.
so_Slice so_string_runes_impl(so_String s, so_rune* buf) {
    so_int n = 0;
    for (so_int i = 0; i < s.len;) {
        so_int w = 0;
        buf[n++] = so_utf8_decode(s, i, &w);
        i += w;
    }
    return (so_Slice){buf, n, n};
}

// utf8_encode encodes a single rune into buf (up to 4 bytes).
// Returns the number of bytes written.
so_int so_utf8_encode(so_rune r, char* buf) {
    if (r < 0x80) {
        buf[0] = (char)r;
        return 1;
    }
    if (r < 0x800) {
        buf[0] = (char)(0xC0 | (r >> 6));
        buf[1] = (char)(0x80 | (r & 0x3F));
        return 2;
    }
    if (r < 0x10000) {
        buf[0] = (char)(0xE0 | (r >> 12));
        buf[1] = (char)(0x80 | ((r >> 6) & 0x3F));
        buf[2] = (char)(0x80 | (r & 0x3F));
        return 3;
    }
    buf[0] = (char)(0xF0 | (r >> 18));
    buf[1] = (char)(0x80 | ((r >> 12) & 0x3F));
    buf[2] = (char)(0x80 | ((r >> 6) & 0x3F));
    buf[3] = (char)(0x80 | (r & 0x3F));
    return 4;
}

// runes_string_impl encodes runes into a UTF-8 buffer and returns a string.
so_String so_runes_string_impl(so_Slice rs, char* buf) {
    so_int pos = 0;
    so_rune* runes = (so_rune*)rs.ptr;
    for (so_int i = 0; i < rs.len; i++) {
        pos += so_utf8_encode(runes[i], buf + pos);
    }
    return (so_String){buf, pos};
}

#ifdef so_build_hosted

// print writes the formatted string to stdout.
// Returns the number of bytes written.
int so_print(const char* format, ...) {
    va_list args;
    va_start(args, format);
    int n = vprintf(format, args);
    va_end(args);
    return n;
}

// println writes the formatted string to stdout with a newline.
// Returns the number of bytes written.
int so_println(const char* format, ...) {
    va_list args;
    va_start(args, format);
    int n = vprintf(format, args);
    va_end(args);
    putchar('\n');
    return n + 1;
}

#else

int so_print(const char* format, ...) {
    (void)format;
    return 0;
}
int so_println(const char* format, ...) {
    (void)format;
    return 0;
}

#endif  // so_build_hosted
