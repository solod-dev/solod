// Copyright 2009 The Go Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package bytes_test

import (
	"testing"

	. "solod.dev/so/bytes"
	"solod.dev/so/unicode"
	"solod.dev/so/unicode/utf8"
)

var lowerTests = []StringTest{
	{"", []byte("")},
	{"abc", []byte("abc")},
	{"AbC123", []byte("abc123")},
	{"azAZ09_", []byte("azaz09_")},
	{"longStrinGwitHmixofsmaLLandcAps", []byte("longstringwithmixofsmallandcaps")},
	{"LONG\u2C6FSTRING\u2C6FWITH\u2C6FNONASCII\u2C6FCHARS", []byte("long\u0250string\u0250with\u0250nonascii\u0250chars")},
	{"\u2C6D\u2C6D\u2C6D\u2C6D\u2C6D", []byte("\u0251\u0251\u0251\u0251\u0251")}, // shrinks one byte per char
	{"A\u0080\U0010FFFF", []byte("a\u0080\U0010FFFF")},                           // test utf8.RuneSelf and utf8.MaxRune
}

func TestToLower(t *testing.T) {
	runStringTests(t, ToLower, "ToLower", lowerTests)
}

var upperTests = []StringTest{
	{"", []byte("")},
	{"ONLYUPPER", []byte("ONLYUPPER")},
	{"abc", []byte("ABC")},
	{"AbC123", []byte("ABC123")},
	{"azAZ09_", []byte("AZAZ09_")},
	{"longStrinGwitHmixofsmaLLandcAps", []byte("LONGSTRINGWITHMIXOFSMALLANDCAPS")},
	{"long\u0250string\u0250with\u0250nonascii\u2C6Fchars", []byte("LONG\u2C6FSTRING\u2C6FWITH\u2C6FNONASCII\u2C6FCHARS")},
	{"\u0250\u0250\u0250\u0250\u0250", []byte("\u2C6F\u2C6F\u2C6F\u2C6F\u2C6F")}, // grows one byte per char
	{"a\u0080\U0010FFFF", []byte("A\u0080\U0010FFFF")},                           // test utf8.RuneSelf and utf8.MaxRune
}

func TestToUpper(t *testing.T) {
	runStringTests(t, ToUpper, "ToUpper", upperTests)
}

func TestMap(t *testing.T) {
	// Run a couple of awful growth/shrinkage tests
	a := tenRunes('a')

	// 1. Grow. This triggers two reallocations in Map.
	maxRune := func(r rune) rune { return unicode.MaxRune }
	m := Map(nil, maxRune, []byte(a))
	expect := tenRunes(unicode.MaxRune)
	if string(m) != expect {
		t.Errorf("growing: expected %q got %q", expect, m)
	}

	// 2. Shrink
	minRune := func(r rune) rune { return 'a' }
	m = Map(nil, minRune, []byte(tenRunes(unicode.MaxRune)))
	expect = a
	if string(m) != expect {
		t.Errorf("shrinking: expected %q got %q", expect, m)
	}

	// 3. Rot13
	m = Map(nil, rot13, []byte("a to zed"))
	expect = "n gb mrq"
	if string(m) != expect {
		t.Errorf("rot13: expected %q got %q", expect, m)
	}

	// 4. Rot13^2
	m = Map(nil, rot13, Map(nil, rot13, []byte("a to zed")))
	expect = "a to zed"
	if string(m) != expect {
		t.Errorf("rot13: expected %q got %q", expect, m)
	}

	// 5. Drop
	dropNotUpper := func(r rune) rune {
		if unicode.Is(unicode.Upper, r) {
			return r
		}
		return -1
	}
	m = Map(nil, dropNotUpper, []byte("Hello, WORLD"))
	expect = "HWORLD"
	if string(m) != expect {
		t.Errorf("drop: expected %q got %q", expect, m)
	}

	// 6. Invalid rune
	invalidRune := func(r rune) rune {
		return utf8.MaxRune + 1
	}
	m = Map(nil, invalidRune, []byte("x"))
	expect = "\uFFFD"
	if string(m) != expect {
		t.Errorf("invalidRune: expected %q got %q", expect, m)
	}
}

// User-defined self-inverse mapping function
func rot13(r rune) rune {
	const step = 13
	if r >= 'a' && r <= 'z' {
		return ((r - 'a' + step) % 26) + 'a'
	}
	if r >= 'A' && r <= 'Z' {
		return ((r - 'A' + step) % 26) + 'A'
	}
	return r
}

func tenRunes(r rune) string {
	runes := make([]rune, 10)
	for i := range runes {
		runes[i] = r
	}
	return string(runes)
}
