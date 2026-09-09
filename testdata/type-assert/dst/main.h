#pragma once
#include "so/builtin/builtin.h"

// -- Types --

typedef struct main_Rect main_Rect;
typedef struct main_Circle main_Circle;
typedef struct main_Canvas main_Canvas;

typedef struct main_Shape {
    void* self;
    so_int (*Area)(void* self);
} main_Shape;

static inline so_int main_Shape_Area(main_Shape self) {
    return self.Area(self.self);
}

typedef struct main_Rect {
    so_int width;
    so_int height;
} main_Rect;

typedef struct main_Circle {
    so_int radius;
} main_Circle;

typedef struct main_Canvas {
    main_Shape shape;
} main_Canvas;

// -- Functions and methods --
so_int main_Rect_Area(void* self);
so_int main_Circle_Area(void* self);
