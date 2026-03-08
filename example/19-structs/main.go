// So's _structs_ are typed collections of fields.
// They're useful for grouping data together to form
// records.
package main

// This `person` struct type has `name` and `age` fields.
type person struct {
	name string
	age  int
}

// `newPerson` constructs a new person struct with the given name.
func newPerson(name string) person {
	// Unlike Go, So is NOT a garbage collected language.
	// You should never return a pointer to a local stack-allocated
	// variable - the stack is freed when the function returns.
	p := person{name: name}
	p.age = 42
	return p
}

func main() {
	// This syntax creates a new struct.
	printPerson(person{"Bob", 20})

	// You can name the fields when initializing a struct.
	printPerson(person{name: "Alice", age: 30})

	// Omitted fields will be zero-valued.
	printPerson(person{name: "Fred"})

	// An `&` prefix yields a pointer to the struct.
	pptr := &person{name: "Ann", age: 40}
	printPerson(*pptr)

	// It's idiomatic to encapsulate new struct creation in constructor functions.
	printPerson(newPerson("Jon"))

	// Access struct fields with a dot.
	s := person{name: "Sean", age: 50}
	println(s.name)

	// You can also use dots with struct pointers - the
	// pointers are automatically dereferenced.
	sp := &s
	println(sp.age)

	// Structs are mutable.
	sp.age = 51
	println(sp.age)

	// If a struct type is only used for a single value, we don't
	// have to give it a name. The value can have an anonymous
	// struct type.
	dog := struct {
		name   string
		isGood bool
	}{
		"Rex",
		true,
	}
	println(dog.name, dog.isGood)
}

func printPerson(p person) {
	println("{", p.name, p.age, "}")
}
