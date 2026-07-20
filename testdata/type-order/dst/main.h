#pragma once
#include "so/builtin/builtin.h"

// -- Types --

typedef struct main_Node main_Node;
typedef struct main_Employee main_Employee;
typedef struct main_Pet main_Pet;
typedef struct main_Point main_Point;
typedef struct main_Rect main_Rect;
typedef struct main_Cell main_Cell;
typedef struct main_Grid main_Grid;
typedef struct main_Origin main_Origin;
typedef struct main_Outer main_Outer;
typedef struct main_Payload main_Payload;
typedef struct main_Reading main_Reading;

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

typedef struct main_Point {
    so_int X;
    so_int Y;
} main_Point;

// Type using a type defined later by value.
typedef struct main_Rect {
    main_Point Min;
    main_Point Max;
} main_Rect;

typedef struct main_Cell {
    so_int v;
} main_Cell;

// Array of a type defined later.
typedef struct main_Grid {
    main_Cell cells[4];
} main_Grid;

typedef struct main_Origin {
    so_int v;
} main_Origin;

// Named type of a type defined later.
typedef main_Origin main_Target;

typedef main_Payload (*main_Handler)(main_Payload);

// Func type held by value: Handler must precede Outer, but the struct types
// in its signature are fine as forward declarations.
typedef struct main_Outer {
    main_Handler handle;
} main_Outer;

typedef struct main_Payload {
    so_int v;
} main_Payload;
typedef so_int main_Meters;

// Pointer to a non-struct type: Meters has no forward declaration,
// so its definition must come first.
typedef struct main_Reading {
    main_Meters* depth;
} main_Reading;
