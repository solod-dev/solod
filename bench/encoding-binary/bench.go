package main

import (
	"solod.dev/so/encoding/binary"
	"solod.dev/so/mem"
	"solod.dev/so/testing"
)

type Sink struct {
	buf [8]byte
	res []byte
}

//so:volatile
var sink Sink

//so:volatile
var sinkInt int

func BE_PutUint64(b *testing.B) {
	b.SetBytes(8)
	for i := range b.N {
		binary.BigEndian.PutUint64(sink.buf[:8], uint64(i))
		sinkInt = i
	}
}

func BE_AppendUint64(b *testing.B) {
	b.SetBytes(8)
	for i := range b.N {
		sink.res = binary.BigEndian.AppendUint64(sink.buf[:0], uint64(i))
	}
}

func LE_PutUint64(b *testing.B) {
	b.SetBytes(8)
	for i := range b.N {
		binary.LittleEndian.PutUint64(sink.buf[:8], uint64(i))
		sinkInt = i
	}
}

func LE_AppendUint64(b *testing.B) {
	b.SetBytes(8)
	for i := range b.N {
		sink.res = binary.LittleEndian.AppendUint64(sink.buf[:0], uint64(i))
	}
}

func main() {
	benchs := []testing.Benchmark{
		{Name: "BE_PutUint64", F: BE_PutUint64},
		{Name: "BE_AppendUint64", F: BE_AppendUint64},
		{Name: "LE_PutUint64", F: LE_PutUint64},
		{Name: "LE_AppendUint64", F: LE_AppendUint64},
	}
	testing.RunBenchmarks(mem.System, benchs)
}
