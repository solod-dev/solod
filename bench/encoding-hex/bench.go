package main

import (
	"solod.dev/so/bytes"
	"solod.dev/so/encoding/hex"
	"solod.dev/so/mem"
	"solod.dev/so/slices"
	"solod.dev/so/testing"
)

//so:volatile
var sink []byte

func Encode_256(b *testing.B) {
	encode(b, 256)
}

func Encode_1024(b *testing.B) {
	encode(b, 1024)
}

func Encode_4096(b *testing.B) {
	encode(b, 4096)
}

func Encode_16384(b *testing.B) {
	encode(b, 16384)
}

func Decode_256(b *testing.B) {
	decode(b, 256)
}

func Decode_1024(b *testing.B) {
	decode(b, 1024)
}

func Decode_4096(b *testing.B) {
	decode(b, 4096)
}

func Decode_16384(b *testing.B) {
	decode(b, 16384)
}

func encode(b *testing.B, size int) {
	b.SetBytes(int64(size))
	alloc := b.Allocator()
	data := []byte{2, 3, 5, 7, 9, 11, 13, 17}
	src := bytes.Repeat(alloc, data, size/8)
	sink = slices.Make[byte](alloc, 2*size)
	for b.Loop() {
		hex.Encode(sink, src)
	}
	slices.Free(alloc, src)
}

func decode(b *testing.B, size int) {
	b.SetBytes(int64(size))
	alloc := b.Allocator()
	data := []byte{'2', 'b', '7', '4', '4', 'f', 'a', 'a'}
	src := bytes.Repeat(alloc, data, size/8)
	sink = slices.Make[byte](alloc, size/2)
	for b.Loop() {
		hex.Decode(sink, src)
	}
	slices.Free(alloc, src)
}

func main() {
	benchs := []testing.Benchmark{
		{Name: "Encode_256", F: Encode_256},
		{Name: "Encode_1024", F: Encode_1024},
		{Name: "Encode_4096", F: Encode_4096},
		{Name: "Encode_16384", F: Encode_16384},
		{Name: "Decode_256", F: Decode_256},
		{Name: "Decode_1024", F: Decode_1024},
		{Name: "Decode_4096", F: Decode_4096},
		{Name: "Decode_16384", F: Decode_16384},
	}
	testing.RunBenchmarks(mem.System, benchs)
}
