#include "main.h"

// -- Types --

typedef struct unode unode;
typedef struct uemployee uemployee;
typedef struct upet upet;
typedef struct upoint upoint;
typedef struct urect urect;

// Unexported self-referencing struct type.
typedef struct unode {
    so_int value;
    unode* next;
} unode;

// Unexported type referencing another type defined later.
typedef struct uemployee {
    so_String name;
    upet* pet;
} uemployee;

typedef struct upet {
    so_String name;
} upet;

typedef struct upoint {
    so_int x;
    so_int y;
} upoint;

// Unexported type using a type defined later by value.
typedef struct urect {
    upoint min;
    upoint max;
} urect;

// -- Implementation --

int main(void) {
    return 0;
}
