// Copyright 2009 The Go Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package strings

import (
	"solod.dev/so/bytealg"
	"solod.dev/so/stringslite"
	"solod.dev/so/unicode/utf8"
)

// asciiSet is a 32-byte value, where each bit represents the presence of a
// given ASCII character in the set. The 128-bits of the lower 16 bytes,
// starting with the least-significant bit of the lowest word to the
// most-significant bit of the highest word, map to the full range of all
// 128 ASCII characters. The 128-bits of the upper 16 bytes will be zeroed,
// ensuring that any non-ASCII character will be reported as not in the set.
// This allocates a total of 32 bytes even though the upper half
// is unused to avoid bounds checks in asciiSet.contains.
type asciiSet struct {
	val [8]uint32
	ok  bool
}

// makeASCIISet creates a set of ASCII characters and reports whether all
// characters in chars are ASCII.
func makeASCIISet(chars string) asciiSet {
	var as asciiSet
	for i := 0; i < len(chars); i++ {
		c := chars[i]
		if c >= utf8.RuneSelf {
			return as
		}
		as.val[c/32] |= 1 << (c % 32)
	}
	as.ok = true
	return as
}

// contains reports whether c is inside the set.
func (as *asciiSet) contains(c byte) bool {
	return (as.val[c/32] & (1 << (c % 32))) != 0
}

// RunePredicate reports whether the rune satisfies a condition.
type RunePredicate func(rune) bool

// Contains reports whether substr is within s.
func Contains(s, substr string) bool {
	return Index(s, substr) >= 0
}

// ContainsAny reports whether any Unicode code points in chars are within s.
func ContainsAny(s, chars string) bool {
	return IndexAny(s, chars) >= 0
}

// ContainsRune reports whether the Unicode code point r is within s.
func ContainsRune(s string, r rune) bool {
	return IndexRune(s, r) >= 0
}

// ContainsFunc reports whether any Unicode code points r within s satisfy f(r).
func ContainsFunc(s string, f RunePredicate) bool {
	return IndexFunc(s, f) >= 0
}

// Index returns the index of the first instance of substr in s, or -1 if substr is not present in s.
func Index(s, substr string) int {
	return stringslite.Index(s, substr)
}

// LastIndex returns the index of the last instance of substr in s, or -1 if substr is not present in s.
func LastIndex(s, substr string) int {
	n := len(substr)
	if n == 0 {
		return len(s)
	} else if n == 1 {
		return bytealg.LastIndexByteString(s, substr[0])
	} else if n == len(s) {
		if substr == s {
			return 0
		}
		return -1
	} else if n > len(s) {
		return -1
	}
	// Rabin-Karp search from the end of the string
	hashss, pow := bytealg.HashStrRev([]byte(substr))
	last := len(s) - n
	var h uint32
	for i := len(s) - 1; i >= last; i-- {
		h = h*bytealg.PrimeRK + uint32(s[i])
	}
	if h == hashss && s[last:] == substr {
		return last
	}
	for i := last - 1; i >= 0; i-- {
		h *= bytealg.PrimeRK
		h += uint32(s[i])
		h -= pow * uint32(s[i+n])
		if h == hashss && s[i:i+n] == substr {
			return i
		}
	}
	return -1
}

// IndexByte returns the index of the first instance of c in s, or -1 if c is not present in s.
func IndexByte(s string, c byte) int {
	return stringslite.IndexByte(s, c)
}

// IndexRune returns the index of the first instance of the Unicode code point
// r, or -1 if rune is not present in s.
// If r is [utf8.RuneError], it returns the first instance of any
// invalid UTF-8 byte sequence.
func IndexRune(s string, r rune) int {
	if 0 <= r && r < utf8.RuneSelf {
		return IndexByte(s, byte(r))
	} else if r == utf8.RuneError {
		for i, r := range s {
			if r == utf8.RuneError {
				return i
			}
		}
		return -1
	} else if !utf8.ValidRune(r) {
		return -1
	}

	// Search for rune r using the last byte of its UTF-8 encoded form.
	// The distribution of the last byte is more uniform compared to the
	// first byte which has a 78% chance of being [240, 243, 244].
	rs := string(r)
	last := len(rs) - 1
	i := last
	fails := 0
	fallback := false
	for i < len(s) {
		if s[i] != rs[last] {
			o := IndexByte(s[i+1:], rs[last])
			if o < 0 {
				return -1
			}
			i += o + 1
		}
		// Step backwards comparing bytes.
		matched := true
		for j := 1; j < len(rs); j++ {
			if s[i-j] != rs[last-j] {
				matched = false
				break
			}
		}
		if matched {
			return i - last
		}
		fails++
		i++
		if fails >= (4+(i>>4)) && i < len(s) {
			fallback = true
			break
		}
	}
	if !fallback {
		return -1
	}

	c0 := rs[last]
	c1 := rs[last-1]
	for ; i < len(s); i++ {
		if s[i] == c0 && s[i-1] == c1 {
			found := true
			for k := 2; k < len(rs); k++ {
				if s[i-k] != rs[last-k] {
					found = false
					break
				}
			}
			if found {
				return i - last
			}
		}
	}
	return -1
}

// IndexAny returns the index of the first instance of any Unicode code point
// from chars in s, or -1 if no Unicode code point from chars is present in s.
func IndexAny(s, chars string) int {
	if chars == "" {
		// Avoid scanning all of s.
		return -1
	}
	if len(chars) == 1 {
		// Avoid scanning all of s.
		r := rune(chars[0])
		if r >= utf8.RuneSelf {
			r = utf8.RuneError
		}
		return IndexRune(s, r)
	}
	if len(s) > 8 {
		if as := makeASCIISet(chars); as.ok {
			for i := 0; i < len(s); i++ {
				if as.contains(s[i]) {
					return i
				}
			}
			return -1
		}
	}
	for i, c := range s {
		if IndexRune(chars, c) >= 0 {
			return i
		}
	}
	return -1
}

// LastIndexByte returns the index of the last instance of c in s, or -1 if c is not present in s.
func LastIndexByte(s string, c byte) int {
	return bytealg.LastIndexByteString(s, c)
}

// IndexFunc returns the index into s of the first Unicode
// code point satisfying f(c), or -1 if none do.
func IndexFunc(s string, f RunePredicate) int {
	return indexFunc(s, f, true)
}

// indexFunc is the same as IndexFunc except that if
// truth==false, the sense of the predicate function is
// inverted.
func indexFunc(s string, f RunePredicate, truth bool) int {
	for i, r := range s {
		if f(r) == truth {
			return i
		}
	}
	return -1
}

// lastIndexFunc is the same as LastIndexFunc except that if
// truth==false, the sense of the predicate function is
// inverted.
func lastIndexFunc(s string, f RunePredicate, truth bool) int {
	for i := len(s); i > 0; {
		r, size := utf8.DecodeLastRuneInString(s[0:i])
		i -= size
		if f(r) == truth {
			return i
		}
	}
	return -1
}
