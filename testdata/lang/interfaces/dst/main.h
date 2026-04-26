#pragma once
#include "so/builtin/builtin.h"

// -- Types --

typedef struct main_Canvas main_Canvas;
typedef struct main_Rect main_Rect;

typedef struct main_Shape {
    void* self;
    so_int (*Area)(void* self);
    so_int (*Perim)(void* self, so_int n);
} main_Shape;

typedef struct main_Canvas {
    so_String name;
    main_Shape shape;
} main_Canvas;

typedef struct main_Rect {
    so_int width;
    so_int height;
} main_Rect;

// -- Functions and methods --
so_int main_Rect_Area(void* self);
so_int main_Rect_Perim(void* self, so_int n);
