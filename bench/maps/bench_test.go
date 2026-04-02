package main

import "testing"

const benchN = 1024

func Benchmark_Set(b *testing.B) {
	b.ReportAllocs()
	for b.Loop() {
		m := make(map[int]int)
		for i := range benchN {
			m[i] = i
		}
	}
}

func Benchmark_Get(b *testing.B) {
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

func Benchmark_Has(b *testing.B) {
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

func Benchmark_Delete(b *testing.B) {
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

func Benchmark_SetDelete(b *testing.B) {
	b.ReportAllocs()
	for b.Loop() {
		m := make(map[int]int)
		for i := range benchN {
			m[i] = i
			delete(m, i)
		}
	}
}
