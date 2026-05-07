package main

// Unexported volatile variable.
//
//so:volatile
var counter int

// Exported volatile variable.
//
//so:volatile
var Counter int

// Unexported thread-local variable.
//
//so:thread_local
var perThread int

// Exported thread-local variable.
//
//so:thread_local
var PerThread int

// Combined volatile + thread-local.
//
//so:volatile
//so:thread_local
var flags int

// Packed struct with so:attr.
//
//so:attr packed
type packed struct {
	a byte
	b int
}

// Struct with multiple attrs.
//
//so:attr packed
//so:attr aligned(16)
type aligned struct {
	x int
}

// Exported struct with so:attr.
//
//so:attr packed
type Exported struct {
	v int
}

// Typedef alias with so:attr.
//
//so:attr aligned(8)
type myInt int

// Exported function with so:attr.
//
//so:attr noinline
func Work() {
}

// Unexported function with so:attr.
//
//so:attr noinline
func helper() {
}

func main() {
	counter = 1
	Counter = 2
	perThread = 3
	PerThread = 4
	flags = 5
	_ = packed{a: 1, b: 2}
	_ = aligned{x: 3}
	_ = Exported{v: 4}
	var m myInt = 5
	_ = m
	Work()
	helper()
}
