package main

import (
	"solod.dev/so/crypto/crand"
	"solod.dev/so/mem"
	"solod.dev/so/testing"
)

//so:volatile
var sinkSlice []byte

//so:volatile
var sinkStr string

func Read_4(b *testing.B) {
	const size = 4
	b.SetBytes(int64(size))
	buf := make([]byte, size)
	for b.Loop() {
		if _, err := crand.Read(buf); err != nil {
			panic(err)
		}
	}
}

func Read_32(b *testing.B) {
	const size = 32
	b.SetBytes(int64(size))
	buf := make([]byte, size)
	for b.Loop() {
		if _, err := crand.Read(buf); err != nil {
			panic(err)
		}
	}
}

func Read_4K(b *testing.B) {
	const size = 4 << 10
	b.SetBytes(int64(size))
	buf := make([]byte, size)
	for b.Loop() {
		if _, err := crand.Read(buf); err != nil {
			panic(err)
		}
	}
}

func Text(b *testing.B) {
	const size = 26
	b.SetBytes(size)
	buf := make([]byte, size)
	for b.Loop() {
		sinkStr = crand.Text(buf)
	}
}

func main() {
	benchs := []testing.Benchmark{
		{Name: "Read_4", F: Read_4},
		{Name: "Read_32", F: Read_32},
		{Name: "Read_4K", F: Read_4K},
		{Name: "Text", F: Text},
	}
	testing.RunBenchmarks(mem.System, benchs)
}
