package main

import (
	"solod.dev/so/math/rand"
	"solod.dev/so/mem"
	"solod.dev/so/testing"
)

//so:volatile
var sink uint64

func SourceUint64(b *testing.B) {
	s := rand.NewPCG(1, 2)
	var t uint64
	for b.Loop() {
		t += s.Uint64()
	}
	sink = uint64(t)
}

func GlobalUint64(b *testing.B) {
	var t uint64
	for b.Loop() {
		t += rand.Uint64()
	}
	sink = t
}

func Uint64(b *testing.B) {
	pcg := rand.NewPCG(1, 2)
	r := rand.New(&pcg)
	var t uint64
	for b.Loop() {
		t += r.Uint64()
	}
	sink = t
}

func Int64N1e9(b *testing.B) {
	pcg := rand.NewPCG(1, 2)
	r := rand.New(&pcg)
	var t int64
	for b.Loop() {
		t += r.Int64N(1e9)
	}
	sink = uint64(t)
}

func Int64N1e18(b *testing.B) {
	pcg := rand.NewPCG(1, 2)
	r := rand.New(&pcg)
	var t int64
	for b.Loop() {
		t += r.Int64N(1e18)
	}
	sink = uint64(t)
}

func Int64N4e18(b *testing.B) {
	pcg := rand.NewPCG(1, 2)
	r := rand.New(&pcg)
	var t int64
	for b.Loop() {
		t += r.Int64N(4e18)
	}
	sink = uint64(t)
}

func Float64(b *testing.B) {
	pcg := rand.NewPCG(1, 2)
	r := rand.New(&pcg)
	var t float64
	for b.Loop() {
		t += r.Float64()
	}
	sink = uint64(t)
}

func main() {
	benchs := []testing.Benchmark{
		{Name: "SourceUint64", F: SourceUint64},
		{Name: "GlobalUint64", F: GlobalUint64},
		{Name: "Uint64", F: Uint64},
		{Name: "Int64N1e9", F: Int64N1e9},
		{Name: "Int64N1e18", F: Int64N1e18},
		{Name: "Int64N4e18", F: Int64N4e18},
		{Name: "Float64", F: Float64},
	}
	testing.RunBenchmarks(mem.System, benchs)
}
