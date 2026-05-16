#pragma once
#include "so/builtin/builtin.h"

// -- Types --

typedef struct sub_Point sub_Point;

typedef struct sub_Point {
    so_int X;
    so_int Y;
} sub_Point;

// -- Result types --

typedef struct sub_PointResult {
    sub_Point val;
    so_Error err;
} sub_PointResult;

// -- Functions and methods --
sub_PointResult sub_MakePoint(so_int x, so_int y);
