package main

import (
	"solod.dev/so/fmt"
	"solod.dev/so/mem"
	"solod.dev/so/path"
	"solod.dev/so/testing"
)

//so:embed bench.h
var bench_h string

var arena *mem.Arena

//so:extern nodecay
var (
	sinkBool bool
	sinkErr  error
	sinkStr  string
)

func Join(b *testing.B) {
	a := b.Allocator()
	parts := []string{"one", "two", "three", "four"}
	for b.Loop() {
		sinkStr = path.Join(a, parts...)
		mem.FreeString(a, sinkStr)
		if arena != nil {
			arena.Reset()
		}
	}
}

func MatchTrue(b *testing.B) {
	pattern := "a*b*c*d*e*/f"
	s := "axbxcxdxexxx/f"
	for b.Loop() {
		sinkBool, sinkErr = path.Match(pattern, s)
	}
}

func MatchFalse(b *testing.B) {
	pattern := "a*b*c*d*e*/f"
	s := "axbxcxdxexxx/fff"
	for b.Loop() {
		sinkBool, sinkErr = path.Match(pattern, s)
	}
}

func main() {
	benchs := []testing.Benchmark{
		{Name: "Join", F: Join},
		{Name: "MatchTrue", F: MatchTrue},
		{Name: "MatchFalse", F: MatchFalse},
	}

	fmt.Println("Malloc-based allocator:")
	testing.RunBenchmarks(mem.System, benchs)

	fmt.Println("Arena allocator:")
	var buf [4096]byte
	a := mem.NewArena(buf[:])
	arena = &a
	testing.RunBenchmarks(arena, benchs)
}
