#include "so/builtin/builtin.h"

// maps_wymum performs 128-bit multiply-and-mix using hardware support.
static inline uint64_t _maps_wymum(uint64_t a, uint64_t b) {
    __uint128_t r = (__uint128_t)a * b;
    return (uint64_t)(r >> 64) ^ (uint64_t)r;
}

// maps_wyr8 reads 8 bytes as a little-endian uint64.
static inline uint64_t _maps_wyr8(const uint8_t* p) {
    uint64_t v;
    memcpy(&v, p, 8);
    return v;
}

// maps_wyr4 reads 4 bytes as a little-endian uint64.
static inline uint64_t _maps_wyr4(const uint8_t* p) {
    uint32_t v;
    memcpy(&v, p, 4);
    return (uint64_t)v;
}

// maps_hash computes wyhash inline using memcpy-based reads
// and __uint128_t multiply.
static inline so_int maps_hash(const void* key, size_t len) {
    const uint8_t* p = (const uint8_t*)key;
    const uint64_t wyp0 = 0xa0761d6478bd642fULL;
    const uint64_t wyp1 = 0xe7037ed1a0b428dbULL;
    uint64_t seed = _maps_wymum(wyp0, wyp1);
    uint64_t a = 0, b = 0;
    if (len > 16) {
        for (size_t i = 0; i + 16 <= len; i += 16) {
            seed = _maps_wymum(_maps_wyr8(p + i) ^ wyp1,
                               _maps_wyr8(p + i + 8) ^ seed);
        }
        a = _maps_wyr8(p + len - 16);
        b = _maps_wyr8(p + len - 8);
    } else if (len >= 4) {
        a = (_maps_wyr4(p) << 32) | _maps_wyr4(p + ((len >> 3) << 2));
        b = (_maps_wyr4(p + len - 4) << 32) |
            _maps_wyr4(p + len - 4 - ((len >> 3) << 2));
    } else if (len > 0) {
        a = ((uint64_t)p[0] << 16) | ((uint64_t)p[len >> 1] << 8) |
            (uint64_t)p[len - 1];
    }
    uint64_t r = _maps_wymum(wyp1 ^ (uint64_t)len, _maps_wymum(a ^ wyp1, b ^ seed));
    return (so_int)(r >> 16);  // upper 48 bits is the hash value
}

// maps_hashString hashes a string key by its content.
static inline so_int maps_hashString(void* key_ptr) {
    so_String* s = (so_String*)key_ptr;
    return maps_hash(s->ptr, s->len);
}

// maps_keyHash hashes a key, dispatching to string or inline hash.
#define maps_keyHash(K, key_ptr) _Generic((K){0}, \
    so_String: maps_hashString(key_ptr),          \
    default: maps_hash((key_ptr), sizeof(K)))
