package main

import (
	"solod.dev/so/mem"
	"solod.dev/so/strings"
)

func main() {
	{
		// Clone.
		s := "hello"
		c := strings.Clone(nil, s)
		if c != s {
			panic("Clone failed")
		}
		mem.FreeString(nil, c)
	}
	{
		// Compare.
		if strings.Compare("abc", "abb") <= 0 {
			panic("Compare failed")
		}
		if strings.Compare("abc", "abd") >= 0 {
			panic("Compare failed")
		}
		if strings.Compare("abc", "abc") != 0 {
			panic("Compare failed")
		}
	}
	{
		// Count.
		n := strings.Count("hello world", "o")
		if n != 2 {
			panic("Count failed")
		}
		n = strings.Count("hello world", "")
		if n != 12 {
			panic("Count failed")
		}
	}
	{
		// Cut.
		before, after := strings.Cut("hello world", " ")
		if before != "hello" || after != "world" {
			panic("Cut failed")
		}
	}
	{
		// CutPrefix and CutSuffix.
		src := "hello world"
		if s, ok := strings.CutPrefix(src, "hello"); !ok || s != " world" {
			panic("CutPrefix failed")
		}
		if s, ok := strings.CutSuffix(src, "world"); !ok || s != "hello " {
			panic("CutSuffix failed")
		}
	}
	{
		// Index and IndexAny.
		idx := strings.Index("hello world", "o")
		if idx != 4 {
			panic("Index failed")
		}
		idx = strings.IndexAny("hello world", "ow")
		if idx != 4 {
			panic("IndexAny failed")
		}
	}
	{
		// Repeat.
		r := strings.Repeat(nil, "abc", 3)
		if r != "abcabcabc" {
			panic("Repeat failed")
		}
		mem.FreeString(nil, r)
	}
	{
		// Replace and ReplaceAll.
		s := "hello world"
		r := strings.Replace(nil, s, "o", "0", 1)
		if r != "hell0 world" {
			panic("Replace failed")
		}
		mem.FreeString(nil, r)
		r = strings.ReplaceAll(nil, s, "o", "0")
		if r != "hell0 w0rld" {
			panic("ReplaceAll failed")
		}
		mem.FreeString(nil, r)
	}
	{
		// Split and Join.
		s := "a,b,c"
		parts := strings.Split(nil, s, ",")
		if len(parts) != 3 || parts[0] != "a" || parts[1] != "b" || parts[2] != "c" {
			panic("Split failed")
		}
		j := strings.Join(nil, parts, ",")
		if j != s {
			panic("Join failed")
		}
		mem.FreeString(nil, j)
		mem.FreeSlice(nil, parts)
	}
	{
		// ToUpper and ToLower.
		s := "Hello, 世界!"
		u := strings.ToUpper(nil, s)
		if u != "HELLO, 世界!" {
			panic("ToUpper failed")
		}
		mem.FreeString(nil, u)
		l := strings.ToLower(nil, s)
		if l != "hello, 世界!" {
			panic("ToLower failed")
		}
		mem.FreeString(nil, l)
	}
	{
		// Trim and TrimSpace.
		s := "  hello world  "
		t := strings.TrimSpace(s)
		if t != "hello world" {
			panic("TrimSpace failed")
		}
		t = strings.Trim(s, " dh")
		if t != "ello worl" {
			panic("Trim failed")
		}
	}
	{
		// Builder.
		var b strings.Builder
		b.WriteString("Hello")
		b.WriteByte(',')
		b.WriteRune(' ')
		b.WriteString("world")
		s := b.String()
		if s != "Hello, world" {
			panic("Builder failed")
		}
		b.Free()
	}
	{
		// Reader.
		r := strings.NewReader("hello world")
		buf := make([]byte, 5)
		n, err := r.Read(buf)
		if err != nil || n != 5 || string(buf) != "hello" {
			panic("Reader failed")
		}
	}
}
