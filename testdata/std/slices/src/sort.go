package main

import (
	"solod.dev/so/cmp"
	"solod.dev/so/slices"
)

func descInt(a, b any) int {
	va := *a.(*int)
	vb := *b.(*int)
	return vb - va
}

func sortTest() {
	var ints = [...]int{74, 59, 238, -784, 9845, 959, 905, 0, 0, 42, 7586, -5467984, 7586}
	var float64s = [...]float64{74.3, 59.0, 238.2, -784.0, 2.3, 9845.768, -959.7485, 905, 7.8, 7.8, 74.3, 59.0, 238.2, -784.0, 2.3}
	var strs = [...]string{"", "Hello", "foo", "bar", "foo", "f00", "%*&^*&^&", "***"}

	{
		// IsSorted: false on unsorted data.
		if slices.IsSorted(ints[:]) {
			panic("IsSorted: unsorted ints")
		}
		if slices.IsSorted(strs[:]) {
			panic("IsSorted: unsorted strs")
		}
		// IsSorted: true on sorted data.
		sorted := []int{1, 2, 3, 4, 5}
		if !slices.IsSorted(sorted) {
			panic("IsSorted: sorted ints")
		}
		sortedStrs := []string{"a", "b", "c"}
		if !slices.IsSorted(sortedStrs) {
			panic("IsSorted: sorted strs")
		}
	}
	{
		// IsSortedFunc: false on unsorted data.
		compare := cmp.FuncFor[int]()
		if slices.IsSortedFunc(ints[:], compare) {
			panic("IsSortedFunc: unsorted ints")
		}
		// IsSortedFunc: true on sorted data.
		sorted := []int{1, 2, 3, 4, 5}
		if !slices.IsSortedFunc(sorted, compare) {
			panic("IsSortedFunc: sorted ints")
		}
	}
	{
		// Sort ints.
		s := slices.Clone(nil, ints[:])
		slices.Sort(s)
		if !slices.IsSorted(s) {
			panic("Sort ints: not sorted")
		}
		if s[0] != -5467984 || s[12] != 9845 {
			panic("Sort ints: wrong values")
		}
		slices.Free(nil, s)
	}
	{
		// Sort float64s.
		s := slices.Clone(nil, float64s[:])
		slices.Sort(s)
		if !slices.IsSorted(s) {
			panic("Sort float64s: not sorted")
		}
		if s[0] != -959.7485 || s[14] != 9845.768 {
			panic("Sort float64s: wrong values")
		}
		slices.Free(nil, s)
	}
	{
		// Sort strings.
		s := slices.Clone(nil, strs[:])
		slices.Sort(s)
		if !slices.IsSorted(s) {
			panic("Sort strings: not sorted")
		}
		if s[0] != "" || s[7] != "foo" {
			panic("Sort strings: wrong values")
		}
		slices.Free(nil, s)
	}
	{
		// SortFunc (reverse order).
		s := slices.Clone(nil, ints[:])
		slices.SortFunc(s, descInt)
		if !slices.IsSortedFunc(s, descInt) {
			panic("SortFunc ints: not sorted")
		}
		if s[0] != 9845 || s[12] != -5467984 {
			panic("SortFunc ints: wrong values")
		}
		slices.Free(nil, s)
	}
	{
		// SortFunc with nil compare.
		type point struct{ x, y int }
		s := []point{{1, 2}, {3, 4}, {2, 3}}
		slices.SortFunc(s, nil)
		if !slices.IsSortedFunc(s, nil) {
			panic("SortFunc with nil: not sorted")
		}
		if s[0].x != 1 || s[0].y != 2 {
			panic("SortFunc with nil: wrong s[0]")
		}
		if s[1].x != 2 || s[1].y != 3 {
			panic("SortFunc with nil: wrong s[1]")
		}
		if s[2].x != 3 || s[2].y != 4 {
			panic("SortFunc with nil: wrong s[2]")
		}
	}
	{
		// SortStableFunc ints.
		s := slices.Clone(nil, ints[:])
		compare := cmp.FuncFor[int]()
		slices.SortStableFunc(s, compare)
		if !slices.IsSorted(s) {
			panic("SortStable ints: not sorted")
		}
		if s[0] != -5467984 || s[12] != 9845 {
			panic("SortStable ints: wrong values")
		}
		slices.Free(nil, s)
	}
	{
		// SortStableFunc float64s.
		s := slices.Clone(nil, float64s[:])
		compare := cmp.FuncFor[float64]()
		slices.SortStableFunc(s, compare)
		if !slices.IsSorted(s) {
			panic("SortStable float64s: not sorted")
		}
		if s[0] != -959.7485 || s[14] != 9845.768 {
			panic("SortStable float64s: wrong values")
		}
		slices.Free(nil, s)
	}
	{
		// SortStableFunc strings.
		s := slices.Clone(nil, strs[:])
		compare := cmp.FuncFor[string]()
		slices.SortStableFunc(s, compare)
		if !slices.IsSorted(s) {
			panic("SortStable strings: not sorted")
		}
		if s[0] != "" || s[7] != "foo" {
			panic("SortStable strings: wrong values")
		}
		slices.Free(nil, s)
	}
}

func minMaxTest() {
	{
		// Min and Max on ints.
		ints := []int{3, 1, 4, 1, 5, 9}
		if slices.Min(ints) != 1 {
			panic("Min ints: wrong value")
		}
		if slices.Max(ints) != 9 {
			panic("Max ints: wrong value")
		}
	}
	{
		// Min and Max on strings.
		strs := []string{"banana", "apple", "cherry"}
		if slices.Min(strs) != "apple" {
			panic("Min strings: wrong value")
		}
		if slices.Max(strs) != "cherry" {
			panic("Max strings: wrong value")
		}
	}
}
