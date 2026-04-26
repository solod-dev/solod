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

func main() {}
