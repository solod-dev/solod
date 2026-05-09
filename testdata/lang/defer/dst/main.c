#include "main.h"

// -- Variables and constants --
static so_int state = 0;

// -- Forward declarations --
static void xopen(so_int* x);
static void xclose(void* a);
static void funcScope(void);
static so_int funcWithReturn(void);
static void blockScope(void);

// -- Implementation --

static void xopen(so_int* x) {
    (*x)++;
}

static void xclose(void* a) {
    so_int* x = (so_int*)a;
    (*x)--;
}

static void funcScope(void) {
    xopen(&state);
    if (state != 1) {
        xclose(&state);
        so_panic("unexpected state");
    }
    xclose(&state);
}

static so_int funcWithReturn(void) {
    xopen(&state);
    if (state != 1) {
        xclose(&state);
        so_panic("unexpected state");
    }
    xclose(&state);
    return 42;
}

static void blockScope(void) {
    {
        xopen(&state);
        if (state != 1) {
            xclose(&state);
            so_panic("unexpected state");
        }
        xclose(&state);
    }
    if (state != 0) {
        so_panic("unexpected state");
    }
    {
        xopen(&state);
        if (state != 1) {
            xclose(&state);
            so_panic("unexpected state");
        }
        xclose(&state);
    }
    if (state != 0) {
        so_panic("unexpected state");
    }
}

int main(void) {
    funcScope();
    if (state != 0) {
        so_panic("unexpected state");
    }
    funcWithReturn();
    if (state != 0) {
        so_panic("unexpected state");
    }
    blockScope();
    if (state != 0) {
        so_panic("unexpected state");
    }
    return 0;
}
