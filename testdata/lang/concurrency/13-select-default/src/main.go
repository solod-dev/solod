package main

func main() {
	ch := make(chan int)

	executed := 0

	select {
	case <-ch:
		executed = 1
	default:
		executed = 2
	}

	if executed != 2 {
		panic("should execute default")
	}
	println("ok: default executed")
}
