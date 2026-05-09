package main

func main() {
	ch1 := make(chan int, 1)
	ch2 := make(chan int, 1)

	ch1 <- 42

	selected := 0
	value := 0

	select {
	case v := <-ch1:
		selected = 1
		value = v
	case v := <-ch2:
		selected = 2
		value = v
	}

	if selected != 1 || value != 42 {
		panic("wrong selection")
	}
	println("ok: selected ch1")
}
