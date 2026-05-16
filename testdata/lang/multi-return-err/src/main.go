package main

import "example/sub"

type Reader interface {
	Read(buf int) (int, error)
}

type File struct {
	size int
}

func makeFile(size int) (File, error) {
	return File{size: size}, nil
}

func (f *File) Read(buf int) (int, error) {
	_ = buf
	return f.size, nil
}

var file File

func returnPtr() (*File, error) {
	return &file, nil
}

type point struct {
	x, y int
}

func makePoint(x, y int) (point, error) {
	return point{x: x, y: y}, nil
}

func makeSubPoint(x, y int) (sub.Point, error) {
	return sub.Point{X: x, Y: y}, nil
}

func divide(a, b int) (int, error) {
	return a / b, nil
}

func returnRune() (rune, error) {
	return 'x', nil
}

func returnString() (string, error) {
	return "hello", nil
}

func returnSlice() ([]int, error) {
	return []int{1, 2, 3}, nil
}

func forwardCall() (int, error) {
	return divide(10, 3)
}

func main() {
	{
		// Destructure into new variables.
		q, err := divide(10, 3)
		_ = q
		_ = err

		// Blank identifier.
		_, err2 := divide(10, 3)
		_ = err2
		r3, _ := divide(10, 3)
		_ = r3

		// Partial reassignment.
		r4, err2 := divide(10, 3)
		_ = r4

		// Assign to existing variables.
		q = 0
		err = nil
		q, err = divide(20, 7)
	}
	{
		// If-init with multi-return.
		f := File{size: 42}
		if n, err := f.Read(64); err != nil {
			_ = n
		}
	}
	{
		// Various return types.
		var err error
		_ = err
		run, err := returnRune()
		_ = run
		str, err := returnString()
		_ = str
		slice, err := returnSlice()
		_ = slice
		// struc, err := returnStruct()
		// _ = struc
		ptr, err := returnPtr()
		_ = ptr
		// iface, err := returnIface()
		// _ = iface
	}
	{
		// Forward call.
		q, err := forwardCall()
		_ = q
		_ = err
	}
	{
		// Custom exported struct + error.
		f, err := makeFile(42)
		if f.size != 42 || err != nil {
			panic("Custom exported struct failed")
		}
	}
	{
		// Custom unexported struct + error.
		p, err := makePoint(1, 2)
		if p.x != 1 || p.y != 2 || err != nil {
			panic("Custom unexported struct failed")
		}
	}
	{
		// Custom struct from another package + error.
		sp1, err := makeSubPoint(1, 2)
		if sp1.X != 1 || sp1.Y != 2 || err != nil {
			panic("Custom struct from another package failed")
		}
		sp2, err := sub.MakePoint(3, 4)
		if sp2.X != 3 || sp2.Y != 4 || err != nil {
			panic("Custom struct from another package failed")
		}
	}
}
