package main

type Data struct {
	Value int
}

func main() {
	ch := make(chan *Data, 1)
	d := &Data{Value: 99}
	ch <- d

	received := <-ch

	if received.Value != 99 {
		panic("wrong pointer value")
	}
	println("ok: pointer channel")
}
