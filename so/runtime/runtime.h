#include <time.h>

#if defined(so_build_darwin) || defined(so_build_netbsd) || defined(so_build_openbsd)
#include <stdlib.h>
#elif defined(so_build_linux) || defined(so_build_freebsd) || defined(so_build_dragonfly)
#include <sys/random.h>
#endif

#include "so/builtin/builtin.h"

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
#else
    so_panic("runtime: cryptographic random not available");
#endif
    return seed;
}

#define runtime_buildVersion so_str(so_version)

#if defined(so_build_darwin)
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
#elif defined(so_build_windows)
#define runtime_GOOS so_str("windows")
#else
#define runtime_GOOS so_str("unknown")
#endif

#if defined(so_build_amd64)
#define runtime_GOARCH so_str("amd64")
#elif defined(so_build_arm64)
#define runtime_GOARCH so_str("arm64")
#else
#define runtime_GOARCH so_str("unknown")
#endif
