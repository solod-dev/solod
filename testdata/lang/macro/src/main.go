package main

//so:embed main.h
var main_h string

//so:inline
func identity[T any](val T) T {
	return val
}

//so:inline
func setPtr[T any](ptr *T, val T) {
	*ptr = val
}

//so:inline
func a[T int](n T) T {
	var some int = 11
	_ = some
	x := b(n) + 1
	return x
}

//so:inline
func b[T int](n T) T {
	var some float64 = 22.2
	_ = some
	x := c(n) + 1
	return x
}

//so:inline
func c[T int](n T) T {
	var some string = "33"
	_ = some
	x := n + 1
	return x
}

//so:extern
type Box[T any] struct {
	val T
}

//so:inline
func (b *Box[T]) set(val T) {
	b.val = val
}

func main() {
	{
		// Function with return.
		x := identity(42)
		if x != int(42) {
			panic("Function with return failed")
		}
	}
	{
		// Function w/o return.
		var y int
		setPtr(&y, 42)
		if y != 42 {
			panic("Function w/o return failed")
		}
	}
	{
		// Nested calls with variable shadowing.
		z := a(42)
		if z != 45 {
			panic("Nested calls failed")
		}
	}
	{
		// Generic method.
		var b Box[int]
		b.set(42)
		if b.val != 42 {
			panic("Generic method failed")
		}
	}
	println("lang/macro ok")
}
