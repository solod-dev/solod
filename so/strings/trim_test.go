// Copyright 2009 The Go Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package strings_test

import (
	"testing"

	. "solod.dev/so/strings"
	"solod.dev/so/unicode"
	"solod.dev/so/unicode/utf8"
)

var trimTests = []struct {
	f            string
	in, arg, out string
}{
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

func TestTrim(t *testing.T) {
	for _, tc := range trimTests {
		name := tc.f
		var f func(string, string) string
		switch name {
		case "Trim":
			f = Trim
		case "TrimLeft":
			f = TrimLeft
		case "TrimRight":
			f = TrimRight
		case "TrimPrefix":
			f = TrimPrefix
		case "TrimSuffix":
			f = TrimSuffix
		default:
			t.Errorf("Undefined trim function %s", name)
		}
		actual := f(tc.in, tc.arg)
		if actual != tc.out {
			t.Errorf("%s(%q, %q) = %q; want %q", name, tc.in, tc.arg, actual, tc.out)
		}
	}
}

var trimSpaceTests = []StringTest{
	{"", ""},
	{"abc", "abc"},
	{space + "abc" + space, "abc"},
	{" ", ""},
	{" \t\r\n \t\t\r\r\n\n ", ""},
	{" \t\r\n x\t\t\r\r\n\n ", "x"},
	{" \u2000\t\r\n x\t\t\r\r\ny\n \u3000", "x\t\t\r\r\ny"},
	{"1 \t\r\n2", "1 \t\r\n2"},
	{" x\x80", "x\x80"},
	{" x\xc0", "x\xc0"},
	{"x \xc0\xc0 ", "x \xc0\xc0"},
	{"x \xc0", "x \xc0"},
	{"x \xc0 ", "x \xc0"},
	{"x \xc0\xc0 ", "x \xc0\xc0"},
	{"x ☺\xc0\xc0 ", "x ☺\xc0\xc0"},
	{"x ☺ ", "x ☺"},
}

func TestTrimSpace(t *testing.T) { runStringTests(t, TrimSpace, "TrimSpace", trimSpaceTests) }

type predicate struct {
	f    func(rune) bool
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

func not(p predicate) predicate {
	return predicate{
		func(r rune) bool {
			return !p.f(r)
		},
		"not " + p.name,
	}
}

var trimFuncTests = []struct {
	f        predicate
	in       string
	trimOut  string
	leftOut  string
	rightOut string
}{
	{isSpace, space + " hello " + space,
		"hello",
		"hello " + space,
		space + " hello"},
	{isDigit, "\u0e50\u0e5212hello34\u0e50\u0e51",
		"hello",
		"hello34\u0e50\u0e51",
		"\u0e50\u0e5212hello"},
	{isUpper, "\u2C6F\u2C6F\u2C6F\u2C6FABCDhelloEF\u2C6F\u2C6FGH\u2C6F\u2C6F",
		"hello",
		"helloEF\u2C6F\u2C6FGH\u2C6F\u2C6F",
		"\u2C6F\u2C6F\u2C6F\u2C6FABCDhello"},
	{not(isSpace), "hello" + space + "hello",
		space,
		space + "hello",
		"hello" + space},
	{not(isDigit), "hello\u0e50\u0e521234\u0e50\u0e51helo",
		"\u0e50\u0e521234\u0e50\u0e51",
		"\u0e50\u0e521234\u0e50\u0e51helo",
		"hello\u0e50\u0e521234\u0e50\u0e51"},
	{isValidRune, "ab\xc0a\xc0cd",
		"\xc0a\xc0",
		"\xc0a\xc0cd",
		"ab\xc0a\xc0"},
	{not(isValidRune), "\xc0a\xc0",
		"a",
		"a\xc0",
		"\xc0a"},
	{isSpace, "",
		"",
		"",
		""},
	{isSpace, " ",
		"",
		"",
		""},
}

func TestTrimFunc(t *testing.T) {
	for _, tc := range trimFuncTests {
		trimmers := []struct {
			name string
			trim func(s string, f RunePredicate) string
			out  string
		}{
			{"TrimFunc", TrimFunc, tc.trimOut},
		}
		for _, trimmer := range trimmers {
			actual := trimmer.trim(tc.in, tc.f.f)
			if actual != trimmer.out {
				t.Errorf("%s(%q, %q) = %q; want %q", trimmer.name, tc.in, tc.f.name, actual, trimmer.out)
			}
		}
	}
}
