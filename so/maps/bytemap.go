// Copyright 2019, Joshua J Baker
// <https://github.com/tidwall/hashmap>

// Permission to use, copy, modify, and/or distribute this software for any
// purpose with or without fee is hereby granted, provided that the above
// copyright notice and this permission notice appear in all copies.

// THE SOFTWARE IS PROVIDED "AS IS" AND THE AUTHOR DISCLAIMS ALL WARRANTIES
// WITH REGARD TO THIS SOFTWARE INCLUDING ALL IMPLIED WARRANTIES OF
// MERCHANTABILITY AND FITNESS. IN NO EVENT SHALL THE AUTHOR BE LIABLE FOR ANY
// SPECIAL, DIRECT, INDIRECT, OR CONSEQUENTIAL DAMAGES OR ANY DAMAGES
// WHATSOEVER RESULTING FROM LOSS OF USE, DATA OR PROFITS, WHETHER IN AN ACTION
// OF CONTRACT, NEGLIGENCE OR OTHER TORTIOUS ACTION, ARISING OUT OF OR IN
// CONNECTION WITH THE USE OR PERFORMANCE OF THIS SOFTWARE.

package maps

import "solod.dev/so/mem"

const loadFactor = 0.85 // must be above 50%

// ByteMap is a Robin Hood hashmap operating on raw byte keys and values.
// Most users will want to use the generic [Map] wrapper type instead.
type ByteMap struct {
	a mem.Allocator

	hdib  []uint64 // one per bucket: bitfield { hash:48 dib:16 }
	keys  []byte   // ksize bytes per bucket, packed
	vals  []byte   // vsize bytes per bucket, packed
	ksize int
	vsize int

	len    int // number of items in the map
	mask   int // mask for indexing into buckets
	growAt int // length at which to grow the map
}

// NewByteMap creates a new ByteMap with the given initial capacity,
// key size, and value size, using the provided allocator (or the
// default allocator if nil). The map automatically grows as needed.
//
// If the allocator is nil, uses the system allocator.
// The caller is responsible for freeing map resources
// with [ByteMap.Free] when done using it.
func NewByteMap(a mem.Allocator, size, ksize, vsize int) ByteMap {
	m := ByteMap{a: a, ksize: ksize, vsize: vsize}
	sz := 8
	for sz < size {
		sz *= 2
	}
	m.hdib = mem.AllocSlice[uint64](m.a, sz, sz)
	m.keys = mem.AllocSlice[byte](m.a, sz*ksize, sz*ksize)
	m.vals = mem.AllocSlice[byte](m.a, sz*vsize, sz*vsize)
	m.mask = sz - 1
	m.growAt = int(float64(sz) * loadFactor)
	return m
}

// Len returns the number of key-value pairs in the map.
func (m *ByteMap) Len() int {
	return m.len
}

// Free frees internal resources used by the map.
// If the map is already freed, does nothing.
// The map must not be used after Free.
func (m *ByteMap) Free() {
	if len(m.hdib) == 0 {
		return
	}
	mem.FreeSlice(m.a, m.hdib)
	mem.FreeSlice(m.a, m.keys)
	mem.FreeSlice(m.a, m.vals)
	m.hdib = nil
	m.keys = nil
	m.vals = nil
	m.len = 0
}

// Resize grows or reallocates the map to hold at least size entries.
func (m *ByteMap) Resize(size int) {
	nmap := NewByteMap(m.a, size, m.ksize, m.vsize)
	rehash(&nmap, m)
	m.Free()
	*m = nmap
}

//so:extern maps_rehash
func rehash(dst, src *ByteMap) {}
