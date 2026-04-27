package main

import (
	"path"
	"testing"
)

func Benchmark_Join(b *testing.B) {
	b.ReportAllocs()
	parts := []string{"one", "two", "three", "four"}
	for b.Loop() {
		sinkStr = path.Join(parts...)
	}
}

func Benchmark_MatchTrue(b *testing.B) {
	pattern := "a*b*c*d*e*/f"
	s := "axbxcxdxexxx/f"
	b.ReportAllocs()
	for b.Loop() {
		sinkBool, sinkErr = path.Match(pattern, s)
	}
}

func Benchmark_MatchFalse(b *testing.B) {
	pattern := "a*b*c*d*e*/f"
	s := "axbxcxdxexxx/fff"
	b.ReportAllocs()
	for b.Loop() {
		sinkBool, sinkErr = path.Match(pattern, s)
	}
}
