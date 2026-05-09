package main

func main() {
	ch := make(chan string, 2)
	ch <- "hello"
	ch <- "world"

	s1 := <-ch
	s2 := <-ch

	if s1 != "hello" || s2 != "world" {
		panic("wrong string values")
	}
	println("ok: string channel")
}
