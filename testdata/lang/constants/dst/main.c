#include "main.h"

// -- Variables and constants --

// File-level constants.
static const so_int fInt = 42;
static const so_String fString = so_str("file");
static const main_HttpStatus statusSecret = 999;
main_Point main_PointZero = (main_Point){.X = main_Zero, .Y = main_Zero};
main_Point main_PointSubZero = (main_Point){.X = sub_Zero, .Y = sub_Zero};

// -- Implementation --

int main(void) {
    {
        // Local constants.
        const int64_t lInt = 500000000;
        (void)lInt;
        const double lFloat = 3e20 / lInt;
        (void)lFloat;
        const so_String lString = so_str("local");
        (void)lString;
    }
    {
        // Using constants in expressions.
        main_HttpStatus status = main_StatusOK;
        (void)(status != main_StatusNotFound);
        main_HttpStatus secret = statusSecret;
        (void)(secret > main_StatusOK);
        main_ServerState state = main_StateConnected;
        (void)so_string_eq(state, main_StateIdle);
    }
    {
        // Using iota constants.
        main_Day day = main_Monday;
        (void)(day == main_Sunday);
    }
    {
        // Using _ on file level is not supported,
        // so silence the unused file-level constants here.
        (void)fInt;
        (void)fString;
    }
    return 0;
}
