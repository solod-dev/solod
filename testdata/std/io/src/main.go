package main

import (
	"solod.dev/so/io"
	"solod.dev/so/mem"
)

type reader struct {
	b []byte
}

func (r *reader) Read(p []byte) (int, error) {
	if len(r.b) == 0 {
		return 0, io.EOF
	}
	n := copy(p, r.b)
	r.b = r.b[n:]
	return n, nil
}

type writer struct {
	b []byte
}

func (w *writer) Write(p []byte) (int, error) {
	w.b = append(w.b, p...)
	return len(p), nil
}

func main() {
	{
		// Copy.
		r := reader{b: []byte("hello world")}
		w := writer{b: make([]byte, 0, 11)}
		if _, err := io.Copy(&w, &r); err != nil {
			panic("Copy failed")
		}
		if string(w.b) != "hello world" {
			panic("Copy failed")
		}
	}
	{
		// CopyN.
		r := reader{b: []byte("hello world")}
		w := writer{b: make([]byte, 0, 5)}
		if _, err := io.CopyN(&w, &r, 5); err != nil {
			panic("CopyN failed")
		}
		if string(w.b) != "hello" {
			panic("CopyN failed")
		}
	}
	{
		// ReadAll.
		r := reader{b: []byte("hello world")}
		buf, err := io.ReadAll(nil, &r)
		if err != nil {
			panic("ReadAll failed")
		}
		if string(buf) != "hello world" {
			panic("ReadAll failed")
		}
		mem.FreeSlice(nil, buf)
	}
	{
		// ReadFull.
		r := reader{b: []byte("hello world")}
		buf := make([]byte, 11)
		if _, err := io.ReadFull(&r, buf); err != nil {
			panic("ReadFull failed")
		}
		if string(buf) != "hello world" {
			panic("ReadFull failed")
		}
	}
	{
		// WriteString.
		w := writer{b: make([]byte, 0, 11)}
		n, err := io.WriteString(&w, "hello world")
		if err != nil {
			panic("WriteString failed")
		}
		if n != 11 || string(w.b) != "hello world" {
			panic("WriteString failed")
		}
	}
	{
		// LimitReader.
		r := reader{b: []byte("hello world")}
		lr := io.LimitReader(&r, 5)
		buf := make([]byte, 5)
		if _, err := lr.Read(buf); err != nil {
			panic("LimitReader failed")
		}
		if string(buf) != "hello" {
			panic("LimitReader failed")
		}
	}
}
