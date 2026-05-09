#include "main.h"

// -- Forward declarations --
static void logger(void);
static void defaults(void);

// -- Implementation --

int main(int argc, char* argv[]) {
    so_String _so_argv[argc];
    so_args_init(argc, argv, _so_argv);
    logger();
    defaults();
    return 0;
}

static void logger(void) {
    // Logger writing to stdout.
    slog_TextHandler h = slog_NewTextHandler((io_Writer){.self = os_Stdout, .Write = os_File_Write}, slog_LevelInfo);
    slog_Logger l = slog_New((slog_Handler){.self = &h, .Enabled = slog_TextHandler_Enabled, .Handle = slog_TextHandler_Handle});
    // Enabled check.
    if (slog_Logger_Enabled(&l, slog_LevelDebug)) {
        so_panic("debug should not be enabled");
    }
    if (!slog_Logger_Enabled(&l, slog_LevelInfo)) {
        so_panic("info should be enabled");
    }
    // Log at info - should appear.
    slog_Logger_Info(&l, so_str("hello world"), (so_Slice){(slog_Attr[2]){slog_String(so_str("user"), so_str("john")), slog_Int(so_str("count"), 42)}, 2, 2});
    // Log at debug - should be filtered.
    slog_Logger_Debug(&l, so_str("hidden"), (so_Slice){(slog_Attr[0]){}, 0, 0});
    // Log with no attrs.
    slog_Logger_Warn(&l, so_str("caution"), (so_Slice){(slog_Attr[0]){}, 0, 0});
    // Log with float and bool attrs.
    slog_Logger_Error(&l, so_str("failure"), (so_Slice){(slog_Attr[2]){slog_Float64(so_str("elapsed"), 1.5), slog_Bool(so_str("retry"), true)}, 2, 2});
    // Log with string that needs quoting.
    slog_Logger_Info(&l, so_str("test quoting"), (so_Slice){(slog_Attr[1]){slog_String(so_str("msg"), so_str("hello world"))}, 1, 1});
}

static void defaults(void) {
    // Default logger should be usable.
    slog_Info(so_str("default test"), (so_Slice){(slog_Attr[1]){slog_Int(so_str("port"), 8080)}, 1, 1});
}
