package main

func main() {
	ch := make(chan int, 3)

	// Can send without blocking
	ch <- 10
	ch <- 20
	ch <- 30

	v1 := <-ch
	v2 := <-ch
	v3 := <-ch

	if v1 != 10 || v2 != 20 || v3 != 30 {
		panic("wrong values")
	}
	println("ok: received all values")
}
