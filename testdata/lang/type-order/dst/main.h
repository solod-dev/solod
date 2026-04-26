#pragma once
#include "so/builtin/builtin.h"

// -- Types --

typedef struct main_Node main_Node;
typedef struct main_Employee main_Employee;
typedef struct main_Pet main_Pet;

// Self-referencing struct type.
typedef struct main_Node {
    so_int value;
    main_Node* next;
} main_Node;

// Type referencing another type defined later.
typedef struct main_Employee {
    so_String name;
    main_Pet* pet;
} main_Employee;

typedef struct main_Pet {
    so_String name;
} main_Pet;
