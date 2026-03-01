package main

type Shape interface {
	Area() int
	Perim(n int) int
}

type Line interface {
	Length() int
}

type Rect struct {
	width, height int
}

func (r Rect) Area() int {
	return r.width * r.height
}

func (r Rect) Perim(n int) int {
	return n * (2*r.width + 2*r.height)
}

func (r *Rect) Length() int {
	return 2*r.width + 2*r.height
}

func calcShape(s Shape) int {
	return s.Perim(2) + s.Area()
}

func calcLine(l Line) int {
	return l.Length()
}

func shapeIsRect(s Shape) bool {
	_, ok := s.(Rect)
	return ok
}

func shapeAsRect(s Shape) Rect {
	_, ok := s.(Rect)
	if !ok {
		return Rect{}
	}
	r := s.(Rect)
	return r
}

func lineIsRect(l Line) bool {
	_, ok := l.(*Rect)
	return ok
}

func lineAsRect(l Line) *Rect {
	_, ok := l.(*Rect)
	if !ok {
		return nil
	}
	r := l.(*Rect)
	return r
}

func main() {
	r := Rect{width: 10, height: 5}
	{
		// Shape interface is implemented by Rect value.
		s := Shape(r)
		var s2 Shape = r
		var _ Shape = s2
		var s3 Shape = &r
		var _ Shape = s3

		calcShape(s)
		calcShape(Shape(r)) // also works
		calcShape(r)        // also works

		_ = shapeIsRect(s)
		rval := shapeAsRect(s)
		_ = rval
	}
	{
		// Line interface is implemented by *Rect pointer.
		l := Line(&r)
		var l2 Line = &r
		_ = l2

		calcLine(l)

		_ = lineIsRect(l)
		rptr := lineAsRect(l)
		_ = rptr
	}
}
