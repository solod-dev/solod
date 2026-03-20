// Copyright 2009 The Go Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package bytes_test

import (
	"fmt"
	"testing"

	. "solod.dev/so/bytes"
	"solod.dev/so/mem"
	"solod.dev/so/unicode"
	"solod.dev/so/unicode/utf8"
)

type predicate struct {
	f    func(r rune) bool
	name string
}

var isSpace = predicate{unicode.IsSpace, "IsSpace"}
var isDigit = predicate{unicode.IsDigit, "IsDigit"}
var isUpper = predicate{unicode.IsUpper, "IsUpper"}
var isValidRune = predicate{
	func(r rune) bool {
		return r != utf8.RuneError
	},
	"IsValidRune",
}

type TrimTest struct {
	f            string
	in, arg, out string
}

var trimTests = []TrimTest{
	{"Trim", "abba", "a", "bb"},
	{"Trim", "abba", "ab", ""},
	{"TrimLeft", "abba", "ab", ""},
	{"TrimRight", "abba", "ab", ""},
	{"TrimLeft", "abba", "a", "bba"},
	{"TrimLeft", "abba", "b", "abba"},
	{"TrimRight", "abba", "a", "abb"},
	{"TrimRight", "abba", "b", "abba"},
	{"Trim", "<tag>", "<>", "tag"},
	{"Trim", "* listitem", " *", "listitem"},
	{"Trim", `"quote"`, `"`, "quote"},
	{"Trim", "\u2C6F\u2C6F\u0250\u0250\u2C6F\u2C6F", "\u2C6F", "\u0250\u0250"},
	{"Trim", "\x80test\xff", "\xff", "test"},
	{"Trim", " Ġ ", " ", "Ġ"},
	{"Trim", " Ġİ0", "0 ", "Ġİ"},
	//empty string tests
	{"Trim", "abba", "", "abba"},
	{"Trim", "", "123", ""},
	{"Trim", "", "", ""},
	{"TrimLeft", "abba", "", "abba"},
	{"TrimLeft", "", "123", ""},
	{"TrimLeft", "", "", ""},
	{"TrimRight", "abba", "", "abba"},
	{"TrimRight", "", "123", ""},
	{"TrimRight", "", "", ""},
	{"TrimRight", "☺\xc0", "☺", "☺\xc0"},
	{"TrimPrefix", "aabb", "a", "abb"},
	{"TrimPrefix", "aabb", "b", "aabb"},
	{"TrimSuffix", "aabb", "a", "aabb"},
	{"TrimSuffix", "aabb", "b", "aab"},
}

type TrimNilTest struct {
	f   string
	in  []byte
	arg string
	out []byte
}

var trimNilTests = []TrimNilTest{
	{"Trim", nil, "", nil},
	{"Trim", []byte{}, "", []byte{}},
	{"Trim", []byte{'a'}, "a", []byte{}},
	{"Trim", []byte{'a', 'a'}, "a", []byte{}},
	{"Trim", []byte{'a'}, "ab", []byte{}},
	{"Trim", []byte{'a', 'b'}, "ab", []byte{}},
	{"Trim", []byte("☺"), "☺", []byte{}},
	{"TrimLeft", nil, "", nil},
	{"TrimLeft", []byte{}, "", []byte{}},
	{"TrimLeft", []byte{'a'}, "a", []byte{}},
	{"TrimLeft", []byte{'a', 'a'}, "a", []byte{}},
	{"TrimLeft", []byte{'a'}, "ab", []byte{}},
	{"TrimLeft", []byte{'a', 'b'}, "ab", []byte{}},
	{"TrimLeft", []byte("☺"), "☺", []byte{}},
	{"TrimRight", nil, "", nil},
	{"TrimRight", []byte{}, "", []byte{}},
	{"TrimRight", []byte{'a'}, "a", []byte{}},
	{"TrimRight", []byte{'a', 'a'}, "a", []byte{}},
	{"TrimRight", []byte{'a'}, "ab", []byte{}},
	{"TrimRight", []byte{'a', 'b'}, "ab", []byte{}},
	{"TrimRight", []byte("☺"), "☺", []byte{}},
	{"TrimPrefix", nil, "", nil},
	{"TrimPrefix", []byte{}, "", []byte{}},
	{"TrimPrefix", []byte{'a'}, "a", []byte{}},
	{"TrimPrefix", []byte("☺"), "☺", []byte{}},
	{"TrimSuffix", nil, "", nil},
	{"TrimSuffix", []byte{}, "", []byte{}},
	{"TrimSuffix", []byte{'a'}, "a", []byte{}},
	{"TrimSuffix", []byte("☺"), "☺", []byte{}},
}

func TestTrim(t *testing.T) {
	toFn := func(name string) (func([]byte, string) []byte, func([]byte, []byte) []byte) {
		switch name {
		case "Trim":
			return Trim, nil
		case "TrimLeft":
			return TrimLeft, nil
		case "TrimRight":
			return TrimRight, nil
		case "TrimPrefix":
			return nil, TrimPrefix
		case "TrimSuffix":
			return nil, TrimSuffix
		default:
			t.Errorf("Undefined trim function %s", name)
			return nil, nil
		}
	}

	for _, tc := range trimTests {
		name := tc.f
		f, fb := toFn(name)
		if f == nil && fb == nil {
			continue
		}
		var actual string
		if f != nil {
			actual = string(f([]byte(tc.in), tc.arg))
		} else {
			actual = string(fb([]byte(tc.in), []byte(tc.arg)))
		}
		if actual != tc.out {
			t.Errorf("%s(%q, %q) = %q; want %q", name, tc.in, tc.arg, actual, tc.out)
		}
	}

	for _, tc := range trimNilTests {
		name := tc.f
		f, fb := toFn(name)
		if f == nil && fb == nil {
			continue
		}
		var actual []byte
		if f != nil {
			actual = f(tc.in, tc.arg)
		} else {
			actual = fb(tc.in, []byte(tc.arg))
		}
		report := func(s []byte) string {
			if s == nil {
				return "nil"
			} else {
				return fmt.Sprintf("%q", s)
			}
		}
		if len(actual) != 0 {
			t.Errorf("%s(%s, %q) returned non-empty value", name, report(tc.in), tc.arg)
		} else {
			actualNil := actual == nil
			outNil := tc.out == nil
			if actualNil != outNil {
				t.Errorf("%s(%s, %q) got nil %t; want nil %t", name, report(tc.in), tc.arg, actualNil, outNil)
			}
		}
	}
}

func TestTrimSpace(t *testing.T) {
	var tests = []StringTest{
		{"", []byte("")},
		{"  a", []byte("a")},
		{"b  ", []byte("b")},
		{"abc", []byte("abc")},
		{space + "abc" + space, []byte("abc")},
		{" ", []byte("")},
		{"\u3000 ", []byte("")},
		{" \u3000", []byte("")},
		{" \t\r\n \t\t\r\r\n\n ", []byte("")},
		{" \t\r\n x\t\t\r\r\n\n ", []byte("x")},
		{" \u2000\t\r\n x\t\t\r\r\ny\n \u3000", []byte("x\t\t\r\r\ny")},
		{"1 \t\r\n2", []byte("1 \t\r\n2")},
		{" x\x80", []byte("x\x80")},
		{" x\xc0", []byte("x\xc0")},
		{"x \xc0\xc0 ", []byte("x \xc0\xc0")},
		{"x \xc0", []byte("x \xc0")},
		{"x \xc0 ", []byte("x \xc0")},
		{"x \xc0\xc0 ", []byte("x \xc0\xc0")},
		{"x ☺\xc0\xc0 ", []byte("x ☺\xc0\xc0")},
		{"x ☺ ", []byte("x ☺")},
	}
	trimSpace := func(a mem.Allocator, b []byte) []byte {
		return TrimSpace(b)
	}
	runStringTests(t, trimSpace, "TrimSpace", tests)
}

func TestTrimFunc(t *testing.T) {
	not := func(p predicate) predicate {
		return predicate{
			func(r rune) bool {
				return !p.f(r)
			},
			"not " + p.name,
		}
	}

	type TrimFuncTest struct {
		f        predicate
		in       string
		trimOut  []byte
		leftOut  []byte
		rightOut []byte
	}

	var tests = []TrimFuncTest{
		{isSpace, space + " hello " + space,
			[]byte("hello"),
			[]byte("hello " + space),
			[]byte(space + " hello")},
		{isDigit, "\u0e50\u0e5212hello34\u0e50\u0e51",
			[]byte("hello"),
			[]byte("hello34\u0e50\u0e51"),
			[]byte("\u0e50\u0e5212hello")},
		{isUpper, "\u2C6F\u2C6F\u2C6F\u2C6FABCDhelloEF\u2C6F\u2C6FGH\u2C6F\u2C6F",
			[]byte("hello"),
			[]byte("helloEF\u2C6F\u2C6FGH\u2C6F\u2C6F"),
			[]byte("\u2C6F\u2C6F\u2C6F\u2C6FABCDhello")},
		{not(isSpace), "hello" + space + "hello",
			[]byte(space),
			[]byte(space + "hello"),
			[]byte("hello" + space)},
		{not(isDigit), "hello\u0e50\u0e521234\u0e50\u0e51helo",
			[]byte("\u0e50\u0e521234\u0e50\u0e51"),
			[]byte("\u0e50\u0e521234\u0e50\u0e51helo"),
			[]byte("hello\u0e50\u0e521234\u0e50\u0e51")},
		{isValidRune, "ab\xc0a\xc0cd",
			[]byte("\xc0a\xc0"),
			[]byte("\xc0a\xc0cd"),
			[]byte("ab\xc0a\xc0")},
		{not(isValidRune), "\xc0a\xc0",
			[]byte("a"),
			[]byte("a\xc0"),
			[]byte("\xc0a")},
		// The nils returned by TrimLeftFunc are odd behavior, but we need
		// to preserve backwards compatibility.
		{isSpace, "",
			[]byte(""),
			[]byte(""),
			[]byte("")},
		{isSpace, " ",
			[]byte(""),
			[]byte(""),
			[]byte("")},
	}

	for _, tc := range tests {
		trimmers := []struct {
			name string
			trim func(s []byte, f RunePredicate) []byte
			out  []byte
		}{
			{"TrimFunc", TrimFunc, tc.trimOut},
		}
		for _, trimmer := range trimmers {
			actual := trimmer.trim([]byte(tc.in), tc.f.f)
			if actual == nil && trimmer.out != nil {
				t.Errorf("%s(%q, %q) = nil; want %q", trimmer.name, tc.in, tc.f.name, trimmer.out)
			}
			if actual != nil && trimmer.out == nil {
				t.Errorf("%s(%q, %q) = %q; want nil", trimmer.name, tc.in, tc.f.name, actual)
			}
			if !Equal(actual, trimmer.out) {
				t.Errorf("%s(%q, %q) = %q; want %q", trimmer.name, tc.in, tc.f.name, actual, trimmer.out)
			}
		}
	}
}
