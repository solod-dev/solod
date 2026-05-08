package main

import "solod.dev/so/cmp"

func main() {
	{
		// Compare numbers.
		a := 11
		b := 22
		if cmp.Compare(a, b) >= 0 {
			panic("Compare failed")
		}
		if cmp.Compare(a, a) != 0 {
			panic("Compare failed")
		}
	}
	{
		// Compare strings.
		a := "hello"
		b := "world"
		if cmp.Compare(a, b) >= 0 {
			panic("Compare failed")
		}
		if cmp.Compare(a, a) != 0 {
			panic("Compare failed")
		}
	}
	{
		// Equal numbers.
		a := 11
		b := 22
		if cmp.Equal(a, b) {
			panic("Equal failed")
		}
		if !cmp.Equal(a, a) {
			panic("Equal failed")
		}
	}
	{
		// Equal strings.
		a := "hello"
		b := "world"
		if cmp.Equal(a, b) {
			panic("Equal failed")
		}
		if !cmp.Equal(a, a) {
			panic("Equal failed")
		}
	}
	{
		// Less numbers.
		a := 11
		b := 22
		if !cmp.Less(a, b) {
			panic("Less failed")
		}
		if cmp.Less(b, a) {
			panic("Less failed")
		}
	}
	{
		// Less strings.
		a := "hello"
		b := "world"
		if !cmp.Less(a, b) {
			panic("Less failed")
		}
		if cmp.Less(b, a) {
			panic("Less failed")
		}
	}
}
