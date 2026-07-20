package main

// Self-referencing struct type.
type Node struct {
	value int
	next  *Node
}

// Type referencing another type defined later.
type Employee struct {
	name string
	pet  *Pet
}

type Pet struct {
	name string
}

// Unexported self-referencing struct type.
type unode struct {
	value int
	next  *unode
}

// Unexported type referencing another type defined later.
type uemployee struct {
	name string
	pet  *upet
}

type upet struct {
	name string
}

// Type using a type defined later by value.
type Rect struct {
	Min, Max Point
}

type Point struct {
	X, Y int
}

// Array of a type defined later.
type Grid struct {
	cells [4]Cell
}

type Cell struct {
	v int
}

// Named type of a type defined later.
type Target Origin

type Origin struct {
	v int
}

// Func type held by value: Handler must precede Outer, but the struct types
// in its signature are fine as forward declarations.
type Outer struct {
	handle Handler
}

type Handler func(Payload) Payload

type Payload struct {
	v int
}

// Pointer to a non-struct type: Meters has no forward declaration,
// so its definition must come first.
type Reading struct {
	depth *Meters
}

type Meters int

// Unexported type using a type defined later by value.
type urect struct {
	min, max upoint
}

type upoint struct {
	x, y int
}

func main() {}
