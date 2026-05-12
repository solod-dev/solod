#include "main.h"

// -- Implementation --

// Scale scales the rectangle by a factor.
void main_Rect_Scale(void* self, so_int factor) {
    main_Rect* r = self;
    r->W *= factor;
    r->H *= factor;
}

int main(void) {
    main_Rect r = (main_Rect){.W = 3, .H = 4};
    (void)main_Rect_Area(r);
    main_Rect_Scale(&r, 2);
    (void)add(1, 2);
    return 0;
}
