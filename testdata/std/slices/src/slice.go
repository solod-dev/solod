package main

import "solod.dev/so/slices"

func sliceTest() {
	{
		// Make a slice.
		s := slices.Make[int](nil, 3)
		s[0] = 11
		s[1] = 22
		s[2] = 33
		if len(s) != 3 || cap(s) != 3 {
			panic("Make failed")
		}
		if s[0] != 11 || s[1] != 22 || s[2] != 33 {
			panic("Make failed")
		}
		slices.Free(nil, s)
	}
	{
		// Append within capacity.
		s := slices.MakeCap[int](nil, 0, 8)
		s = slices.Append(nil, s, 10, 20, 30)
		if len(s) != 3 || s[0] != 10 || s[1] != 20 || s[2] != 30 {
			panic("Append failed")
		}
		slices.Free(nil, s)
	}
	{
		// Append that triggers growth.
		s := slices.MakeCap[int](nil, 0, 2)
		s = slices.Append(nil, s, 1, 2)
		s = slices.Append(nil, s, 3, 4, 5)
		if len(s) != 5 || s[0] != 1 || s[4] != 5 {
			panic("Append grow failed")
		}
		slices.Free(nil, s)
	}
	{
		// Append to nil slice.
		var s []int
		s = slices.Append(nil, s, 10, 20, 30)
		if len(s) != 3 || s[0] != 10 || s[1] != 20 || s[2] != 30 {
			panic("Append to nil failed")
		}
		slices.Free(nil, s)
	}
	{
		// Extend from another slice.
		s := slices.MakeCap[int](nil, 0, 8)
		other := []int{100, 200, 300}
		s = slices.Extend(nil, s, other)
		if len(s) != 3 || s[0] != 100 || s[2] != 300 {
			panic("Extend failed")
		}
		slices.Free(nil, s)
	}
	{
		// Extend a nil slice.
		var s []int
		other := []int{10, 20, 30}
		s = slices.Extend(nil, s, other)
		if len(s) != 3 || s[0] != 10 || s[1] != 20 || s[2] != 30 {
			panic("Extend to nil failed")
		}
		slices.Free(nil, s)
	}
	{
		// Clone a slice.
		s1 := []int{11, 22, 33}
		s2 := slices.Clone(nil, s1)
		s2[0] = 99
		if s1[0] != 11 || s2[0] != 99 {
			panic("Clone failed")
		}
		slices.Free(nil, s2)
	}
	{
		// Equal slices.
		s1 := []int{1, 2, 3}
		s2 := []int{1, 2, 3}
		s3 := []int{1, 2, 4}
		s4 := []int{1, 2}
		s5 := []int{}
		var s6 []int = nil
		if !slices.Equal(s1, s2) {
			panic("want s1 == s2")
		}
		if slices.Equal(s1, s3) {
			panic("want s1 != s3")
		}
		if slices.Equal(s1, s4) {
			panic("want s1 != s4")
		}
		if !slices.Equal(s5, s6) {
			panic("want empty and nil slices equal")
		}
	}
	{
		// Equal string slices.
		s1 := []string{"a", "b", "c"}
		s2 := []string{"a", "b", "c"}
		s3 := []string{"a", "b", "d"}
		if !slices.Equal(s1, s2) {
			panic("want s1 == s2")
		}
		if slices.Equal(s1, s3) {
			panic("want s1 != s3")
		}
	}
	{
		// Equal struct slices.
		type point struct {
			x, y int
		}
		s1 := []point{{1, 2}, {3, 4}}
		s2 := []point{{1, 2}, {3, 4}}
		s3 := []point{{1, 2}, {3, 5}}
		if !slices.Equal(s1, s2) {
			panic("want s1 == s2")
		}
		if slices.Equal(s1, s3) {
			panic("want s1 != s3")
		}
	}
	{
		// Index of an element.
		ints := []int{10, 20, 30, 20}
		if slices.Index(ints, 20) != 1 {
			panic("Index failed")
		}
		if slices.Index(ints, 40) != -1 {
			panic("Index failed")
		}
		strs := []string{"a", "b", "c", "b"}
		if slices.Index(strs, "b") != 1 {
			panic("Index failed")
		}
		if slices.Index(strs, "d") != -1 {
			panic("Index failed")
		}
	}
	{
		// Contains an element.
		ints := []int{10, 20, 30, 20}
		if !slices.Contains(ints, 20) {
			panic("Contains failed")
		}
		if slices.Contains(ints, 40) {
			panic("Contains failed")
		}
		strs := []string{"a", "b", "c", "b"}
		if !slices.Contains(strs, "b") {
			panic("Contains failed")
		}
		if slices.Contains(strs, "d") {
			panic("Contains failed")
		}
	}
}
