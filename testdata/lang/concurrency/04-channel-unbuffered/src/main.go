package main

var received int

func sender(ch chan int) {
	ch <- 42
}

func main() {
	ch := make(chan int)
	go sender(ch)
	received = <-ch

	if received != 42 {
		panic("expected 42")
	}
	println("ok: received =", received)
}
