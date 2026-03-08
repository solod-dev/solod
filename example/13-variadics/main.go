// _Variadic functions_ can be called with any number
// of trailing arguments. For example, `println`
// is a common variadic function.
package main

// Here's a function that will take an arbitrary number
// of `int`s as arguments.
func sum(nums ...int) {
	print("sum = ")
	total := 0
	// Within the function, the type of `nums` is
	// equivalent to `[]int`. We can call `len(nums)`,
	// iterate over it with `range`, etc.
	for _, num := range nums {
		total += num
	}
	println(total)
}

func main() {
	// Variadic functions can be called in the usual way
	// with individual arguments.
	sum(1, 2)
	sum(1, 2, 3)

	// If you already have multiple args in a slice,
	// apply them to a variadic function using
	// `func(slice...)` like this.
	nums := []int{1, 2, 3, 4}
	sum(nums...)
}
