package os

import (
	"unsafe"

	"solod.dev/so/c"
)

// Args holds the command-line arguments, starting with the program name.
//
//so:extern
var Args []string

// Getegid returns the numeric effective group id of the caller.
func Getegid() int {
	gid := getegid()
	return int(gid)
}

// Geteuid returns the numeric effective user id of the caller.
func Geteuid() int {
	uid := geteuid()
	return int(uid)
}

// Getgid returns the numeric group id of the caller.
func Getgid() int {
	gid := getgid()
	return int(gid)
}

// Getpid returns the process id of the caller.
func Getpid() int {
	pid := getpid()
	return int(pid)
}

// Getppid returns the process id of the caller's parent.
func Getppid() int {
	ppid := getppid()
	return int(ppid)
}

// Getuid returns the numeric user id of the caller.
func Getuid() int {
	uid := getuid()
	return int(uid)
}

// Getwd returns an absolute path name corresponding to the
// current directory.
//
// Writes the result into buf. Panics if buf is empty.
// The returned string is a view into buf.
func Getwd(buf []byte) (string, error) {
	bufPtr := unsafe.SliceData(buf)
	if bufPtr == nil {
		panic("os: empty buffer")
	}
	cwd := getcwd((*c.Char)(bufPtr), uintptr(len(buf)))
	if cwd == nil {
		return "", mapError()
	}
	return c.String(cwd), nil
}

// Hostname returns the host name reported by the kernel.
//
// Writes the result into buf. Panics if buf is empty.
// The returned string is a view into buf.
func Hostname(buf []byte) (string, error) {
	name := (*c.Char)(unsafe.SliceData(buf))
	if name == nil {
		panic("os: empty buffer")
	}
	if gethostname(name, uintptr(len(buf))) != 0 {
		return "", mapError()
	}
	return c.String(name), nil
}

// Exit causes the current program to exit with the given status code.
// Conventionally, code zero indicates success, non-zero an error.
func Exit(code int) {
	exit(c.Int(code))
}
