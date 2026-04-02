package main

import (
	"solod.dev/so/fmt"
	"solod.dev/so/maps"
	"solod.dev/so/mem"
	"solod.dev/so/strings"
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

// strKeys holds pre-generated string keys for string benchmarks.
var strKeys []string

func initStrKeys() {
	strKeys = mem.AllocSlice[string](nil, N, N)
	buf := fmt.NewBuffer(32)
	for i := range N {
		strKeys[i] = strings.Clone(nil, fmt.Sprintf(buf, "key-%d", i))
	}
}

func freeStrKeys() {
	for i := range N {
		mem.FreeString(nil, strKeys[i])
	}
	mem.FreeSlice(nil, strKeys)
}

func IntSet(b *testing.B) {
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

func IntGet(b *testing.B) {
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

func IntHas(b *testing.B) {
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

func IntDelete(b *testing.B) {
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

func IntSetDel(b *testing.B) {
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

func StrSet(b *testing.B) {
	a := b.Allocator()
	for b.Loop() {
		m := maps.New[string, int](a, 0)
		for i := range N {
			m.Set(strKeys[i], i)
		}
		m.Free()
		if arena != nil {
			arena.Reset()
		}
	}
}

func StrGet(b *testing.B) {
	m := maps.New[string, int](nil, N)
	for i := range N {
		m.Set(strKeys[i], i)
	}
	defer m.Free()
	for b.Loop() {
		for i := range N {
			sinkInt = m.Get(strKeys[i])
		}
	}
}

func StrHas(b *testing.B) {
	m := maps.New[string, int](nil, N)
	for i := range N {
		m.Set(strKeys[i], i)
	}
	defer m.Free()
	for b.Loop() {
		for i := range N {
			sinkBool = m.Has(strKeys[i])
		}
	}
}

func StrDelete(b *testing.B) {
	a := b.Allocator()
	for b.Loop() {
		m := maps.New[string, int](a, N)
		for i := range N {
			m.Set(strKeys[i], i)
		}
		for i := range N {
			m.Delete(strKeys[i])
		}
		m.Free()
		if arena != nil {
			arena.Reset()
		}
	}
}

func StrSetDel(b *testing.B) {
	a := b.Allocator()
	for b.Loop() {
		m := maps.New[string, int](a, 0)
		for i := range N {
			m.Set(strKeys[i], i)
			m.Delete(strKeys[i])
		}
		m.Free()
		if arena != nil {
			arena.Reset()
		}
	}
}

func main() {
	initStrKeys()
	defer freeStrKeys()

	intBenchs := []testing.Benchmark{
		{Name: "IntSet", F: IntSet},
		{Name: "IntGet", F: IntGet},
		{Name: "IntHas", F: IntHas},
		{Name: "IntDelete", F: IntDelete},
		{Name: "IntSetDel", F: IntSetDel},
	}

	strBenchs := []testing.Benchmark{
		{Name: "StrSet", F: StrSet},
		{Name: "StrGet", F: StrGet},
		{Name: "StrHas", F: StrHas},
		{Name: "StrDelete", F: StrDelete},
		{Name: "StrSetDel", F: StrSetDel},
	}

	fmt.Println("Malloc-based allocator:")
	testing.RunBenchmarks(mem.System, intBenchs)
	testing.RunBenchmarks(mem.System, strBenchs)

	fmt.Println("Arena allocator:")
	const size = 4 << 20
	buf := mem.AllocSlice[byte](nil, size, size)
	defer mem.FreeSlice(nil, buf)
	a := mem.NewArena(buf[:])
	arena = &a
	testing.RunBenchmarks(arena, intBenchs)
	testing.RunBenchmarks(arena, strBenchs)
}
