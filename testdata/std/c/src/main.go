package main

import "solod.dev/so/c"

//so:embed main.h
var main_h string

//so:extern
func isalpha(ch int32) bool

//so:extern
func get_cstring(s string) *c.ConstChar

func main() {
	{
		// Return `const char*` from C.
		cstr := get_cstring("Hello, C!")
		str := c.String(cstr)
		if str != "Hello, C!" {
			panic("unexpected string: " + str)
		}
	}
	{
		// Use header included via so:include.c
		if !isalpha('a') {
			panic("isalpha failed")
		}
	}
	{
		// Typed C expression.
		nan := c.Val[float64]("NAN")
		if nan == nan {
			panic("nan == nan")
		}
		x := c.Val[float64]("sqrt(49)")
		if x != 7 {
			panic("x != 7")
		}
	}
	{
		// Raw C block.
		var b int
		c.Raw(`
		int a = 7;
		b = a * a;
		b = sqrt(b);
		`)
		if b != 7 {
			panic("b != 7")
		}
	}
}
