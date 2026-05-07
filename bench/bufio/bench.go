package main

import (
	"solod.dev/so/bufio"
	"solod.dev/so/bytes"
	"solod.dev/so/mem"
	"solod.dev/so/slices"
	"solod.dev/so/strings"
	"solod.dev/so/testing"
)

//so:volatile
var sinkStr string

//so:volatile
var sinkInt int64

func ReaderBuf(b *testing.B) {
	a := b.Allocator()
	data := slices.Make[byte](nil, 16<<10)
	defer slices.Free(nil, data)
	r := bytes.NewReader(data)

	for b.Loop() {
		r.Reset(data)
		src := bufio.NewReader(a, &r)
		dst := bytes.NewBuffer(a, nil)
		sinkInt, _ = src.WriteTo(&dst)
		dst.Free()
		src.Free()
	}
}

func ReaderUnbuf(b *testing.B) {
	a := b.Allocator()
	data := slices.Make[byte](nil, 16<<10)
	defer slices.Free(nil, data)
	src := bytes.NewReader(data)

	for b.Loop() {
		src.Reset(data)
		dst := bytes.NewBuffer(a, nil)
		sinkInt, _ = src.WriteTo(&dst)
		dst.Free()
	}
}

func WriterBuf(b *testing.B) {
	a := b.Allocator()
	data := slices.Make[byte](nil, 16<<10)
	defer slices.Free(nil, data)
	r := bytes.NewReader(data)

	for b.Loop() {
		r.Reset(data)
		dstBuf := bytes.NewBuffer(a, nil)
		dst := bufio.NewWriter(a, &dstBuf)
		sinkInt, _ = dst.ReadFrom(&r)
		dst.Free()
		dstBuf.Free()
	}
}

func WriterUnbuf(b *testing.B) {
	a := b.Allocator()
	data := slices.Make[byte](nil, 16<<10)
	defer slices.Free(nil, data)
	r := bytes.NewReader(data)

	for b.Loop() {
		r.Reset(data)
		dst := bytes.NewBuffer(a, nil)
		sinkInt, _ = dst.ReadFrom(&r)
		dst.Free()
	}
}

func Scanner(b *testing.B) {
	const text = "Lorem ipsum dolor sit amet,\nconsectetur adipiscing elit,\nsed do eiusmod tempor incididunt ut labore\net dolore magna aliqua."
	a := b.Allocator()
	for b.Loop() {
		r := strings.NewReader(text)
		sc := bufio.NewScanner(a, &r)
		for sc.Scan() {
			sinkStr = sc.Text()
		}
		if err := sc.Err(); err != nil {
			panic(err)
		}
		sc.Free()
	}
}

func main() {
	benchs := []testing.Benchmark{
		{Name: "ReaderBuf", F: ReaderBuf},
		{Name: "ReaderUnbuf", F: ReaderUnbuf},
		{Name: "WriterBuf", F: WriterBuf},
		{Name: "WriterUnbuf", F: WriterUnbuf},
		{Name: "Scanner", F: Scanner},
	}
	testing.RunBenchmarks(mem.System, benchs)
}
