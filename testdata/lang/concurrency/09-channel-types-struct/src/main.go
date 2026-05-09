package main

type Point struct {
	X int
	Y int
}

func main() {
	ch := make(chan Point, 1)
	ch <- Point{X: 10, Y: 20}

	p := <-ch

	if p.X != 10 || p.Y != 20 {
		panic("wrong struct values")
	}
	println("ok: struct channel")
}
