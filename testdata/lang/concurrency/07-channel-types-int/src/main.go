package main

func main() {
	ch := make(chan int, 2)
	ch <- 123
	ch <- 456

	a := <-ch
	b := <-ch

	if a != 123 || b != 456 {
		panic("wrong int values")
	}
	println("ok: int channel")
}
