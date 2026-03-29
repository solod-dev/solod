package main

import (
	"solod.dev/so/io"
	"solod.dev/so/mem"
	"solod.dev/so/os"
)

func seekTest() {
	{
		// Seek.
		name := "test_seek.txt"
		f, err := os.Create(name)
		if err != nil {
			panic("Create failed")
		}
		defer os.Remove(name)
		f.Write([]byte("abcdef"))
		pos, err := f.Seek(0, io.SeekStart)
		if err != nil {
			panic("Seek failed")
		}
		if pos != 0 {
			panic("Seek: wrong position")
		}
		buf := make([]byte, 6)
		n, err := f.Read(buf)
		if err != nil {
			panic("Read after Seek failed")
		}
		if string(buf[:n]) != "abcdef" {
			panic("Seek: wrong data")
		}
		f.Close()
	}
	{
		// ReadAt.
		name := "test_readat.txt"
		err := os.WriteFile(name, []byte("hello world"), 0o666)
		if err != nil {
			panic("WriteFile failed")
		}
		defer os.Remove(name)
		f, err := os.Open(name)
		if err != nil {
			panic("Open failed")
		}
		buf := make([]byte, 5)
		n, err := f.ReadAt(buf, 6)
		if err != nil {
			panic("ReadAt failed")
		}
		if n != 5 {
			panic("ReadAt: wrong count")
		}
		if string(buf[:n]) != "world" {
			panic("ReadAt: wrong data")
		}
		f.Close()
	}
	{
		// WriteAt.
		name := "test_writeat.txt"
		f, err := os.Create(name)
		if err != nil {
			panic("Create failed")
		}
		defer os.Remove(name)
		f.Write([]byte("hello world"))
		_, err = f.WriteAt([]byte("WORLD"), 6)
		if err != nil {
			panic("WriteAt failed")
		}
		f.Close()

		b, err := os.ReadFile(nil, name)
		if err != nil {
			panic("ReadFile failed")
		}
		defer mem.FreeSlice(nil, b)
		if string(b) != "hello WORLD" {
			panic("WriteAt: wrong data")
		}
	}
}
