package main

import (
	"github.com/nalgeon/solod/so/fmt"
	"github.com/nalgeon/solod/so/strings"
)

func main() {
	{
		// Print.
		n, err := fmt.Print("hello", "world")
		if err != nil {
			panic("Print failed")
		}
		if n != 11 {
			panic("Print: wrong count")
		}
	}
	{
		// Println.
		n, err := fmt.Println("hello", "world")
		if err != nil {
			panic("Println failed")
		}
		if n != 12 {
			panic("Println: wrong count")
		}
	}
	{
		// Printf.
		s := "world"
		d := 42
		n, err := fmt.Printf("s = %s, d = %d\n", s, d)
		if err != nil {
			panic("Printf failed")
		}
		if n != 18 {
			panic("Printf: wrong count")
		}
	}
	{
		// Fprintf.
		var sb strings.Builder
		s := "world"
		n, err := fmt.Fprintf(&sb, "hello %s", s)
		if err != nil {
			panic("Fprintf failed")
		}
		if n != 11 {
			panic("Fprintf: wrong count")
		}
		if sb.String() != "hello world" {
			panic("Fprintf: wrong output")
		}
		sb.Free()
	}
	{
		// Sscanf.
		var a int32
		var b int32
		n, err := fmt.Sscanf("42 7", "%d %d", &a, &b)
		if err != nil {
			panic("Sscanf failed")
		}
		if n != 2 {
			panic("Sscanf: wrong count")
		}
		if a != 42 || b != 7 {
			panic("Sscanf: wrong values")
		}
	}
	{
		// Fscanf.
		r := strings.NewReader(nil, "100 200")
		var a int32
		var b int32
		n, err := fmt.Fscanf(&r, "%d %d", &a, &b)
		if err != nil {
			panic("Fscanf failed")
		}
		if n != 2 {
			panic("Fscanf: wrong count")
		}
		if a != 100 || b != 200 {
			panic("Fscanf: wrong values")
		}
	}
}
