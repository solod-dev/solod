// Copyright 2009 The Go Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package bytes

import (
	"solod.dev/so/mem"
	"solod.dev/so/unicode/utf8"
)

// Split slices s into all subslices separated by sep and returns a slice of
// the subslices between those separators.
// If sep is empty, Split splits after each UTF-8 sequence.
// It is equivalent to SplitN with a count of -1.
//
// To split around the first instance of a separator, see [Cut].
func Split(a mem.Allocator, s, sep []byte) [][]byte {
	return genSplit(a, s, sep, 0, -1)
}

// SplitN slices s into subslices separated by sep and returns a slice of
// the subslices between those separators.
// If sep is empty, SplitN splits after each UTF-8 sequence.
// The count determines the number of subslices to return:
//   - n > 0: at most n subslices; the last subslice will be the unsplit remainder;
//   - n == 0: the result is nil (zero subslices);
//   - n < 0: all subslices.
//
// To split around the first instance of a separator, see [Cut].
func SplitN(a mem.Allocator, s, sep []byte, n int) [][]byte {
	return genSplit(a, s, sep, 0, n)
}

// Generic split: splits after each instance of sep,
// including sepSave bytes of sep in the subslices.
//
// The returned slice is allocated; the caller owns it.
func genSplit(a mem.Allocator, s, sep []byte, sepSave, n int) [][]byte {
	if n == 0 {
		return [][]byte{}
	}
	if len(sep) == 0 {
		return explode(a, s, n)
	}
	if n < 0 {
		n = Count(s, sep) + 1
	}
	if n > len(s)+1 {
		n = len(s) + 1
	}

	res := mem.AllocSlice[[]byte](a, n, n)
	n--
	i := 0
	for i < n {
		m := Index(s, sep)
		if m < 0 {
			break
		}
		res[i] = s[: m+sepSave : m+sepSave]
		s = s[m+len(sep):]
		i++
	}
	res[i] = s
	return res[:i+1]
}

// explode splits s into a slice of UTF-8 sequences, one per Unicode code point
// (still slices of bytes), up to a maximum of n byte slices. Invalid UTF-8
// sequences are chopped into individual bytes.
//
// The returned slice is allocated; the caller owns it.
func explode(a mem.Allocator, s []byte, n int) [][]byte {
	if n <= 0 || n > len(s) {
		n = len(s)
	}
	res := mem.AllocSlice[[]byte](a, n, n)
	var size int
	na := 0
	for len(s) > 0 {
		if na+1 >= n {
			res[na] = s
			na++
			break
		}
		_, size = utf8.DecodeRune(s)
		res[na] = s[0:size:size]
		s = s[size:]
		na++
	}
	return res[0:na]
}
