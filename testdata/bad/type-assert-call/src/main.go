package main

type Shape interface {
	Area() int
}

type Rect struct {
	width, height int
}

func (r *Rect) Area() int {
	return r.width * r.height
}

func get(s Shape) Shape {
	return s
}

func main() {
	r := Rect{width: 10, height: 5}
	var s Shape = &r
	p, ok := get(s).(*Rect)
	_ = p
	_ = ok
}
