#include "person.ext.h"

int64_t account_inc_balance(Account* a, int64_t amount) {
    int64_t balBefore = a->balance;
    so_byte* flags = a->flags.ptr;
    printf("name = %s balance = %" PRId64 " flags[0] = %u\n",
           a->name.ptr, balBefore, a->flags.len > 0 ? flags[0] : 0);
    a->balance += amount;
    return balBefore;
}

void account_set_name(Account* a, so_String name) {
    a->name = name;
}

void write_acc(Account* a, const char* fmt, ...) {
    va_list args;
    va_start(args, fmt);
    printf("Account %s: ", a->name.ptr);
    vprintf(fmt, args);
    printf("\n");
    va_end(args);
}
