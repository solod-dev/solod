#include "main.h"

// -- Implementation --

int main(void) {
    {
        // Big endian.
        const uint64_t n1 = 0x0123456789abcdef;
        const uint64_t n2 = 0xfedcba9876543210;
        so_Slice buf = so_make_slice(so_byte, 8, 8);
        binary_BE_PutUint64(binary_BigEndian, buf, n1);
        {
            uint64_t got = binary_BE_Uint64(binary_BigEndian, buf);
            if (got != n1) {
                so_panic("BigEndian: invalid decoded n1");
            }
        }
        buf = binary_BE_AppendUint64(binary_BigEndian, so_slice(so_byte, buf, 0, 0), n2);
        {
            uint64_t got = binary_BE_Uint64(binary_BigEndian, buf);
            if (got != n2) {
                so_panic("BigEndian: invalid decoded n2");
            }
        }
    }
    {
        // Little endian.
        const uint64_t n1 = 0x0123456789abcdef;
        const uint64_t n2 = 0xfedcba9876543210;
        so_Slice buf = so_make_slice(so_byte, 8, 8);
        binary_LE_PutUint64(binary_LittleEndian, buf, n1);
        {
            uint64_t got = binary_LE_Uint64(binary_LittleEndian, buf);
            if (got != n1) {
                so_panic("LittleEndian: invalid decoded n1");
            }
        }
        buf = binary_LE_AppendUint64(binary_LittleEndian, so_slice(so_byte, buf, 0, 0), n2);
        {
            uint64_t got = binary_LE_Uint64(binary_LittleEndian, buf);
            if (got != n2) {
                so_panic("LittleEndian: invalid decoded n2");
            }
        }
    }
    return 0;
}
