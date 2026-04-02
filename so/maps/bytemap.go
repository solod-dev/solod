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

const (
	loadFactor  = 0.85                      // must be above 50%
	dibBitSize  = 16                        // 0xFFFF
	hashBitSize = 64 - dibBitSize           // 0xFFFFFFFFFFFF
	maxDIB      = ^uint64(0) >> hashBitSize // max 65,535
)

// KeyEqualFn is a function that checks two key byte slices for equality.
type KeyEqualFn func(a, b []byte) bool

// EqualBytes returns true if two byte slices are equal.
//
//so:extern
func EqualBytes(a, b []byte) bool {
	return false
}

// ByteMap is a Robin Hood hashmap operating on raw byte keys and values.
// Most users will want to use the generic [Map] wrapper type instead.
type ByteMap struct {
	a       mem.Allocator
	equalFn KeyEqualFn

	hdib  []hashDIB // one per bucket: bitfield { hash:48 dib:16 }
	keys  []byte    // ksize bytes per bucket, packed
	vals  []byte    // vsize bytes per bucket, packed
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
	m.hdib = mem.AllocSlice[hashDIB](m.a, sz, sz)
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

func (m *ByteMap) set(hash int, key, value []byte) {
	ehdib := makeHDIB(hash, 1)
	ekey := make([]byte, m.ksize)
	eval := make([]byte, m.vsize)
	tmpk := make([]byte, m.ksize)
	tmpv := make([]byte, m.vsize)
	copy(ekey, key)
	copy(eval, value)
	i := hash & m.mask
	for {
		if m.hdib[i].dib() == 0 {
			m.hdib[i] = ehdib
			copy(m.keyAt(i), ekey)
			copy(m.valAt(i), eval)
			m.len++
			return
		}
		if ehdib.hash() == m.hdib[i].hash() && m.equalFn(ekey, m.keyAt(i)) {
			copy(m.valAt(i), eval)
			return
		}
		if m.hdib[i].dib() < ehdib.dib() {
			tmp := ehdib
			ehdib = m.hdib[i]
			m.hdib[i] = tmp
			copy(tmpk, ekey)
			copy(ekey, m.keyAt(i))
			copy(m.keyAt(i), tmpk)
			copy(tmpv, eval)
			copy(eval, m.valAt(i))
			copy(m.valAt(i), tmpv)
		}
		i = (i + 1) & m.mask
		ehdib = ehdib.setDIB(ehdib.dib() + 1)
	}
}

// Resize grows or reallocates the map to hold at least size entries.
func (m *ByteMap) Resize(size int) {
	nmap := NewByteMap(m.a, size, m.ksize, m.vsize)
	nmap.equalFn = m.equalFn
	nbuckets := len(m.hdib)
	for i := range nbuckets {
		if m.hdib[i].dib() > 0 {
			nmap.set(m.hdib[i].hash(), m.keyAt(i), m.valAt(i))
		}
	}
	m.Free()
	*m = nmap
}

func (m *ByteMap) keyAt(i int) []byte {
	off := i * m.ksize
	return m.keys[off : off+m.ksize]
}

func (m *ByteMap) valAt(i int) []byte {
	off := i * m.vsize
	return m.vals[off : off+m.vsize]
}

// hashDIB is a compact struct that stores both the hash
// and DIB (distance to initial bucket) for a map entry.
type hashDIB uint64

func makeHDIB(hash, dib int) hashDIB {
	val := (uint64(hash) << dibBitSize) | (uint64(dib) & maxDIB)
	return hashDIB(val)
}

func (hdib hashDIB) setDIB(d int) hashDIB {
	val := (uint64(hdib) >> dibBitSize << dibBitSize) | (uint64(d) & maxDIB)
	return hashDIB(val)
}

func (hdib hashDIB) hash() int {
	return int(uint64(hdib) >> dibBitSize)
}

func (hdib hashDIB) dib() int {
	return int(uint64(hdib) & maxDIB)
}
