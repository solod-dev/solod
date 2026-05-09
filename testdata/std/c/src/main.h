#if __STDC_HOSTED__

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

#endif // __STDC_HOSTED__

static inline const char* get_cstring(const char* s) {
    return s;
}
