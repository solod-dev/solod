package main

import (
	"solod.dev/so/mem"
	"solod.dev/so/slices"
)

func main() {
	{
		// Append within capacity.
		s := mem.AllocSlice[int](nil, 0, 8)
		s = slices.Append(nil, s, 10, 20, 30)
		if len(s) != 3 || s[0] != 10 || s[1] != 20 || s[2] != 30 {
			panic("Append: unexpected value")
		}
		mem.FreeSlice(nil, s)
	}
	{
		// Append that triggers growth.
		s := mem.AllocSlice[int](nil, 0, 2)
		s = slices.Append(nil, s, 1, 2)
		s = slices.Append(nil, s, 3, 4, 5)
		if len(s) != 5 || s[0] != 1 || s[4] != 5 {
			panic("Append grow: unexpected value")
		}
		mem.FreeSlice(nil, s)
	}
	{
		// Extend from another slice.
		s := mem.AllocSlice[int](nil, 0, 8)
		other := []int{100, 200, 300}
		s = slices.Extend(nil, s, other)
		if len(s) != 3 || s[0] != 100 || s[2] != 300 {
			panic("Extend: unexpected value")
		}
		mem.FreeSlice(nil, s)
	}
	{
		// Clone a slice.
		s1 := []int{11, 22, 33}
		s2 := slices.Clone(nil, s1)
		s2[0] = 99
		if s1[0] != 11 || s2[0] != 99 {
			panic("Clone: unexpected value")
		}
		mem.FreeSlice(nil, s2)
	}
}
