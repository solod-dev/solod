package main

import (
	"github.com/nalgeon/solod/so/mem"
	"github.com/nalgeon/solod/so/stringslite"
)

func main() {
	{
		// Clone.
		s := "hello"
		c := stringslite.Clone(nil, s)
		if c != s {
			panic("Clone failed")
		}
		mem.FreeString(nil, c)
	}
	{
		// Cut.
		before, after := stringslite.Cut("hello world", " ")
		if before != "hello" || after != "world" {
			panic("Cut failed")
		}
	}
	{
		// CutPrefix.
		after, found := stringslite.CutPrefix("hello world", "hello ")
		if after != "world" || !found {
			panic("CutPrefix failed")
		}
	}
	{
		// CutSuffix.
		before, found := stringslite.CutSuffix("hello world", " world")
		if before != "hello" || !found {
			panic("CutSuffix failed")
		}
	}
	{
		// HasPrefix.
		if !stringslite.HasPrefix("hello world", "hello") {
			panic("HasPrefix failed")
		}
		if stringslite.HasPrefix("hello world", "world") {
			panic("HasPrefix failed")
		}
	}
	{
		// HasSuffix.
		if !stringslite.HasSuffix("hello world", "world") {
			panic("HasSuffix failed")
		}
		if stringslite.HasSuffix("hello world", "hello") {
			panic("HasSuffix failed")
		}
	}
	{
		// Index.
		idx := stringslite.Index("hello world", "world")
		if idx != 6 {
			panic("Index failed")
		}
	}
	{
		// IndexByte.
		idx := stringslite.IndexByte("hello world", 'o')
		if idx != 4 {
			panic("IndexByte failed")
		}
	}
	{
		// TrimPrefix.
		s := stringslite.TrimPrefix("hello world", "hello ")
		if s != "world" {
			panic("TrimPrefix failed")
		}
	}
	{
		// TrimSuffix.
		s := stringslite.TrimSuffix("hello world", " world")
		if s != "hello" {
			panic("TrimSuffix failed")
		}
	}
}
