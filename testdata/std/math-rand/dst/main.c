#include "main.h"

// -- Implementation --

int main(void) {
    {
        // Int.
        rand_PCG pcg = rand_NewPCG(1, 2);
        rand_Rand r = rand_New((rand_Source){.self = &pcg, .Uint64 = rand_PCG_Uint64});
        so_int n1 = rand_Rand_Int(&r);
        if (n1 < 0) {
            so_panic("negative Int()");
        }
        so_int n2 = rand_Rand_Int(&r);
        if (n2 < 0) {
            so_panic("negative Int()");
        }
        if (n1 == n2) {
            so_panic("same Int() twice in a row");
        }
        so_println("%" PRIdINT " %" PRIdINT, n1, n2);
    }
    {
        // Float64.
        rand_PCG pcg = rand_NewPCG(1, 2);
        rand_Rand r = rand_New((rand_Source){.self = &pcg, .Uint64 = rand_PCG_Uint64});
        double f1 = rand_Rand_Float64(&r);
        if (f1 < 0 || f1 >= 1) {
            so_panic("Float64() out of range");
        }
        double f2 = rand_Rand_Float64(&r);
        if (f2 < 0 || f2 >= 1) {
            so_panic("Float64() out of range");
        }
        if (f1 == f2) {
            so_panic("same Float64() twice in a row");
        }
        so_println("%f %f", f1, f2);
    }
    {
        // Global functions.
        so_int n1 = rand_IntN(100);
        if (n1 < 0 || n1 >= 100) {
            so_panic("IntN() out of range");
        }
        so_int n2 = rand_IntN(100);
        if (n2 < 0 || n2 >= 100) {
            so_panic("IntN() out of range");
        }
        so_println("%" PRIdINT " %" PRIdINT, n1, n2);
    }
}
