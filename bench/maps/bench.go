package main

import (
	"solod.dev/so/fmt"
	"solod.dev/so/maps"
	"solod.dev/so/mem"
	"solod.dev/so/testing"
)

//so:embed bench.h
var bench_h string

var arena *mem.Arena

//so:extern nodecay
var (
	sinkInt  int
	sinkBool bool
)

const N = 1024

func Set(b *testing.B) {
	a := b.Allocator()
	for b.Loop() {
		m := maps.New[int, int](a, 0)
		for i := range N {
			m.Set(i, i)
		}
		m.Free()
		if arena != nil {
			arena.Reset()
		}
	}
}

func Get(b *testing.B) {
	m := maps.New[int, int](nil, N)
	for i := range N {
		m.Set(i, i)
	}
	defer m.Free()
	for b.Loop() {
		for i := range N {
			sinkInt = m.Get(i)
		}
	}
}

func Has(b *testing.B) {
	m := maps.New[int, int](nil, N)
	for i := range N {
		m.Set(i, i)
	}
	defer m.Free()
	for b.Loop() {
		for i := range N {
			sinkBool = m.Has(i)
		}
	}
}

func Delete(b *testing.B) {
	a := b.Allocator()
	for b.Loop() {
		m := maps.New[int, int](a, N)
		for i := range N {
			m.Set(i, i)
		}
		for i := range N {
			m.Delete(i)
		}
		m.Free()
		if arena != nil {
			arena.Reset()
		}
	}
}

func SetDelete(b *testing.B) {
	a := b.Allocator()
	for b.Loop() {
		m := maps.New[int, int](a, 0)
		for i := range N {
			m.Set(i, i)
			m.Delete(i)
		}
		m.Free()
		if arena != nil {
			arena.Reset()
		}
	}
}

func main() {
	benchs := []testing.Benchmark{
		{Name: "Set", F: Set},
		{Name: "Get", F: Get},
		{Name: "Has", F: Has},
		{Name: "Delete", F: Delete},
		{Name: "SetDelete", F: SetDelete},
	}

	fmt.Println("Malloc-based allocator:")
	testing.RunBenchmarks(mem.System, benchs)

	fmt.Println("Arena allocator:")
	var buf [1024 * 1024]byte
	a := mem.NewArena(buf[:])
	arena = &a
	testing.RunBenchmarks(arena, benchs)
}
