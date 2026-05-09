#include "main.h"

// -- Implementation --

int main(void) {
    so_println("%.*s %.*s %.*s", runtime_Version().len, runtime_Version().ptr, runtime_GOOS.len, runtime_GOOS.ptr, runtime_GOARCH.len, runtime_GOARCH.ptr);
    {
        // Version.
        so_String v = runtime_Version();
        if (so_len(v) == 0) {
            so_panic("Empty version");
        }
    }
    {
        // GOOS.
        so_String os = runtime_GOOS;
        if (so_string_ne(os, so_str("darwin")) && so_string_ne(os, so_str("linux")) && so_string_ne(os, so_str("windows")) && so_string_ne(os, so_str("wasip1"))) {
            so_panic("Unexpected GOOS");
        }
    }
    {
        // GOARCH.
        so_String arch = runtime_GOARCH;
        if (so_string_ne(arch, so_str("amd64")) && so_string_ne(arch, so_str("arm64")) && so_string_ne(arch, so_str("386")) && so_string_ne(arch, so_str("riscv64")) && so_string_ne(arch, so_str("wasm"))) {
            so_panic("Unexpected GOARCH");
        }
    }
    return 0;
}
