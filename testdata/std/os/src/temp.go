package main

import (
	"solod.dev/so/mem"
	"solod.dev/so/os"
	"solod.dev/so/strings"
)

func tempTest() {
	buf := make([]byte, os.MaxPathLen)
	{
		// TempDir.
		td := os.TempDir()
		if len(td) == 0 {
			panic("TempDir: empty")
		}
	}
	{
		// CreateTemp.
		f, err := os.CreateTemp(buf, "", "sotest")
		if err != nil {
			panic("CreateTemp failed")
		}
		name := f.Name()
		if len(name) == 0 {
			panic("CreateTemp: empty name")
		}
		// Name should contain the pattern prefix.
		if !strings.Contains(name, "sotest") {
			panic("CreateTemp: name missing pattern")
		}
		f.Write([]byte("temp data"))
		f.Close()

		// Verify the file exists.
		b, err := os.ReadFile(nil, name)
		if err != nil {
			panic("ReadFile temp failed")
		}
		if string(b) != "temp data" {
			panic("CreateTemp: wrong data")
		}
		mem.FreeSlice(nil, b)
		os.Remove(name)
	}
	{
		// CreateTemp with specific dir.
		td := os.TempDir()
		f, err := os.CreateTemp(buf, td, "myprefix")
		if err != nil {
			panic("CreateTemp dir failed")
		}
		name := f.Name()
		if !strings.Contains(name, "myprefix") {
			panic("CreateTemp dir: missing pattern")
		}
		if !strings.HasPrefix(name, td) {
			panic("CreateTemp dir: wrong dir")
		}
		f.Close()
		os.Remove(name)
	}
	{
		// MkdirTemp.
		dir, err := os.MkdirTemp(buf, "", "sotest")
		if err != nil {
			panic("MkdirTemp failed")
		}
		if len(dir) == 0 {
			panic("MkdirTemp: empty")
		}
		if !strings.Contains(dir, "sotest") {
			panic("MkdirTemp: name missing pattern")
		}

		// Verify it's a directory.
		fi, err := os.Stat(dir)
		if err != nil {
			panic("Stat MkdirTemp failed")
		}
		if !fi.IsDir() {
			panic("MkdirTemp: not a directory")
		}

		os.Remove(dir)
	}
}
