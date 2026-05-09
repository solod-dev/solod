package main

import "solod.dev/so/time"

func producer(ch chan int, count int) {
	for i := 0; i < count; i++ {
		ch <- i
	}
	close(ch)
}

func consumer(ch chan int, result *int) {
	sum := 0
	for {
		v, ok := <-ch
		if !ok {
			break
		}
		sum += v
	}
	*result = sum
}

func main() {
	ch := make(chan int, 5)
	var result int

	go producer(ch, 10)
	consumer(ch, &result)

	time.Sleep(time.Millisecond * 10)

	if result != 45 { // 0+1+2+...+9
		panic("expected sum=45")
	}
	println("ok: producer-consumer pattern")
}
