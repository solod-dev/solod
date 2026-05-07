package main

import (
	"solod.dev/so/fmt"
	"solod.dev/so/maps"
	"solod.dev/so/mem"
	"solod.dev/so/strings"
	"solod.dev/so/testing"
)

var arena *mem.Arena

//so:volatile
var sinkInt int

//so:volatile
var sinkBool bool

// nKeys is the number of map keys to use in benchmarks.
const nKeys = 1024

// strKeys holds pre-generated string keys for string benchmarks.
var strKeys []string

func initStrKeys() {
	strKeys = mem.AllocSlice[string](nil, nKeys, nKeys)
	buf := fmt.NewBuffer(32)
	for i := range nKeys {
		strKeys[i] = strings.Clone(nil, fmt.Sprintf(buf, "key-%d", i))
	}
}

func freeStrKeys() {
	for i := range nKeys {
		mem.FreeString(nil, strKeys[i])
	}
	mem.FreeSlice(nil, strKeys)
}

func StdIntSet(b *testing.B) {
	a := b.Allocator()
	for b.Loop() {
		m := maps.New[int, int](a, 0)
		for i := range nKeys {
			m.Set(i, i)
		}
		m.Free()
		if arena != nil {
			arena.Reset()
		}
	}
}

func StdIntPre(b *testing.B) {
	a := b.Allocator()
	for b.Loop() {
		m := maps.New[int, int](a, nKeys)
		for i := range nKeys {
			m.Set(i, i)
		}
		m.Free()
		if arena != nil {
			arena.Reset()
		}
	}
}

func StdIntGet(b *testing.B) {
	m := maps.New[int, int](nil, nKeys)
	for i := range nKeys {
		m.Set(i, i)
	}
	defer m.Free()
	for b.Loop() {
		for i := range nKeys {
			sinkInt = m.Get(i)
		}
	}
}

func StdIntDel(b *testing.B) {
	a := b.Allocator()
	for b.Loop() {
		m := maps.New[int, int](a, nKeys)
		for i := range nKeys {
			m.Set(i, i)
		}
		for i := range nKeys {
			m.Delete(i)
		}
		m.Free()
		if arena != nil {
			arena.Reset()
		}
	}
}

func StdStrSet(b *testing.B) {
	a := b.Allocator()
	for b.Loop() {
		m := maps.New[string, int](a, 0)
		for i := range nKeys {
			m.Set(strKeys[i], i)
		}
		m.Free()
		if arena != nil {
			arena.Reset()
		}
	}
}

func StdStrPre(b *testing.B) {
	a := b.Allocator()
	for b.Loop() {
		m := maps.New[string, int](a, nKeys)
		for i := range nKeys {
			m.Set(strKeys[i], i)
		}
		m.Free()
		if arena != nil {
			arena.Reset()
		}
	}
}

func StdStrGet(b *testing.B) {
	m := maps.New[string, int](nil, nKeys)
	for i := range nKeys {
		m.Set(strKeys[i], i)
	}
	defer m.Free()
	for b.Loop() {
		for i := range nKeys {
			sinkInt = m.Get(strKeys[i])
		}
	}
}

func StdStrDel(b *testing.B) {
	a := b.Allocator()
	for b.Loop() {
		m := maps.New[string, int](a, nKeys)
		for i := range nKeys {
			m.Set(strKeys[i], i)
		}
		for i := range nKeys {
			m.Delete(strKeys[i])
		}
		m.Free()
		if arena != nil {
			arena.Reset()
		}
	}
}

func BtlIntSet(b *testing.B) {
	for b.Loop() {
		bltIntSet(nKeys) // alloca only frees when the function returns
	}
}

func bltIntSet(nKeys int) {
	m := make(map[int]int, nKeys)
	for i := range nKeys {
		m[i] = i
	}
	sinkInt = m[0]
}

func BtlIntGet(b *testing.B) {
	m := make(map[int]int, nKeys)
	for i := range nKeys {
		m[i] = i
	}
	for b.Loop() {
		for i := range nKeys {
			sinkInt = m[i]
		}
	}
}

func BtlStrSet(b *testing.B) {
	for b.Loop() {
		bltStrSet(nKeys) // alloca only frees when the function returns
	}
}

func bltStrSet(nKeys int) {
	m := make(map[string]int, nKeys)
	for i := range nKeys {
		m[strKeys[i]] = i
	}
	sinkInt = m[strKeys[0]]
}

func BtlStrGet(b *testing.B) {
	m := make(map[string]int, nKeys)
	for i := range nKeys {
		m[strKeys[i]] = i
	}
	for b.Loop() {
		for i := range nKeys {
			sinkInt = m[strKeys[i]]
		}
	}
}

func main() {
	initStrKeys()
	defer freeStrKeys()

	benchs := []testing.Benchmark{
		{Name: "IntSet", F: StdIntSet},
		{Name: "IntPre", F: StdIntPre},
		{Name: "IntGet", F: StdIntGet},
		{Name: "IntDel", F: StdIntDel},
		{Name: "StrSet", F: StdStrSet},
		{Name: "StrPre", F: StdStrPre},
		{Name: "StrGet", F: StdStrGet},
		{Name: "StrDel", F: StdStrDel},
	}

	fmt.Println("Malloc-based allocator:")
	testing.RunBenchmarks(mem.System, benchs)

	fmt.Println("Arena allocator:")
	const size = 4 << 20
	buf := mem.AllocSlice[byte](nil, size, size)
	defer mem.FreeSlice(nil, buf)
	a := mem.NewArena(buf[:])
	arena = &a
	testing.RunBenchmarks(arena, benchs)

	builtinBenchs := []testing.Benchmark{
		{Name: "IntSet", F: BtlIntSet},
		{Name: "IntGet", F: BtlIntGet},
		{Name: "StrSet", F: BtlStrSet},
		{Name: "StrGet", F: BtlStrGet},
	}
	fmt.Println("Built-in map:")
	testing.RunBenchmarks(mem.System, builtinBenchs)
}
