#include "main.h"

// -- Variables and constants --
static so_int state = 0;

// -- Implementation --

int main(void) {
    if (state != 42) {
        so_panic("init() did not run");
    }
    so_println("%s", "ok");
}

static void __attribute__((constructor)) main_init() {
    state = 42;
}
