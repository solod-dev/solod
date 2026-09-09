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

func main() {
	r := Rect{width: 10, height: 5}
	var s Shape = &r
	if true {
		s, ok := s.(*Rect)
		_ = s
		_ = ok
	}
}
