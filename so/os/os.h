#include <stdbool.h>
#include <stdio.h>
#include <errno.h>

#define os_file FILE

// File represents an open file descriptor.
typedef struct {
    os_file* fd;
    bool closed;
} os_File;

// Error codes.
#define os_EACCES EACCES
#define os_EEXIST EEXIST
#define os_EISDIR EISDIR
#define os_ENOENT ENOENT
#define os_ENOTDIR ENOTDIR
#define os_EPERM EPERM
