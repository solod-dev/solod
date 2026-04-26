#pragma once
#include "so/builtin/builtin.h"

// -- Types --

typedef struct main_Rect main_Rect;

// Rect is a rectangle.
typedef struct main_Rect {
    so_int W;
    so_int H;
} main_Rect;

// -- Functions and methods --

// Area returns the area of the rectangle.
//
static inline so_int main_Rect_Area(main_Rect r) {
    return r.W * r.H;
}

// Scale scales the rectangle by a factor.
void main_Rect_Scale(void* self, so_int factor);

static inline so_int add(so_int a, so_int b) {
    return a + b;
}
