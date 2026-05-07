package main

import (
	"solod.dev/so/fmt"
	"solod.dev/so/mem"
	"solod.dev/so/strings"
	"solod.dev/so/testing"
)

var arena *mem.Arena
var someStr = "some string sdljlk jsklj3lkjlk djlkjw"

//so:volatile
var sink string

const numWrite = 16

func Write_AutoGrow(b *testing.B) {
	a := b.Allocator()
	someBytes := []byte(someStr)
	for b.Loop() {
		buf := strings.NewBuilder(a)
		for range numWrite {
			buf.Write(someBytes)
		}
		sink = buf.String()
		buf.Free()
		if arena != nil {
			arena.Reset()
		}
	}
}

func Write_PreGrow(b *testing.B) {
	a := b.Allocator()
	someBytes := []byte(someStr)
	for b.Loop() {
		buf := strings.NewBuilder(a)
		buf.Grow(len(someBytes) * numWrite)
		for range numWrite {
			buf.Write(someBytes)
		}
		sink = buf.String()
		buf.Free()
		if arena != nil {
			arena.Reset()
		}
	}
}

func WriteString_AutoGrow(b *testing.B) {
	a := b.Allocator()
	for b.Loop() {
		buf := strings.NewBuilder(a)
		for range numWrite {
			buf.WriteString(someStr)
		}
		sink = buf.String()
		buf.Free()
		if arena != nil {
			arena.Reset()
		}
	}
}

func WriteString_PreGrow(b *testing.B) {
	a := b.Allocator()
	for b.Loop() {
		buf := strings.NewBuilder(a)
		buf.Grow(len(someStr) * numWrite)
		for range numWrite {
			buf.WriteString(someStr)
		}
		sink = buf.String()
		buf.Free()
		if arena != nil {
			arena.Reset()
		}
	}
}

func main() {
	benchs := []testing.Benchmark{
		{Name: "WriteB_AutoGrow", F: Write_AutoGrow},
		{Name: "WriteB_PreGrow", F: Write_PreGrow},
		{Name: "WriteS_AutoGrow", F: WriteString_AutoGrow},
		{Name: "WriteS_PreGrow", F: WriteString_PreGrow},
	}

	fmt.Println("Malloc-based allocator:")
	testing.RunBenchmarks(mem.System, benchs)

	fmt.Println("Arena allocator:")
	var buf [4096]byte
	a := mem.NewArena(buf[:])
	arena = &a
	testing.RunBenchmarks(arena, benchs)
}
