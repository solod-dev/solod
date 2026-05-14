#include "main.h"

// -- Implementation --

int main(void) {
    {
        // Writer -> Buffer -> Reader pipeline.
        bytes_Buffer buf = {0};
        buf = bytes_NewBuffer((mem_Allocator){0}, (so_Slice){0});
        bufio_Writer w = bufio_NewWriter((mem_Allocator){0}, (io_Writer){.self = &buf, .Write = bytes_Buffer_Write});
        bufio_Writer_WriteString(&w, so_str("Hello, "));
        bufio_Writer_WriteString(&w, so_str("World!"));
        bufio_Writer_WriteByte(&w, '\n');
        bufio_Writer_Flush(&w);
        bufio_Writer_Free(&w);
        strings_Reader sr = strings_NewReader(bytes_Buffer_String(&buf));
        bufio_Reader r = bufio_NewReader((mem_Allocator){0}, (io_Reader){.self = &sr, .Read = strings_Reader_Read});
        so_R_str_err _res1 = bufio_Reader_ReadString(&r, '\n');
        so_String line = _res1.val;
        so_Error err = _res1.err;
        if (err.self != NULL) {
            so_panic("ReadString failed");
        }
        if (so_string_ne(line, so_str("Hello, World!\n"))) {
            so_panic("unexpected line");
        }
        mem_FreeString((mem_Allocator){0}, line);
        bufio_Reader_Free(&r);
        bytes_Buffer_Free(&buf);
    }
    {
        // ReadByte and UnreadByte.
        strings_Reader sr = strings_NewReader(so_str("abc"));
        bufio_Reader r = bufio_NewReader((mem_Allocator){0}, (io_Reader){.self = &sr, .Read = strings_Reader_Read});
        so_R_byte_err _res2 = bufio_Reader_ReadByte(&r);
        so_byte b = _res2.val;
        so_Error err = _res2.err;
        if (err.self != NULL || b != 'a') {
            so_panic("ReadByte failed");
        }
        err = bufio_Reader_UnreadByte(&r);
        if (err.self != NULL) {
            so_panic("UnreadByte failed");
        }
        so_R_byte_err _res3 = bufio_Reader_ReadByte(&r);
        b = _res3.val;
        err = _res3.err;
        if (err.self != NULL || b != 'a') {
            so_panic("UnreadByte re-read failed");
        }
        bufio_Reader_Free(&r);
    }
    {
        // Peek.
        strings_Reader sr = strings_NewReader(so_str("hello"));
        bufio_Reader r = bufio_NewReader((mem_Allocator){0}, (io_Reader){.self = &sr, .Read = strings_Reader_Read});
        so_R_slice_err _res4 = bufio_Reader_Peek(&r, 3);
        so_Slice p = _res4.val;
        so_Error err = _res4.err;
        if (err.self != NULL || so_string_ne(so_bytes_string(p), so_str("hel"))) {
            so_panic("Peek failed");
        }
        bufio_Reader_Free(&r);
    }
    {
        // WriteRune.
        bytes_Buffer buf = {0};
        buf = bytes_NewBuffer((mem_Allocator){0}, (so_Slice){0});
        bufio_Writer w = bufio_NewWriter((mem_Allocator){0}, (io_Writer){.self = &buf, .Write = bytes_Buffer_Write});
        bufio_Writer_WriteRune(&w, U'A');
        bufio_Writer_Flush(&w);
        if (so_string_ne(bytes_Buffer_String(&buf), so_str("A"))) {
            so_panic("WriteRune failed");
        }
        bufio_Writer_Free(&w);
        bytes_Buffer_Free(&buf);
    }
    {
        // Scanner.
        strings_Reader sr = strings_NewReader(so_str("line1\nline2\n"));
        bufio_Scanner s = bufio_NewScanner((mem_Allocator){0}, (io_Reader){.self = &sr, .Read = strings_Reader_Read});
        so_int count = 0;
        for (; bufio_Scanner_Scan(&s);) {
            if (count == 0 && so_string_ne(bufio_Scanner_Text(&s), so_str("line1"))) {
                so_panic("Scanner: expected line1");
            }
            if (count == 1 && so_string_ne(bufio_Scanner_Text(&s), so_str("line2"))) {
                so_panic("Scanner: expected line2");
            }
            count++;
        }
        if (count != 2) {
            so_panic("Scanner: expected 2 lines");
        }
        bufio_Scanner_Free(&s);
    }
    return 0;
}
