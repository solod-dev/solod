package main

// Primitive types.
type ID int          // not a different type
type AliasedID = int // also int
type AlsoID ID       // also int
type Rune rune

// Complex types.
type Name string
type IntArray [3]int
type IntSlice []int
type IntPtr *int
type Any interface{}
type Empty struct{}

// Struct type.
type Person struct {
	name string
	age  int
}

func newPerson(name string) Person {
	p := Person{name: name}
	p.age = 42
	return p
}

// Alias for a struct type.
type Human = Person
type Employee Person

// Inner struct.
type Benchmark struct {
	name string
	loop struct {
		n int
		i int
	}
}

func main() {
	{
		// Primitive types.
		var id ID = 123
		_ = id

		var aid AliasedID = 456
		_ = aid

		var alsoID AlsoID = 789
		_ = alsoID

		var r Rune = 'A'
		_ = r
	}
	{
		// Complex types.
		var n Name = "Alice"
		_ = n

		var arr IntArray = [3]int{1, 2, 3}
		_ = arr

		var slice IntSlice = []int{4, 5, 6}
		_ = slice
	}
	{
		// Struct types.
		bob := Person{"Bob", 20}
		_ = bob

		alice := Person{name: "Alice", age: 30}
		_ = alice

		fred := Person{name: "Fred"}
		_ = fred

		ann := &Person{name: "Ann", age: 40}
		*ann = newPerson("Jon")
		_ = ann

		var sean Person
		sean.name = "Sean"
		sean.age = 50
		sp := &sean
		sp.age = 51
		_ = sean
	}
	{
		// Anonymous struct type.
		dog := struct {
			name   string
			isGood bool
		}{
			"Rex",
			true,
		}
		_ = dog
	}
	{
		// Named struct type inside a function.
		type Point struct {
			x, y int
		}
		p := Point{1, 2}
		_ = p
	}
	{
		// Inner struct.
		b1 := Benchmark{name: "Test"}
		b1.loop.n = 100
		if b1.loop.n != 100 {
			panic("b1.loop.n != 100")
		}
		b2 := Benchmark{name: "Test2", loop: struct{ n, i int }{n: 200, i: 10}}
		if b2.loop.n != 200 {
			panic("b2.loop.n != 200")
		}
		b3 := Benchmark{name: "Test3", loop: struct{ n, i int }{300, 30}}
		if b3.loop.n != 300 {
			panic("b3.loop.n != 300")
		}
		var b4 Benchmark
		if b4.loop.n != 0 {
			panic("b4.loop.n != 0")
		}
	}
}
