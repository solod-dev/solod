// Copyright 2009 The Go Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

// Package strings implements simple functions to manipulate UTF-8 encoded strings.
package strings

import (
	"solod.dev/so/bytealg"
	"solod.dev/so/mem"
	"solod.dev/so/stringslite"
	"solod.dev/so/unicode/utf8"
)

// Clone returns a fresh copy of s.
//
// It guarantees to make a copy of s into a new allocation,
// which can be important when retaining only a small substring
// of a much larger string. Using Clone can help such programs
// use less memory. Of course, since using Clone makes a copy,
// overuse of Clone can make programs use more memory.
//
// Clone should typically be used only rarely, and only when
// profiling indicates that it is needed.
//
// The returned string is allocated; the caller owns it.
func Clone(a mem.Allocator, s string) string {
	return stringslite.Clone(a, s)
}

// Compare returns an integer comparing two strings lexicographically.
// The result will be 0 if a == b, -1 if a < b, and +1 if a > b.
//
// Use Compare when you need to perform a three-way comparison (with
// [slices.SortFunc], for example). It is usually clearer and always faster
// to use the built-in string comparison operators ==, <, >, and so on.
func Compare(a, b string) int {
	return bytealg.Compare([]byte(a), []byte(b))
}

// Count counts the number of non-overlapping instances of substr in s.
// If substr is an empty string, Count returns 1 + the number of Unicode code points in s.
func Count(s, substr string) int {
	// special case
	if len(substr) == 0 {
		return utf8.RuneCountInString(s) + 1
	}
	if len(substr) == 1 {
		return bytealg.CountString(s, substr[0])
	}
	n := 0
	for {
		i := Index(s, substr)
		if i == -1 {
			return n
		}
		n++
		s = s[i+len(substr):]
	}
}

// Cut slices s around the first instance of sep,
// returning the text before and after sep.
// If sep does not appear in s, cut returns s, "".
func Cut(s, sep string) (string, string) {
	return stringslite.Cut(s, sep)
}

// CutPrefix returns s without the provided leading prefix string
// and reports whether it found the prefix.
// If s doesn't start with prefix, CutPrefix returns s, false.
// If prefix is the empty string, CutPrefix returns s, true.
func CutPrefix(s, prefix string) (string, bool) {
	return stringslite.CutPrefix(s, prefix)
}

// CutSuffix returns s without the provided ending suffix string
// and reports whether it found the suffix.
// If s doesn't end with suffix, CutSuffix returns s, false.
// If suffix is the empty string, CutSuffix returns s, true.
func CutSuffix(s, suffix string) (string, bool) {
	return stringslite.CutSuffix(s, suffix)
}

// Join concatenates the elements of its first argument to create a single string.
// The separator string sep is placed between elements in the resulting string.
//
// The returned string is allocated; the caller owns it.
func Join(a mem.Allocator, elems []string, sep string) string {
	if len(elems) == 0 {
		return ""
	} else if len(elems) == 1 {
		return stringslite.Clone(a, elems[0])
	}

	var n int
	if len(sep) > 0 {
		if len(sep) >= maxInt/(len(elems)-1) {
			panic("strings: Join output length overflow")
		}
		n += len(sep) * (len(elems) - 1)
	}
	for _, elem := range elems {
		if len(elem) > maxInt-n {
			panic("strings: Join output length overflow")
		}
		n += len(elem)
	}

	b := Builder{a: a}
	b.Grow(n)
	b.WriteString(elems[0])
	for _, s := range elems[1:] {
		b.WriteString(sep)
		b.WriteString(s)
	}
	return b.String()
}

// HasPrefix reports whether the string s begins with prefix.
func HasPrefix(s, prefix string) bool {
	return stringslite.HasPrefix(s, prefix)
}

// HasSuffix reports whether the string s ends with suffix.
func HasSuffix(s, suffix string) bool {
	return stringslite.HasSuffix(s, suffix)
}

// Replace returns a copy of the string s with the first n
// non-overlapping instances of old replaced by new.
//
// If old is empty, it matches at the beginning of the string
// and after each UTF-8 sequence, yielding up to k+1 replacements
// for a k-rune string.
//
// If n < 0, there is no limit on the number of replacements.
//
// The returned string is allocated; the caller owns it.
func Replace(a mem.Allocator, s, old, new string, n int) string {
	if old == new || n == 0 {
		return stringslite.Clone(a, s)
	}

	// Compute number of replacements.
	if m := Count(s, old); m == 0 {
		return stringslite.Clone(a, s)
	} else if n < 0 || m < n {
		n = m
	}

	// Apply replacements to buffer.
	b := Builder{a: a}
	b.Grow(len(s) + n*(len(new)-len(old)))
	start := 0
	if len(old) > 0 {
		for range n {
			j := start + Index(s[start:], old)
			b.WriteString(s[start:j])
			b.WriteString(new)
			start = j + len(old)
		}
	} else { // len(old) == 0
		b.WriteString(new)
		for range n - 1 {
			_, wid := utf8.DecodeRuneInString(s[start:])
			j := start + wid
			b.WriteString(s[start:j])
			b.WriteString(new)
			start = j
		}
	}
	b.WriteString(s[start:])
	return b.String()
}

// ReplaceAll returns a copy of the string s with all
// non-overlapping instances of old replaced by new.
// If old is empty, it matches at the beginning of the string
// and after each UTF-8 sequence, yielding up to k+1 replacements
// for a k-rune string.
//
// The returned string is allocated; the caller owns it.
func ReplaceAll(a mem.Allocator, s, old, new string) string {
	return Replace(a, s, old, new, -1)
}
