#pragma once
#include "so/builtin/builtin.h"

// -- Embeds --

typedef struct {
    int val;
} main_Box;

// -- Functions and methods --

#define identity(T, val_) ({ \
    val_; \
})

#define setPtr(T, ptr_, val_) do { \
    *ptr_ = val_; \
} while (0)

#define increment(T, n_) ({ \
    T _n = n_; \
    _n = _n + 1; \
    _n = _n + 1; \
    _n; \
})

#define a(T, n_) ({ \
    so_int _some = 11; \
    (void)_some; \
    T _x = b(T, (n_)) + 1; \
    _x; \
})

#define b(T, n_) ({ \
    double _some = 22.2; \
    (void)_some; \
    T _x = c(T, (n_)) + 1; \
    _x; \
})

#define c(T, n_) ({ \
    so_String _some = so_str("33"); \
    (void)_some; \
    T _x = n_ + 1; \
    _x; \
})

#define work(T, v_) ({ \
    (so_R_ptr_err){.val = v_, .err = (so_Error){0}}; \
})

#define main_Box_set(T, b_, val_) do { \
    b_->val = val_; \
} while (0)
