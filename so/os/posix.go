package os

import "solod.dev/so/c"

//so:include <dirent.h>
//so:include <fcntl.h>
//so:include.c <sys/stat.h>
//so:include.c <unistd.h>

//so:extern S_IFMT
const sIFMT = 0170000 // type of file mask
//so:extern S_IFIFO
const sIFIFO = 0010000 // named pipe (fifo)
//so:extern S_IFCHR
const sIFCHR = 0020000 // character special
//so:extern S_IFDIR
const sIFDIR = 0040000 // directory
//so:extern S_IFBLK
const sIFBLK = 0060000 // block special
//so:extern S_IFREG
const sIFREG = 0100000 // regular
//so:extern S_IFLNK
const sIFLNK = 0120000 // symbolic link
//so:extern S_IFSOCK
const sIFSOCK = 0140000 // socket
//so:extern S_ISUID
const sISUID = 0004000 // set user id on execution
//so:extern S_ISGID
const sISGID = 0002000 // set group id on execution
//so:extern S_ISVTX
const sISVTX = 0001000 // directory restricted delete

//so:extern UTIME_OMIT
const utimeOmit = 0 // sentinel: leave file time unchanged

//so:extern
type gid_t uint32 // group ID
//so:extern
type mode_t uint32 // file mode bits
//so:extern
type pid_t int32 // process ID
//so:extern
type uid_t uint32 // user ID

//so:extern
type os_statResult struct {
	size    int64
	mode    mode_t
	modSec  int64
	modNsec int64
	dev     uint64
	ino     uint64
	ok      bool
}

//so:extern DIR
type os_dir struct{}

// dirent d_type constants.
//
//so:extern DT_UNKNOWN
const dtUnknown = 0

//so:extern DT_FIFO
const dtFIFO = 1

//so:extern DT_CHR
const dtCHR = 2

//so:extern DT_DIR
const dtDIR = 4

//so:extern DT_BLK
const dtBLK = 6

//so:extern DT_REG
const dtREG = 8

//so:extern DT_LNK
const dtLNK = 10

//so:extern DT_SOCK
const dtSOCK = 12

//so:extern
type os_readdirResult struct {
	nameLen int32
	dtype   uint8
	ok      bool
}

// int chdir(const char* path);
//
//so:extern
func chdir(path string) c.Int {
	_ = path
	return 0
}

// int chmod(const char* path, mode_t mode);
//
//so:extern
func chmod(path string, mode mode_t) c.Int {
	_, _ = path, mode
	return 0
}

// int chown(const char* path, uid_t owner, gid_t group);
//
//so:extern
func chown(path string, uid uid_t, gid gid_t) c.Int {
	_, _, _ = path, uid, gid
	return 0
}

// int closedir(DIR *dirp);
//
//so:extern
func closedir(dir *os_dir) c.Int {
	_ = dir
	return 0
}

// int close(int fd);
//
//so:extern close
func fdclose(fd c.Int) c.Int {
	_ = fd
	return 0
}

// FILE* fdopen(int fd, const char *mode);
//
//so:extern
func fdopen(fd c.Int, mode string) *os_file {
	_, _ = fd, mode
	return &os_file{}
}

// char* getcwd(char *buf, size_t size);
//
//so:extern
func getcwd(buf *c.Char, size uintptr) *c.Char {
	_, _ = buf, size
	return nil
}

// uid_t geteuid(void);
//
//so:extern
func geteuid() uid_t {
	return 0
}

// gid_t getegid(void);
//
//so:extern
func getegid() gid_t {
	return 0
}

// gid_t getgid(void);
//
//so:extern
func getgid() gid_t {
	return 0
}

// pid_t getpid(void);
//
//so:extern
func getpid() pid_t {
	return 0
}

// pid_t getppid(void);
//
//so:extern
func getppid() pid_t {
	return 0
}

// uid_t getuid(void);
//
//so:extern
func getuid() uid_t {
	return 0
}

// int gethostname(char* name, size_t namelen);
//
//so:extern
func gethostname(name *c.Char, namelen uintptr) c.Int {
	_, _ = name, namelen
	return 0
}

// int lchown(const char* path, uid_t owner, gid_t group);
//
//so:extern
func lchown(path string, uid uid_t, gid gid_t) c.Int {
	_, _, _ = path, uid, gid
	return 0
}

// int link(const char *path1, const char *path2);
//
//so:extern
func link(old, new string) c.Int {
	_, _ = old, new
	return 0
}

// lstat wrapper (fills os_statResult).
// int lstat(const char* path, struct stat* buf);
//
//so:extern
func os_lstat(path string) os_statResult {
	_ = path
	return os_statResult{size: 42, mode: 0o777, ok: true}
}

// int mkdir(const char* path, mode_t mode);
//
//so:extern
func mkdir(path string, mode mode_t) c.Int {
	_, _ = path, mode
	return 0
}

// int mkstemp(char *template);
//
//so:extern
func mkstemp(tmpl *c.Char) c.Int {
	_ = tmpl
	return 0
}

// char* mkdtemp(char *template);
//
//so:extern
func mkdtemp(tmpl *c.Char) *c.Char {
	_ = tmpl
	b := []byte("example/tmpdir")
	return (*c.Char)(&b[0])
}

// DIR* opendir(const char *name);
//
//so:extern
func opendir(name string) *os_dir {
	_ = name
	return nil
}

// int open(const char* path, int oflag, ...);
//
//so:extern open
func posixOpen(path string, flags c.Int, mode c.Int) c.Int {
	_, _, _ = path, flags, mode
	return 42
}

// os_readdir_next reads the next directory entry into buf.
//
//so:extern
func os_readdir_next(dir *os_dir, buf *c.Char, bufsize uintptr) os_readdirResult {
	_, _, _ = dir, buf, bufsize
	return os_readdirResult{}
}

// ssize_t readlink(const char* path, char* buf, size_t bufsize);
//
//so:extern
func readlink(path string, buf *c.Char, bufsize uintptr) int {
	_, _, _ = path, buf, bufsize
	return 0
}

// stat wrapper (fills os_statResult).
// int stat(const char* path, struct stat* buf);
//
//so:extern
func os_stat(path string) os_statResult {
	_ = path
	return os_statResult{}
}

// int symlink(const char *path1, const char *path2);
//
//so:extern
func symlink(old, new string) c.Int {
	_, _ = old, new
	return 0
}

// int truncate(const char* path, off_t length);
//
//so:extern
func truncate(path string, size int64) c.Int {
	_, _ = path, size
	return 0
}

// int setenv(const char *name, const char *value, int overwrite);
//
//so:extern
func setenv(key, value string, overwrite c.Int) c.Int {
	_, _, _ = key, value, overwrite
	return 0
}

// int unsetenv(const char *name);
//
//so:extern
func unsetenv(key string) c.Int {
	_ = key
	return 0
}

// utimensat wrapper (passes timespec values as separate arguments).
//
//so:extern
func os_utimens(path string, asec, ansec, msec, mnsec int64) c.Int {
	_, _, _, _, _ = path, asec, ansec, msec, mnsec
	return 0
}
