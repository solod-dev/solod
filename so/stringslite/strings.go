// Copyright 2024 The Go Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

// Package stringslite implements a subset of strings,
// only using packages that may be imported by "os".
// It's only meant for use by packages in the standard library.
// Others should use the strings package instead.
//
// Based on the [internal/stringslite] package.
//
// [internal/stringslite]: https://github.com/golang/go/blob/go1.26.1/src/internal/stringslite/strings.go
package stringslite

import (
	"github.com/nalgeon/solod/so/bytealg"
	"github.com/nalgeon/solod/so/mem"
)

func Clone(a mem.Allocator, s string) string {
	if len(s) == 0 {
		return ""
	}
	b := mem.AllocSlice[byte](a, len(s), len(s))
	copy(b, s)
	return string(b)
}

func Cut(s, sep string) (string, string) {
	if i := Index(s, sep); i >= 0 {
		return s[:i], s[i+len(sep):]
	}
	return s, ""
}

func CutPrefix(s, prefix string) (string, bool) {
	if !HasPrefix(s, prefix) {
		return s, false
	}
	return s[len(prefix):], true
}

func CutSuffix(s, suffix string) (string, bool) {
	if !HasSuffix(s, suffix) {
		return s, false
	}
	return s[:len(s)-len(suffix)], true
}

func HasPrefix(s, prefix string) bool {
	return len(s) >= len(prefix) && s[:len(prefix)] == prefix
}

func HasSuffix(s, suffix string) bool {
	return len(s) >= len(suffix) && s[len(s)-len(suffix):] == suffix
}

func Index(s, substr string) int {
	n := len(substr)
	if n == 0 {
		return 0
	} else if n == 1 {
		return IndexByte(s, substr[0])
	} else if n == len(s) {
		if substr == s {
			return 0
		}
		return -1
	} else if n > len(s) {
		return -1
	}

	c0 := substr[0]
	c1 := substr[1]
	i := 0
	t := len(s) - n + 1
	fails := 0
	for i < t {
		if s[i] != c0 {
			o := IndexByte(s[i+1:t], c0)
			if o < 0 {
				return -1
			}
			i += o + 1
		}
		if s[i+1] == c1 && s[i:i+n] == substr {
			return i
		}
		i++
		fails++
		if fails >= (4+(i>>4)) && i < t {
			// See comment in [bytes.Index].
			j := bytealg.IndexRabinKarp([]byte(s[i:]), []byte(substr))
			if j < 0 {
				return -1
			}
			return i + j
		}
	}
	return -1
}

func IndexByte(s string, c byte) int {
	return bytealg.IndexByteString(s, c)
}

func TrimPrefix(s, prefix string) string {
	if HasPrefix(s, prefix) {
		return s[len(prefix):]
	}
	return s
}

func TrimSuffix(s, suffix string) string {
	if HasSuffix(s, suffix) {
		return s[:len(s)-len(suffix)]
	}
	return s
}
