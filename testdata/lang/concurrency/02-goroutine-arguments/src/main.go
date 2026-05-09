package main

import "solod.dev/so/time"

var result int

func compute(a int, b int, c int) {
	result = a + b*c
}

func main() {
	x := 10
	go compute(x, 20, 3) // args evaluated before goroutine starts
	x = 999              // changing x doesn't affect goroutine

	time.Sleep(time.Millisecond * 10)

	if result != 70 { // 10 + 20*3
		panic("wrong result")
	}
	println("ok: result =", result)
}
