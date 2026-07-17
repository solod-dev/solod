#include "main.h"

// -- Types --

typedef struct movie movie;

// A function pointer field with a reserved parameter name.
typedef struct movie {
    so_int (*rate)(so_int);
} movie;

// An interface method with a reserved parameter name.
typedef struct rater {
    void* self;
    so_int (*rate)(void* self, so_int);
} rater;

// -- Forward declarations --
static so_int scale(so_int long_, so_int register_);
static so_int shadow(so_int long_);

// -- Variables and constants --

// An exported identifier gets a package prefix,
// so it doesn't need mangling.
so_int main_NULL = 0;

// -- Implementation --

// C keywords used as parameter names.
static so_int scale(so_int long_, so_int register_) {
    so_int total = long_ * register_;
    return total;
}

// A mangled parameter (long -> long_) and a same-named local in a nested
// block are a legal C shadow, not a collision, so both are accepted.
static so_int shadow(so_int long_) {
    if (long_ > 0) {
        so_int long_ = 99;
        return long_;
    }
    return long_;
}

int main(void) {
    // C keywords used as local variables.
    so_int long_ = 10;
    so_int short_ = 20;
    so_int value = scale(long_, short_);
    (void)value;
    (void)shadow(value);
    // The name should be mangled everywhere it is used.
    for (so_int bool_ = 0; bool_ < long_; bool_++) {
        so_int b = bool_;
        (void)b;
    }
    // Reference the reserved-parameter types so they are emitted.
    movie m = {0};
    rater r = {0};
    (void)m;
    (void)r;
    return 0;
}
