#include "main.h"

// -- Implementation --

so_int main_Rect_Area(void* self) {
    main_Rect* r = self;
    return r->width * r->height;
}

so_int main_Circle_Area(void* self) {
    main_Circle* c = self;
    return 3 * c->radius * c->radius;
}

int main(void) {
    main_Rect r = (main_Rect){.width = 10, .height = 5};
    main_Circle c = (main_Circle){.radius = 2};
    {
        // Both targets on a matching type.
        main_Shape s = (main_Shape){.self = &r, .Area = main_Rect_Area};
        bool ok = (s.Area == main_Rect_Area);
        main_Rect* p = ok ? (main_Rect*)s.self : NULL;
        if (!ok) {
            so_panic("want ok");
        }
        if (main_Rect_Area(p) != 50) {
            so_panic("p.Area() != 50");
        }
    }
    {
        // A failed assertion gives a nil pointer.
        main_Shape s = (main_Shape){.self = &c, .Area = main_Circle_Area};
        bool ok = (s.Area == main_Rect_Area);
        main_Rect* p = ok ? (main_Rect*)s.self : NULL;
        if (ok) {
            so_panic("want !ok");
        }
        if (p != NULL) {
            so_panic("want nil p");
        }
    }
    {
        // Blank ok target.
        main_Shape s = (main_Shape){.self = &r, .Area = main_Rect_Area};
        main_Rect* p = (s.Area == main_Rect_Area) ? (main_Rect*)s.self : NULL;
        if (main_Rect_Area(p) != 50) {
            so_panic("p.Area() != 50");
        }
        main_Shape s2 = (main_Shape){.self = &c, .Area = main_Circle_Area};
        main_Rect* q = (s2.Area == main_Rect_Area) ? (main_Rect*)s2.self : NULL;
        if (q != NULL) {
            so_panic("want nil q");
        }
    }
    {
        // Blank value target.
        main_Shape s = (main_Shape){.self = &r, .Area = main_Rect_Area};
        bool ok = (s.Area == main_Rect_Area);
        if (!ok) {
            so_panic("want ok");
        }
    }
    {
        // Plain assignment.
        main_Shape s = (main_Shape){.self = &r, .Area = main_Rect_Area};
        main_Rect* p = NULL;
        bool ok = false;
        ok = (s.Area == main_Rect_Area);
        p = ok ? (main_Rect*)s.self : NULL;
        if (!ok || main_Rect_Area(p) != 50) {
            so_panic("want ok and p.Area() == 50");
        }
        s = (main_Shape){.self = &c, .Area = main_Circle_Area};
        ok = (s.Area == main_Rect_Area);
        p = ok ? (main_Rect*)s.self : NULL;
        if (ok || p != NULL) {
            so_panic("want !ok and nil p");
        }
    }
    {
        // A selector operand is read two times.
        main_Canvas canvas = (main_Canvas){.shape = (main_Shape){.self = &r, .Area = main_Rect_Area}};
        bool ok = (canvas.shape.Area == main_Rect_Area);
        main_Rect* p = ok ? (main_Rect*)canvas.shape.self : NULL;
        if (!ok || main_Rect_Area(p) != 50) {
            so_panic("want ok and p.Area() == 50");
        }
    }
    {
        // A redeclared target keeps its type.
        main_Shape s = (main_Shape){.self = &r, .Area = main_Rect_Area};
        bool ok = false;
        ok = (s.Area == main_Rect_Area);
        main_Rect* p = ok ? (main_Rect*)s.self : NULL;
        if (!ok || main_Rect_Area(p) != 50) {
            so_panic("want ok and p.Area() == 50");
        }
    }
    {
        // A direct assertion casts without a check.
        main_Shape s = (main_Shape){.self = &r, .Area = main_Rect_Area};
        main_Rect* p = (main_Rect*)s.self;
        if (main_Rect_Area(p) != 50) {
            so_panic("p.Area() != 50");
        }
    }
    {
        // A blank value target with a plain assignment.
        main_Shape s = (main_Shape){.self = &r, .Area = main_Rect_Area};
        bool ok = false;
        ok = (s.Area == main_Rect_Area);
        if (!ok) {
            so_panic("want ok");
        }
    }
    {
        // An assertion in the init of an if statement.
        main_Shape s = (main_Shape){.self = &c, .Area = main_Circle_Area};
        {
            bool ok = (s.Area == main_Rect_Area);
            if (ok) {
                so_panic("want !ok");
            }
        }
    }
    {
        // A nil interface matches no concrete type.
        main_Shape s = {};
        bool ok = (s.Area == main_Rect_Area);
        main_Rect* p = ok ? (main_Rect*)s.self : NULL;
        if (ok || p != NULL) {
            so_panic("want !ok and nil p");
        }
    }
    return 0;
}
