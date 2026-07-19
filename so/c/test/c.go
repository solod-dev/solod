package main

import (
	"solod.dev/so/c"
	"solod.dev/so/testing"
)

//so:embed main.h
var main_h string

//so:extern
func isalpha(ch int32) bool {
	_ = ch
	return false
}

//so:extern
func get_cstring(s string) *c.ConstChar {
	_ = s
	return nil
}

func TestAssert(t *testing.T) {
	_ = t
	a, b := 11, 11
	c.Assert(a == b, "a != b")
}

func TestString(t *testing.T) {
	cstr := get_cstring("Hello, C!")
	str := c.String(cstr)
	if str != "Hello, C!" {
		t.Error("String = " + str + ", want Hello, C!")
	}
}

func TestExtern(t *testing.T) {
	if !isalpha('a') {
		t.Error("isalpha('a') = false")
	}
}

func TestVal(t *testing.T) {
	nan := c.Val[float64]("NAN")
	if nan == nan {
		t.Error("NAN == NAN")
	}
	x := c.Val[float64]("sqrt(49)")
	if x != 7 {
		t.Error("sqrt(49) != 7")
	}
}

func TestRaw(t *testing.T) {
	var b int
	c.Raw(`
	int a = 7;
	b = a * a;
	b = sqrt(b);
	`)
	if b != 7 {
		t.Error("Raw block: b != 7")
	}
}

func TestCString(t *testing.T) {
	_ = t
	s := "hello"
	p := (*c.ConstChar)(c.CString(s))
	_ = p
}

func TestNumericTypes(t *testing.T) {
	var i c.Int = 42
	var u c.UInt = c.UInt(i)
	var l c.Long = c.Long(u)
	var ul c.ULong = c.ULong(l)
	var ll c.LongLong = c.LongLong(ul)
	var ull c.ULongLong = c.ULongLong(ll)
	if ull != 42 {
		t.Error("ull != 42")
	}
}
