#pragma once
#include "so/builtin/builtin.h"

// -- Types --

typedef struct main_Sum main_Sum;

typedef struct main_Sum {
    so_int v;
} main_Sum;

// -- Functions and methods --
void main_Sum_Add(void* self, so_Slice nums);
