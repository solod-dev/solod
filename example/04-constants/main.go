// So supports _constants_ of character, string, boolean,
// and numeric values.
package main

import (
	"github.com/nalgeon/solod/so/math"
)

// `const` declares a constant value.
const s string = "constant"

func main() {
	println(s)

	// A `const` statement can also appear inside a
	// function body.
	const n = 500000000

	// Constant expressions perform arithmetic with
	// arbitrary precision.
	const d = 3e20 / n
	println(d)

	// A numeric constant has no type until it's given
	// one, such as by an explicit conversion.
	println(int64(d))

	// A number can be given a type by using it in a
	// context that requires one, such as a variable
	// assignment or function call. For example, here
	// `math.Sqrt` expects a `float64`.
	println(math.Sqrt(n))
}
