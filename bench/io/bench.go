package main

import (
	"solod.dev/so/bytes"
	"solod.dev/so/fmt"
	"solod.dev/so/io"
	"solod.dev/so/mem"
	"solod.dev/so/testing"
)

func CopyNSmall(b *testing.B) {
	a := b.Allocator()

	bs := make([]byte, 512+1)
	rd := bytes.NewReader(bs)
	buf := bytes.NewBuffer(a, nil)
	defer buf.Free()

	for b.Loop() {
		io.CopyN(&buf, &rd, 512)
		rd.Reset(bs)
	}
}

func CopyNLarge(b *testing.B) {
	a := b.Allocator()

	bs := make([]byte, 32*1024+1)
	rd := bytes.NewReader(bs)
	buf := bytes.NewBuffer(a, nil)
	defer buf.Free()

	for b.Loop() {
		io.CopyN(&buf, &rd, 32*1024)
		rd.Reset(bs)
	}
}

func main() {
	benchs := []testing.Benchmark{
		{Name: "CopyNSmall", F: CopyNSmall},
		{Name: "CopyNLarge", F: CopyNLarge},
	}

	fmt.Println("Malloc-based allocator:")
	testing.RunBenchmarks(mem.System, benchs)
}
