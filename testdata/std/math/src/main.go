package main

import (
	"solod.dev/so/math"
)

func main() {
	x := math.Sqrt(49)
	if x != 7 {
		panic("want x == 7")
	}
}
