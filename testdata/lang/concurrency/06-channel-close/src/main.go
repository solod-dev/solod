package main

func producer(ch chan int) {
	for i := 0; i < 5; i++ {
		ch <- i
	}
	close(ch)
}

func main() {
	ch := make(chan int)
	go producer(ch)

	sum := 0
	count := 0
	for {
		v, ok := <-ch
		if !ok {
			break
		}
		sum += v
		count++
	}

	if sum != 10 { // 0+1+2+3+4
		panic("expected sum=10")
	}
	if count != 5 {
		panic("expected count=5")
	}
	println("ok: sum =", sum)
}
