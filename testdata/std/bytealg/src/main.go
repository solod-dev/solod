package main

import "solod.dev/so/bytealg"

func main() {
	{
		// IndexRabinKarp.
		b := []byte("go is fun")
		idx := bytealg.IndexRabinKarp(b, []byte("is"))
		if idx != 3 {
			panic("IndexRabinKarp failed")
		}
	}
	{
		// LastIndexRabinKarp.
		b := []byte("hello")
		idx := bytealg.LastIndexRabinKarp(b, []byte("l"))
		if idx != 3 {
			panic("LastIndexRabinKarp failed")
		}
	}
	{
		// Compare.
		b := []byte("abc")
		if bytealg.Compare(b, []byte("abb")) <= 0 {
			panic("Compare failed")
		}
		if bytealg.Compare(b, []byte("abd")) >= 0 {
			panic("Compare failed")
		}
		if bytealg.Compare(b, []byte("abc")) != 0 {
			panic("Compare failed")
		}
	}
	{
		// Count and CountString.
		b := []byte("hello world")
		n := bytealg.Count(b, 'o')
		if n != 2 {
			panic("Count failed")
		}
		s := "hello world"
		n = bytealg.CountString(s, 'o')
		if n != 2 {
			panic("CountString failed")
		}
	}
	{
		// Equal.
		a := []byte("hello")
		b := []byte("hello")
		if !bytealg.Equal(a, b) {
			panic("Equal failed")
		}
		c := []byte("world")
		if bytealg.Equal(a, c) {
			panic("Equal failed")
		}
	}
	{
		// IndexByte and IndexByteString.
		b := []byte("hello")
		idx := bytealg.IndexByte(b, 'l')
		if idx != 2 {
			panic("IndexByte failed")
		}
		s := "hello"
		idx = bytealg.IndexByteString(s, 'l')
		if idx != 2 {
			panic("IndexByteString failed")
		}
	}
	{
		// LastIndexByte and LastIndexByteString.
		b := []byte("hello")
		idx := bytealg.LastIndexByte(b, 'l')
		if idx != 3 {
			panic("LastIndexByte failed")
		}
		s := "hello"
		idx = bytealg.LastIndexByteString(s, 'l')
		if idx != 3 {
			panic("LastIndexByteString failed")
		}
	}
}
