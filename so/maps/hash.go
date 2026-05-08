package maps

import (
	"unsafe"

	"solod.dev/so/c"
	_ "solod.dev/so/math/bits"
	"solod.dev/so/mem"
)

// wymum performs 128-bit multiply-and-mix.
//
//so:extern maps_wymum
func wymum(a, b uint64) uint64 {
	return a * b // for testing
}

// wyr8 reads 8 bytes as a little-endian uint64.
//
//so:inline
func wyr8(p any) uint64 {
	var v uint64
	mem.Copy(&v, p, 8)
	return v
}

// wyr4 reads 4 bytes as a little-endian uint32.
//
//so:inline
func wyr4(p any) uint64 {
	var v uint32
	mem.Copy(&v, p, 4)
	return uint64(v)
}

// hash computes wyhash with a per-map seed.
//
//so:inline
func maps_hash(key any, length int, seed uint64) int {
	p := key.(*byte)
	wyp0 := uint64(0xa0761d6478bd642f)
	wyp1 := uint64(0xe7037ed1a0b428db)
	seed = wymum(seed^wyp0, wyp1)
	var a, b uint64
	if length > 16 {
		for i := 0; i+16 <= length; i += 16 {
			seed = wymum(wyr8(c.PtrAdd(p, i))^wyp1,
				wyr8(c.PtrAdd(p, i+8))^seed)
		}
		a = wyr8(c.PtrAdd(p, length-16))
		b = wyr8(c.PtrAdd(p, length-8))
	} else if length >= 4 {
		a = (wyr4(p) << 32) | wyr4(c.PtrAdd(p, (length>>3)<<2))
		b = (wyr4(c.PtrAdd(p, length-4)) << 32) |
			wyr4(c.PtrAdd(p, length-4-((length>>3)<<2)))
	} else if length > 0 {
		a = (uint64(*p) << 16) |
			(uint64(*c.PtrAdd(p, length>>1)) << 8) |
			uint64(*c.PtrAdd(p, length-1))
	}
	r := wymum(wyp1^uint64(length), wymum(a^wyp1, b^seed))
	return int(r >> 16) // upper 48 bits is the hash value
}

// hashString hashes a string key by its content.
//
//so:inline
func maps_hashString(keyPtr any, seed uint64) int {
	s := *keyPtr.(*string)
	return maps_hash(unsafe.StringData(s), len(s), seed)
}
