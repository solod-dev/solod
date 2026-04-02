package main

import (
	"fmt"
	"testing"
)

const benchN = 1024

var strBenchKeys [benchN]string

func init() {
	for i := range benchN {
		strBenchKeys[i] = fmt.Sprintf("key-%d", i)
	}
}

func Benchmark_IntSet(b *testing.B) {
	b.ReportAllocs()
	for b.Loop() {
		m := make(map[int]int)
		for i := range benchN {
			m[i] = i
		}
	}
}

func Benchmark_IntGet(b *testing.B) {
	b.ReportAllocs()
	m := make(map[int]int, benchN)
	for i := range benchN {
		m[i] = i
	}
	for b.Loop() {
		for i := range benchN {
			sinkInt = m[i]
		}
	}
}

func Benchmark_IntHas(b *testing.B) {
	b.ReportAllocs()
	m := make(map[int]int, benchN)
	for i := range benchN {
		m[i] = i
	}
	for b.Loop() {
		for i := range benchN {
			_, sinkBool = m[i]
		}
	}
}

func Benchmark_IntDelete(b *testing.B) {
	b.ReportAllocs()
	for b.Loop() {
		m := make(map[int]int, benchN)
		for i := range benchN {
			m[i] = i
		}
		for i := range benchN {
			delete(m, i)
		}
	}
}

func Benchmark_IntSetDel(b *testing.B) {
	b.ReportAllocs()
	for b.Loop() {
		m := make(map[int]int)
		for i := range benchN {
			m[i] = i
			delete(m, i)
		}
	}
}

func Benchmark_StrSet(b *testing.B) {
	b.ReportAllocs()
	for b.Loop() {
		m := make(map[string]int)
		for i := range benchN {
			m[strBenchKeys[i]] = i
		}
	}
}

func Benchmark_StrGet(b *testing.B) {
	b.ReportAllocs()
	m := make(map[string]int, benchN)
	for i := range benchN {
		m[strBenchKeys[i]] = i
	}
	for b.Loop() {
		for i := range benchN {
			sinkInt = m[strBenchKeys[i]]
		}
	}
}

func Benchmark_StrHas(b *testing.B) {
	b.ReportAllocs()
	m := make(map[string]int, benchN)
	for i := range benchN {
		m[strBenchKeys[i]] = i
	}
	for b.Loop() {
		for i := range benchN {
			_, sinkBool = m[strBenchKeys[i]]
		}
	}
}

func Benchmark_StrDelete(b *testing.B) {
	b.ReportAllocs()
	for b.Loop() {
		m := make(map[string]int, benchN)
		for i := range benchN {
			m[strBenchKeys[i]] = i
		}
		for i := range benchN {
			delete(m, strBenchKeys[i])
		}
	}
}

func Benchmark_StrSetDel(b *testing.B) {
	b.ReportAllocs()
	for b.Loop() {
		m := make(map[string]int)
		for i := range benchN {
			m[strBenchKeys[i]] = i
			delete(m, strBenchKeys[i])
		}
	}
}
