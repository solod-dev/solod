package main

func main() {
	ch := make(chan int, 1)

	sent := false

	select {
	case ch <- 99:
		sent = true
	}

	if !sent {
		panic("should have sent")
	}

	v := <-ch
	if v != 99 {
		panic("wrong value")
	}
	println("ok: select send")
}
