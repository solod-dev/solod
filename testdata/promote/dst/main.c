#include "main.h"

// -- Variables and constants --
so_int main_defaultCap = 8;

// -- Implementation --

// newCounter is called from an inline (header) NewStats,
// so it needs to be promoted to the header.
//
main_counter main_newCounter(void) {
    return (main_counter){.val = 0};
}

// inc is called from an inline (header) Stats.Inc method,
// so it needs to be promoted to the header.
//
void main_counter_inc(void* self) {
    main_counter* c = self;
    c->val++;
}

// GetCounter is a non-inline function that returns a promoted type.
main_counter main_GetCounter(void) {
    return main_newCounter();
}

int main(void) {
    main_Stats w = main_NewStats();
    main_Stats_Inc(&w);
    (void)(main_alias){.val = 1};
    (void)main_defaultCap;
    (void)main_version;
    (void)main_GetCounter();
    return 0;
}
