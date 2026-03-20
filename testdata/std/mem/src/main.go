package main

import (
	"solod.dev/so/mem"
)

type Point struct {
	x, y int
}

func withDefer() {
	p := mem.Alloc[Point](nil)
	defer mem.Free(nil, p)

	p.x = 11
	p.y = 22
	if p.x != 11 || p.y != 22 {
		panic("unexpected value")
	}
}

func main() {
	{
		// TryAlloc and Free.
		p, err := mem.TryAlloc[Point](mem.System)
		if err != nil {
			panic("Alloc: allocation failed")
		}
		p.x = 11
		p.y = 22
		if p.x != 11 || p.y != 22 {
			panic("Alloc: unexpected value")
		}
		mem.Free(mem.System, p)
	}
	{
		// TryAllocSlice and FreeSlice.
		slice, err := mem.TryAllocSlice[int](mem.System, 3, 3)
		if err != nil {
			panic("AllocSlice: allocation failed")
		}
		slice[0] = 11
		slice[1] = 22
		slice[2] = 33
		if slice[0] != 11 || slice[1] != 22 || slice[2] != 33 {
			panic("AllocSlice: unexpected value")
		}
		mem.FreeSlice(mem.System, slice)
	}
	{
		// Alloc/Free with default allocator.
		p := mem.Alloc[Point](nil)
		p.x = 11
		p.y = 22
		if p.x != 11 || p.y != 22 {
			panic("New: unexpected value")
		}
		mem.Free(nil, p)
	}
	{
		// AllocSlice/FreeSlice with default allocator.
		slice := mem.AllocSlice[int](nil, 3, 3)
		slice[0] = 11
		slice[1] = 22
		slice[2] = 33
		if slice[0] != 11 || slice[1] != 22 || slice[2] != 33 {
			panic("NewSlice: unexpected value")
		}
		mem.FreeSlice(nil, slice)
	}
	{
		// Free with nil or an empty slice.
		var p *Point
		mem.Free(nil, p)
		var empty []int
		mem.FreeSlice(nil, empty)
	}
	{
		// Free string.
		b := mem.AllocSlice[byte](nil, 3, 3)
		b[0] = 'h'
		b[1] = 'i'
		b[2] = '!'
		s1 := string(b)
		mem.FreeString(nil, s1)
		s2 := ""
		mem.FreeString(nil, s2)
	}
	withDefer()
}
