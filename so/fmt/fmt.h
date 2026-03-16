#include <stdarg.h>
#include <stdio.h>
#include "so/builtin/builtin.h"
#include "so/io/io.h"

// BufSize is the size of the internal formatting buffer in bytes.
#define fmt_BufSize 1024

extern so_Error fmt_ErrPrint;  // print failure
extern so_Error fmt_ErrScan;   // scan failure
extern so_Error fmt_ErrSize;   // buffer size exceeded

// Print writes its arguments to standard output, separated by spaces.
// It returns the number of bytes written and any write error encountered.
#define fmt_Print(...) fmt_print(false, __VA_ARGS__, NULL)
// Println is like Print but adds a newline at the end.
#define fmt_Println(...) fmt_print(true, __VA_ARGS__, NULL)
so_Result fmt_print(int newline, ...);

// Printf formats according to a format specifier and writes to standard output.
// It returns the number of bytes written and any write error encountered.
so_Result fmt_Printf(const char* format, ...);

// Fprintf formats according to a format specifier and writes to w.
// It returns the number of bytes written and any write error encountered.
// Returns [ErrSize] if the output size exceeds BufSize.
so_Result fmt_Fprintf(io_Writer w, const char* format, ...);

// Scanf scans text read from standard input, storing successive
// space-separated values into successive arguments as determined by the format.
// It returns the number of items successfully scanned.
so_Result fmt_Scanf(const char* format, ...);

// Sscanf scans the argument string, storing successive space-separated
// values into successive arguments as determined by the format.
// It returns the number of items successfully scanned.
so_Result fmt_Sscanf(const char* str, const char* format, ...);

// Fscanf scans text read from r, storing successive space-separated
// values into successive arguments as determined by the format.
// It returns the number of items successfully scanned.
so_Result fmt_Fscanf(io_Reader r, const char* format, ...);
