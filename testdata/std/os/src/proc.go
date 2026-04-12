package main

import (
	"solod.dev/so/os"
)

func procTest() {
	{
		// Getpid.
		pid := os.Getpid()
		if pid <= 0 {
			panic("Getpid: invalid")
		}
	}
	{
		// Getppid.
		ppid := os.Getppid()
		if ppid < 0 {
			panic("Getppid: invalid")
		}
	}
	{
		// Getuid.
		uid := os.Getuid()
		if uid < 0 {
			panic("Getuid: invalid")
		}
	}
	{
		// Geteuid.
		euid := os.Geteuid()
		if euid < 0 {
			panic("Geteuid: invalid")
		}
	}
	{
		// Getgid.
		gid := os.Getgid()
		if gid < 0 {
			panic("Getgid: invalid")
		}
	}
	{
		// Getegid.
		egid := os.Getegid()
		if egid < 0 {
			panic("Getegid: invalid")
		}
	}
	{
		// Getwd.
		var wdBuf [os.MaxPathLen]byte
		wd, err := os.Getwd(wdBuf[:])
		if err != nil {
			panic("Getwd failed")
		}
		if len(wd) == 0 {
			panic("Getwd: empty")
		}
		// Should start with '/'.
		if wd[0] != '/' {
			panic("Getwd: not absolute")
		}
	}
	{
		// Hostname.
		var hostBuf [os.MaxNameLen]byte
		name, err := os.Hostname(hostBuf[:])
		if err != nil {
			panic("Hostname failed")
		}
		if len(name) == 0 {
			panic("Hostname: empty")
		}
	}
}
