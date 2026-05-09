package main

import "solod.dev/so/time"

var sum int

func add(n int) {
	sum += n
}

func main() {
	sum = 0
	go add(10)
	go add(20)
	go add(30)

	time.Sleep(time.Millisecond * 50)

	if sum != 60 {
		panic("expected sum=60")
	}
	println("ok: sum =", sum)
}
