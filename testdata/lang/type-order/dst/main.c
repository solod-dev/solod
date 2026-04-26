#include "main.h"

// -- Types --

typedef struct unode unode;
typedef struct uemployee uemployee;
typedef struct upet upet;

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

// -- Implementation --

int main(void) {
}
