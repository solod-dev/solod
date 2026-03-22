// Copyright 2009 The Go Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package bytes

import (
	"solod.dev/so/unicode"
	"solod.dev/so/unicode/utf8"
)

var asciiSpace = [256]uint8{'\t': 1, '\n': 1, '\v': 1, '\f': 1, '\r': 1, ' ': 1}

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

// contains reports whether c is inside the set.
func (as *asciiSet) contains(c byte) bool {
	return (as.val[c/32] & (1 << (c % 32))) != 0
}

// RunePredicate reports whether the rune satisfies a condition.
type RunePredicate func(rune) bool

// Trim returns a subslice of s by slicing off all leading and
// trailing UTF-8-encoded code points contained in cutset.
func Trim(s []byte, cutset string) []byte {
	if len(s) == 0 {
		return s
	}
	if cutset == "" {
		return s
	}
	if len(cutset) == 1 && cutset[0] < utf8.RuneSelf {
		return trimLeftByte(trimRightByte(s, cutset[0]), cutset[0])
	}
	if as := makeASCIISet(cutset); as.ok {
		return trimLeftASCII(trimRightASCII(s, &as), &as)
	}
	return trimLeftUnicode(trimRightUnicode(s, cutset), cutset)
}

// TrimFunc returns a subslice of s by slicing off all leading and trailing
// UTF-8-encoded code points c that satisfy f(c).
func TrimFunc(s []byte, f RunePredicate) []byte {
	return trimRightFunc(trimLeftFunc(s, f), f)
}

// TrimLeft returns a subslice of s by slicing off all leading
// UTF-8-encoded code points contained in cutset.
func TrimLeft(s []byte, cutset string) []byte {
	if len(s) == 0 {
		return s
	}
	if cutset == "" {
		return s
	}
	if len(cutset) == 1 && cutset[0] < utf8.RuneSelf {
		return trimLeftByte(s, cutset[0])
	}
	if as := makeASCIISet(cutset); as.ok {
		return trimLeftASCII(s, &as)
	}
	return trimLeftUnicode(s, cutset)
}

// TrimPrefix returns s without the provided leading prefix string.
// If s doesn't start with prefix, s is returned unchanged.
func TrimPrefix(s, prefix []byte) []byte {
	if HasPrefix(s, prefix) {
		return s[len(prefix):]
	}
	return s
}

// TrimRight returns a subslice of s by slicing off all trailing
// UTF-8-encoded code points that are contained in cutset.
func TrimRight(s []byte, cutset string) []byte {
	if len(s) == 0 || cutset == "" {
		return s
	}
	if len(cutset) == 1 && cutset[0] < utf8.RuneSelf {
		return trimRightByte(s, cutset[0])
	}
	if as := makeASCIISet(cutset); as.ok {
		return trimRightASCII(s, &as)
	}
	return trimRightUnicode(s, cutset)
}

// TrimSpace returns a subslice of s by slicing off all leading and
// trailing white space, as defined by Unicode.
func TrimSpace(s []byte) []byte {
	// Fast path for ASCII: look for the first ASCII non-space byte.
	for lo, c := range s {
		if c >= utf8.RuneSelf {
			// If we run into a non-ASCII byte, fall back to the
			// slower unicode-aware method on the remaining bytes.
			return TrimFunc(s[lo:], unicode.IsSpace)
		}
		if asciiSpace[c] != 0 {
			continue
		}
		s = s[lo:]
		// Now look for the first ASCII non-space byte from the end.
		for hi := len(s) - 1; hi >= 0; hi-- {
			c := s[hi]
			if c >= utf8.RuneSelf {
				return TrimFunc(s[:hi+1], unicode.IsSpace)
			}
			if asciiSpace[c] == 0 {
				// At this point, s[:hi+1] starts and ends with ASCII
				// non-space bytes, so we're done. Non-ASCII cases have
				// already been handled above.
				return s[:hi+1]
			}
		}
	}
	return []byte{}
}

// TrimSuffix returns s without the provided trailing suffix string.
// If s doesn't end with suffix, s is returned unchanged.
func TrimSuffix(s, suffix []byte) []byte {
	if HasSuffix(s, suffix) {
		return s[:len(s)-len(suffix)]
	}
	return s
}

func trimLeftASCII(s []byte, as *asciiSet) []byte {
	for len(s) > 0 {
		if !as.contains(s[0]) {
			break
		}
		s = s[1:]
	}
	if len(s) == 0 {
		return []byte{}
	}
	return s
}

func trimLeftByte(s []byte, c byte) []byte {
	for len(s) > 0 && s[0] == c {
		s = s[1:]
	}
	if len(s) == 0 {
		return []byte{}
	}
	return s
}

// trimLeftFunc treats s as UTF-8-encoded bytes and returns a subslice of s by slicing off
// all leading UTF-8-encoded code points c that satisfy f(c).
func trimLeftFunc(s []byte, f RunePredicate) []byte {
	i := indexFunc(s, f, false)
	if i == -1 {
		return []byte{}
	}
	return s[i:]
}

func trimLeftUnicode(s []byte, cutset string) []byte {
	for len(s) > 0 {
		r, n := utf8.DecodeRune(s)
		if !containsRune(cutset, r) {
			break
		}
		s = s[n:]
	}
	if len(s) == 0 {
		return []byte{}
	}
	return s
}

func trimRightASCII(s []byte, as *asciiSet) []byte {
	for len(s) > 0 {
		if !as.contains(s[len(s)-1]) {
			break
		}
		s = s[:len(s)-1]
	}
	return s
}

func trimRightByte(s []byte, c byte) []byte {
	for len(s) > 0 && s[len(s)-1] == c {
		s = s[:len(s)-1]
	}
	return s
}

// trimRightFunc returns a subslice of s by slicing off all trailing
// UTF-8-encoded code points c that satisfy f(c).
func trimRightFunc(s []byte, f RunePredicate) []byte {
	i := lastIndexFunc(s, f, false)
	if i >= 0 && s[i] >= utf8.RuneSelf {
		_, wid := utf8.DecodeRune(s[i:])
		i += wid
	} else {
		i++
	}
	return s[0:i]
}

func trimRightUnicode(s []byte, cutset string) []byte {
	for len(s) > 0 {
		r, n := rune(s[len(s)-1]), 1
		if r >= utf8.RuneSelf {
			r, n = utf8.DecodeLastRune(s)
		}
		if !containsRune(cutset, r) {
			break
		}
		s = s[:len(s)-n]
	}
	return s
}

// containsRune is a simplified version of strings.ContainsRune
// to avoid importing the strings package.
func containsRune(s string, r rune) bool {
	for _, c := range s {
		if c == r {
			return true
		}
	}
	return false
}

// indexFunc is the same as IndexFunc except that if
// truth==false, the sense of the predicate function is
// inverted.
func indexFunc(s []byte, f RunePredicate, truth bool) int {
	start := 0
	for start < len(s) {
		r, wid := utf8.DecodeRune(s[start:])
		if f(r) == truth {
			return start
		}
		start += wid
	}
	return -1
}

// lastIndexFunc is the same as LastIndexFunc except that if
// truth==false, the sense of the predicate function is
// inverted.
func lastIndexFunc(s []byte, f RunePredicate, truth bool) int {
	for i := len(s); i > 0; {
		r, size := rune(s[i-1]), 1
		if r >= utf8.RuneSelf {
			r, size = utf8.DecodeLastRune(s[0:i])
		}
		i -= size
		if f(r) == truth {
			return i
		}
	}
	return -1
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
