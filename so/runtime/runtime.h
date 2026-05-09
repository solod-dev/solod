#include "so/builtin/builtin.h"

#if __STDC_HOSTED__

#if defined(so_build_darwin) || defined(so_build_netbsd) || defined(so_build_openbsd)
#include <stdlib.h>
#elif defined(so_build_linux) || defined(so_build_freebsd) || defined(so_build_dragonfly)
#include <sys/random.h>
#elif defined(so_build_wasm)
#include <unistd.h>
#endif

// Seed returns a random 64-bit seed.
static inline uint64_t runtime_Seed(void) {
    uint64_t seed = 0;
#if defined(so_build_darwin) || defined(so_build_netbsd) || defined(so_build_openbsd)
    arc4random_buf(&seed, sizeof(seed));
#elif defined(so_build_linux) || defined(so_build_freebsd) || defined(so_build_dragonfly)
    ssize_t n = getrandom(&seed, sizeof(seed), 0);
    if (n != sizeof(seed)) {
        so_panic("runtime: cryptographic random not available");
    }
#elif defined(so_build_wasm)
    if (getentropy(&seed, sizeof(seed)) != 0) {
        so_panic("runtime: cryptographic random not available");
    }
#else
    so_panic("runtime: cryptographic random not available");
#endif
    return seed;
}

#else

// Deterministic xorshift64 fallback for freestanding environments.
static inline uint64_t runtime_Seed(void) {
    static uint64_t rng_state = 0xdeadbeefcafebabeULL;
    uint64_t x = rng_state;
    x ^= x << 13; x ^= x >> 7; x ^= x << 17;
    rng_state = x;
    return x;
}

#endif  // __STDC_HOSTED__

#define runtime_buildVersion so_str(so_version)

#if !__STDC_HOSTED__
#define runtime_GOOS so_str("bare")
#elif defined(so_build_darwin)
#define runtime_GOOS so_str("darwin")
#elif defined(so_build_linux)
#define runtime_GOOS so_str("linux")
#elif defined(so_build_freebsd)
#define runtime_GOOS so_str("freebsd")
#elif defined(so_build_netbsd)
#define runtime_GOOS so_str("netbsd")
#elif defined(so_build_openbsd)
#define runtime_GOOS so_str("openbsd")
#elif defined(so_build_dragonfly)
#define runtime_GOOS so_str("dragonfly")
#elif defined(so_build_wasm)
#define runtime_GOOS so_str("wasip1")
#elif defined(so_build_windows)
#define runtime_GOOS so_str("windows")
#else
#define runtime_GOOS so_str("unknown")
#endif

#if defined(so_build_amd64)
#define runtime_GOARCH so_str("amd64")
#elif defined(so_build_arm64)
#define runtime_GOARCH so_str("arm64")
#elif defined(so_build_riscv64)
#define runtime_GOARCH so_str("riscv64")
#elif defined(so_build_i386)
#define runtime_GOARCH so_str("386")
#elif defined(so_build_wasm32)
#define runtime_GOARCH so_str("wasm")
#else
#define runtime_GOARCH so_str("unknown")
#endif
