#pragma once
#include "so/builtin/builtin.h"
#include "so/c/c.h"

// -- Embeds --

#ifdef so_build_hosted

#include <ctype.h>
#include <math.h>

#else

#define NAN (0.0 / 0.0)
static inline int isalpha(int ch) {
    return (ch >= 'A' && ch <= 'Z') || (ch >= 'a' && ch <= 'z');
}
static inline double sqrt(double x) {
    if (x < 0) return NAN;
    if (x == 0) return 0;
    double guess = x / 2;
    for (int i = 0; i < 10; i++) {
        guess = (guess + x / guess) / 2;
    }
    return guess;
}

#endif  // so_build_hosted

static inline const char* get_cstring(const char* s) {
    return s;
}
