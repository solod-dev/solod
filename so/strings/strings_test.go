// Copyright 2009 The Go Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package strings_test

import (
	"testing"
	"unsafe"

	. "solod.dev/so/strings"
	"solod.dev/so/unicode/utf8"
)

var abcd = "abcd"
var faces = "☺☻☹"
var commas = "1,2,3,4"
var dots = "1....2....3....4"

const space = "\t\v\r\f\n\u0085\u00a0\u2000\u3000"

// Test case for any function which accepts and returns a single string.
type StringTest struct {
	in, out string
}

func TestClone(t *testing.T) {
	var emptyString string

	var cloneTests = []string{
		"",
		Clone(nil, ""),
		Repeat(nil, "a", 42)[:0],
		"short",
		Repeat(nil, "a", 42),
	}
	for _, input := range cloneTests {
		clone := Clone(nil, input)
		if clone != input {
			t.Errorf("Clone(%q) = %q; want %q", input, clone, input)
		}

		if len(input) != 0 && unsafe.StringData(clone) == unsafe.StringData(input) {
			t.Errorf("Clone(%q) return value should not reference inputs backing memory.", input)
		}

		if len(input) == 0 && unsafe.StringData(clone) != unsafe.StringData(emptyString) {
			t.Errorf("Clone(%#v) return value should be equal to empty string.", unsafe.StringData(input))
		}
	}
}

func TestCount(t *testing.T) {
	var tests = []struct {
		s, sep string
		num    int
	}{
		{"", "", 1},
		{"", "notempty", 0},
		{"notempty", "", 9},
		{"smaller", "not smaller", 0},
		{"12345678987654321", "6", 2},
		{"611161116", "6", 3},
		{"notequal", "NotEqual", 0},
		{"equal", "equal", 1},
		{"abc1231231123q", "123", 3},
		{"11111", "11", 2},
	}

	for _, tt := range tests {
		if num := Count(tt.s, tt.sep); num != tt.num {
			t.Errorf("Count(%q, %q) = %d, want %d", tt.s, tt.sep, num, tt.num)
		}
	}
}

func TestCut(t *testing.T) {
	var tests = []struct {
		s, sep        string
		before, after string
		found         bool
	}{
		{"abc", "b", "a", "c", true},
		{"abc", "a", "", "bc", true},
		{"abc", "c", "ab", "", true},
		{"abc", "abc", "", "", true},
		{"abc", "", "", "abc", true},
		{"abc", "d", "abc", "", false},
		{"", "d", "", "", false},
		{"", "", "", "", true},
	}

	for _, tt := range tests {
		if before, after := Cut(tt.s, tt.sep); before != tt.before || after != tt.after {
			t.Errorf("Cut(%q, %q) = %q, %q, want %q, %q", tt.s, tt.sep, before, after, tt.before, tt.after)
		}
	}
}

func TestCutPrefix(t *testing.T) {
	var tests = []struct {
		s, sep string
		after  string
		found  bool
	}{
		{"abc", "a", "bc", true},
		{"abc", "abc", "", true},
		{"abc", "", "abc", true},
		{"abc", "d", "abc", false},
		{"", "d", "", false},
		{"", "", "", true},
	}

	for _, tt := range tests {
		if after, found := CutPrefix(tt.s, tt.sep); after != tt.after || found != tt.found {
			t.Errorf("CutPrefix(%q, %q) = %q, %v, want %q, %v", tt.s, tt.sep, after, found, tt.after, tt.found)
		}
	}
}

func TestCutSuffix(t *testing.T) {
	var tests = []struct {
		s, sep string
		before string
		found  bool
	}{
		{"abc", "bc", "a", true},
		{"abc", "abc", "", true},
		{"abc", "", "abc", true},
		{"abc", "d", "abc", false},
		{"", "d", "", false},
		{"", "", "", true},
	}
	for _, tt := range tests {
		if before, found := CutSuffix(tt.s, tt.sep); before != tt.before || found != tt.found {
			t.Errorf("CutSuffix(%q, %q) = %q, %v, want %q, %v", tt.s, tt.sep, before, found, tt.before, tt.found)
		}
	}
}

func TestReplace(t *testing.T) {
	var tests = []struct {
		in       string
		old, new string
		n        int
		out      string
	}{
		{"hello", "l", "L", 0, "hello"},
		{"hello", "l", "L", -1, "heLLo"},
		{"hello", "x", "X", -1, "hello"},
		{"", "x", "X", -1, ""},
		{"radar", "r", "<r>", -1, "<r>ada<r>"},
		{"", "", "<>", -1, "<>"},
		{"banana", "a", "<>", -1, "b<>n<>n<>"},
		{"banana", "a", "<>", 1, "b<>nana"},
		{"banana", "a", "<>", 1000, "b<>n<>n<>"},
		{"banana", "an", "<>", -1, "b<><>a"},
		{"banana", "ana", "<>", -1, "b<>na"},
		{"banana", "", "<>", -1, "<>b<>a<>n<>a<>n<>a<>"},
		{"banana", "", "<>", 10, "<>b<>a<>n<>a<>n<>a<>"},
		{"banana", "", "<>", 6, "<>b<>a<>n<>a<>n<>a"},
		{"banana", "", "<>", 5, "<>b<>a<>n<>a<>na"},
		{"banana", "", "<>", 1, "<>banana"},
		{"banana", "a", "a", -1, "banana"},
		{"banana", "a", "a", 1, "banana"},
		{"☺☻☹", "", "<>", -1, "<>☺<>☻<>☹<>"},
	}

	for _, tt := range tests {
		if s := Replace(nil, tt.in, tt.old, tt.new, tt.n); s != tt.out {
			t.Errorf("Replace(%q, %q, %q, %d) = %q, want %q", tt.in, tt.old, tt.new, tt.n, s, tt.out)
		}
		if tt.n == -1 {
			s := ReplaceAll(nil, tt.in, tt.old, tt.new)
			if s != tt.out {
				t.Errorf("ReplaceAll(%q, %q, %q) = %q, want %q", tt.in, tt.old, tt.new, s, tt.out)
			}
		}
	}
}

func TestRunes(t *testing.T) {
	runesEqual := func(a, b []rune) bool {
		if len(a) != len(b) {
			return false
		}
		for i, r := range a {
			if r != b[i] {
				return false
			}
		}
		return true
	}

	var tests = []struct {
		in    string
		out   []rune
		lossy bool
	}{
		{"", []rune{}, false},
		{" ", []rune{32}, false},
		{"ABC", []rune{65, 66, 67}, false},
		{"abc", []rune{97, 98, 99}, false},
		{"\u65e5\u672c\u8a9e", []rune{26085, 26412, 35486}, false},
		{"ab\x80c", []rune{97, 98, 0xFFFD, 99}, true},
		{"ab\xc0c", []rune{97, 98, 0xFFFD, 99}, true},
	}

	for _, tt := range tests {
		a := []rune(tt.in)
		if !runesEqual(a, tt.out) {
			t.Errorf("[]rune(%q) = %v; want %v", tt.in, a, tt.out)
			continue
		}
		if !tt.lossy {
			// can only test reassembly if we didn't lose information
			s := string(a)
			if s != tt.in {
				t.Errorf("string([]rune(%q)) = %x; want %x", tt.in, s, tt.in)
			}
		}
	}
}

// Execute f on each test case.  funcName should be the name of f; it's used
// in failure reports.
func runStringTests(t *testing.T, f func(string) string, funcName string, testCases []StringTest) {
	for _, tc := range testCases {
		actual := f(tc.in)
		if actual != tc.out {
			t.Errorf("%s(%q) = %q; want %q", funcName, tc.in, actual, tc.out)
		}
	}
}

func equal(m string, s1, s2 string, t *testing.T) bool {
	if s1 == s2 {
		return true
	}
	e1 := Split(nil, s1, "")
	e2 := Split(nil, s2, "")
	for i, c1 := range e1 {
		if i >= len(e2) {
			break
		}
		r1, _ := utf8.DecodeRuneInString(c1)
		r2, _ := utf8.DecodeRuneInString(e2[i])
		if r1 != r2 {
			t.Errorf("%s diff at %d: U+%04X U+%04X", m, i, r1, r2)
		}
	}
	return false
}
