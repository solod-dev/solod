//go:build ignore
#include <dirent.h>
#include <errno.h>
#include <fcntl.h>
#include <string.h>
#include <sys/stat.h>
#include "so/builtin/builtin.h"

// Stat result - flat struct filled by C helpers.
typedef struct {
    int64_t size;
    uint32_t mode;
    int64_t modSec;
    int64_t modNsec;
    uint64_t dev;
    uint64_t ino;
    bool ok;
} os_statResult;

// os_lstat fills result from lstat().
static os_statResult os_lstat(const char* path) {
    struct stat st;
    if (lstat(path, &st) != 0) return (os_statResult){.ok = false};
    return (os_statResult){
        .size = st.st_size,
        .mode = st.st_mode,
        .modSec = st.st_mtime,
        .modNsec = 0,  // fields differ on macos and linux, set to 0 for now
        .dev = st.st_dev,
        .ino = st.st_ino,
        .ok = true,
    };
}

// os_stat fills result from stat().
static os_statResult os_stat(const char* path) {
    struct stat st;
    if (stat(path, &st) != 0) return (os_statResult){.ok = false};
    return (os_statResult){
        .size = st.st_size,
        .mode = st.st_mode,
        .modSec = st.st_mtime,
        .modNsec = 0,  // fields differ on macos and linux, set to 0 for now
        .dev = st.st_dev,
        .ino = st.st_ino,
        .ok = true,
    };
}

// os_utimens sets access and modification times using utimensat.
// A tv_nsec of UTIME_OMIT leaves the corresponding time unchanged.
static int os_utimens(const char* path, int64_t asec, int64_t ansec, int64_t msec, int64_t mnsec) {
    struct timespec times[2] = {
        {.tv_sec = asec, .tv_nsec = ansec},
        {.tv_sec = msec, .tv_nsec = mnsec},
    };
    return utimensat(AT_FDCWD, path, times, 0);
}

// readdir result - one entry at a time.
typedef struct {
    int32_t nameLen;
    uint8_t dtype;
    bool ok;
} os_readdirResult;

// os_readdir_next reads the next directory entry.
// Copies d_name into buf. Returns {nameLen, dtype, ok}.
static os_readdirResult os_readdir_next(DIR* dir, char* buf, size_t bufsize) {
    errno = 0;
    struct dirent* ent = readdir(dir);
    if (ent == NULL) return (os_readdirResult){.ok = false};
    size_t n = strlen(ent->d_name);
    if (n >= bufsize) n = bufsize - 1;
    memcpy(buf, ent->d_name, n);
    buf[n] = '\0';
    return (os_readdirResult){.nameLen = (int32_t)n, .dtype = ent->d_type, .ok = true};
}
