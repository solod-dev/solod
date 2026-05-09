package main

import "solod.dev/so/time"

var counter int

func increment(n int) {
	for i := 0; i < n; i++ {
		counter++
	}
}

func main() {
	counter = 0
	go increment(100)
	time.Sleep(time.Millisecond * 50)

	if counter != 100 {
		panic("expected counter=100")
	}
	println("ok: counter =", counter)
}
