#include "so/builtin/builtin.h"

// equalBytes compares two byte keys for equality.
static inline bool maps_equalBytes(so_Slice a, so_Slice b) {
    return a.len == b.len && (a.len == 0 || memcmp(a.ptr, b.ptr, a.len) == 0);
}

// equalString compares two string keys by their content.
static inline bool maps_equalString(so_Slice a, so_Slice b) {
    return so_string_eq(*(so_String*)a.ptr, *(so_String*)b.ptr);
}

// equal compares two typed key pointers for equality.
#define maps_keyEqual(K, a, b)                                       \
    _Generic((K){0},                                                 \
        so_String: so_string_eq(*(so_String*)(a), *(so_String*)(b)), \
        default: memcmp((a), (b), sizeof(K)) == 0)

// Map is a generic hashmap similar to Go's built-in map[K]V.
typedef maps_ByteMap maps_Map;

// New creates a new Map with the given minimal capacity
// using the provided allocator (or the default allocator if nil).
#define maps_New(K, V, a, size) ({                                    \
    maps_ByteMap _m = maps_NewByteMap((a), (size), (so_int)sizeof(K), \
                                      (so_int)sizeof(V));             \
    _m.equalFn = _Generic((K){0},                                     \
        so_String: maps_equalString,                                  \
        default: maps_equalBytes);                                    \
    _m;                                                               \
})

// Has returns true if the given key is in the map.
#define maps_Map_Has(K, V, m, key) ({                     \
    K _key = (key);                                       \
    bool _found = false;                                  \
    maps_ByteMap* _m = (m);                               \
    if (_m->hdib.len > 0) {                               \
        so_int _hash = maps_keyHash(K, &_key);            \
        so_int _i = _hash & _m->mask;                     \
        uint64_t* _hdib = (uint64_t*)_m->hdib.ptr;        \
        K* _keys = (K*)_m->keys.ptr;                      \
        so_int _dist = 1;                                 \
        for (;;) {                                        \
            uint64_t _ehdib = _hdib[_i];                  \
            if ((so_int)(_ehdib & 0xFFFF) < _dist) break; \
            if ((so_int)(_ehdib >> 16) == _hash &&        \
                maps_keyEqual(K, &_key, &_keys[_i])) {    \
                _found = true;                            \
                break;                                    \
            }                                             \
            _i = (_i + 1) & _m->mask;                     \
            _dist++;                                      \
        }                                                 \
    }                                                     \
    _found;                                               \
})

// Get returns the value for the given key,
// or the zero value if the key is not in the map.
#define maps_Map_Get(K, V, m, key) ({                     \
    K _key = (key);                                       \
    V _val;                                               \
    memset(&_val, 0, sizeof(V));                          \
    maps_ByteMap* _m = (m);                               \
    if (_m->hdib.len > 0) {                               \
        so_int _hash = maps_keyHash(K, &_key);            \
        so_int _i = _hash & _m->mask;                     \
        uint64_t* _hdib = (uint64_t*)_m->hdib.ptr;        \
        K* _keys = (K*)_m->keys.ptr;                      \
        so_int _dist = 1;                                 \
        for (;;) {                                        \
            uint64_t _ehdib = _hdib[_i];                  \
            if ((so_int)(_ehdib & 0xFFFF) < _dist) break; \
            if ((so_int)(_ehdib >> 16) == _hash &&        \
                maps_keyEqual(K, &_key, &_keys[_i])) {    \
                _val = ((V*)_m->vals.ptr)[_i];            \
                break;                                    \
            }                                             \
            _i = (_i + 1) & _m->mask;                     \
            _dist++;                                      \
        }                                                 \
    }                                                     \
    _val;                                                 \
})

// Set sets the value for the given key,
// overwriting any existing value.
#define maps_Map_Set(K, V, m, key, value)                      \
    do {                                                       \
        K _key = (key);                                        \
        V _val = (value);                                      \
        maps_ByteMap* _m = (m);                                \
        if (_m->len >= _m->growAt) {                           \
            maps_ByteMap_Resize(_m, (so_int)_m->hdib.len * 2); \
        }                                                      \
        so_int _hash = maps_keyHash(K, &_key);                 \
        uint64_t _ehdib = ((uint64_t)_hash << 16) | 1;         \
        so_int _i = _hash & _m->mask;                          \
        uint64_t* _hdib = (uint64_t*)_m->hdib.ptr;             \
        K* _keys = (K*)_m->keys.ptr;                           \
        V* _vals = (V*)_m->vals.ptr;                           \
        K _ekey = _key;                                        \
        V _eval = _val;                                        \
        for (;;) {                                             \
            if ((_hdib[_i] & 0xFFFF) == 0) {                   \
                _hdib[_i] = _ehdib;                            \
                _keys[_i] = _ekey;                             \
                _vals[_i] = _eval;                             \
                _m->len++;                                     \
                break;                                         \
            }                                                  \
            if ((_ehdib >> 16) == (_hdib[_i] >> 16) &&         \
                maps_keyEqual(K, &_ekey, &_keys[_i])) {        \
                _vals[_i] = _eval;                             \
                break;                                         \
            }                                                  \
            if ((_hdib[_i] & 0xFFFF) < (_ehdib & 0xFFFF)) {    \
                uint64_t _tmphdib = _ehdib;                    \
                _ehdib = _hdib[_i];                            \
                _hdib[_i] = _tmphdib;                          \
                K _tmpk = _ekey;                               \
                _ekey = _keys[_i];                             \
                _keys[_i] = _tmpk;                             \
                V _tmpv = _eval;                               \
                _eval = _vals[_i];                             \
                _vals[_i] = _tmpv;                             \
            }                                                  \
            _i = (_i + 1) & _m->mask;                          \
            _ehdib++;                                          \
        }                                                      \
    } while (0)

// Delete removes the key and its value from the map.
// If the key is not in the map, does nothing.
#define maps_Map_Delete(K, V, m, key)                        \
    do {                                                     \
        K _key = (key);                                      \
        maps_ByteMap* _m = (m);                              \
        if (_m->hdib.len == 0) break;                        \
        so_int _hash = maps_keyHash(K, &_key);               \
        so_int _i = _hash & _m->mask;                        \
        uint64_t* _hdib = (uint64_t*)_m->hdib.ptr;           \
        K* _keys = (K*)_m->keys.ptr;                         \
        V* _vals = (V*)_m->vals.ptr;                         \
        so_int _dist = 1;                                    \
        for (;;) {                                           \
            if ((so_int)(_hdib[_i] & 0xFFFF) < _dist) break; \
            if ((so_int)(_hdib[_i] >> 16) == _hash &&        \
                maps_keyEqual(K, &_key, &_keys[_i])) {       \
                for (;;) {                                   \
                    so_int _prev = _i;                       \
                    _i = (_i + 1) & _m->mask;                \
                    if ((_hdib[_i] & 0xFFFF) <= 1) {         \
                        _hdib[_prev] = 0;                    \
                        memset(&_keys[_prev], 0, sizeof(K)); \
                        memset(&_vals[_prev], 0, sizeof(V)); \
                        break;                               \
                    }                                        \
                    _hdib[_prev] = _hdib[_i] - 1;            \
                    _keys[_prev] = _keys[_i];                \
                    _vals[_prev] = _vals[_i];                \
                }                                            \
                _m->len--;                                   \
                break;                                       \
            }                                                \
            _i = (_i + 1) & _m->mask;                        \
            _dist++;                                         \
        }                                                    \
    } while (0)

// Len returns the number of key-value pairs in the map.
#define maps_Map_Len(K, V, m) \
    maps_ByteMap_Len(m)

// Free frees internal resources used by the map.
// If the map is already freed, does nothing.
#define maps_Map_Free(K, V, m) \
    maps_ByteMap_Free(m)
