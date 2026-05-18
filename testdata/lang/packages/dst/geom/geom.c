#include "geom.h"

// -- Forward declarations --
static double rectArea(double width, double height);

// -- Implementation --

static double rectArea(double width, double height) {
    return width * height;
}

double geom_RectArea(double width, double height) {
    return rectArea(width, height);
}
