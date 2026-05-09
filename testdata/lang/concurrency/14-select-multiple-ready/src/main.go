package main

func main() {
	ch1 := make(chan int, 1)
	ch2 := make(chan int, 1)

	ch1 <- 10
	ch2 <- 20

	sum := 0

	// Both ready - randomly selects
	select {
	case v := <-ch1:
		sum += v
	case v := <-ch2:
		sum += v
	}

	// One of them was consumed
	if sum != 10 && sum != 20 {
		panic("wrong sum")
	}
	println("ok: selected one of two ready channels")
}
