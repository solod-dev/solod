package main

import (
	"solod.dev/so/mem"
	"solod.dev/so/testing"
	"solod.dev/so/uuid"
)

//so:volatile
var sink uuid.UUID

//so:volatile
var sinkStr string

func NewV4(b *testing.B) {
	for b.Loop() {
		sink = uuid.NewV4()
	}
}

func NewV7(b *testing.B) {
	for b.Loop() {
		sink = uuid.NewV7()
	}
}

func String(b *testing.B) {
	u := uuid.MustParse("f81d4fae-7dec-11d0-a765-00a0c91e6bf6")
	buf := make([]byte, uuid.UUIDLen)
	for b.Loop() {
		sinkStr = u.String(buf)
	}
}

func ParseSuccess(b *testing.B) {
	for b.Loop() {
		sink = uuid.MustParse("f81d4fae-7dec-11d0-a765-00a0c91e6bf6")
	}
}

func ParseError(b *testing.B) {
	for b.Loop() {
		sink, _ = uuid.Parse("00000000-0000-0000-0000-00000000000X")
	}
}

func main() {
	benchs := []testing.Benchmark{
		{Name: "NewV4", F: NewV4},
		{Name: "NewV7", F: NewV7},
		{Name: "String", F: String},
		{Name: "ParseSuccess", F: ParseSuccess},
		{Name: "ParseError", F: ParseError},
	}
	testing.RunBenchmarks(mem.System, benchs)
}
