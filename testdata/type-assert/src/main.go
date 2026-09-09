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

type Circle struct {
	radius int
}

func (c *Circle) Area() int {
	return 3 * c.radius * c.radius
}

type Canvas struct {
	shape Shape
}

func main() {
	r := Rect{width: 10, height: 5}
	c := Circle{radius: 2}
	{
		// Both targets on a matching type.
		var s Shape = &r
		p, ok := s.(*Rect)
		if !ok {
			panic("want ok")
		}
		if p.Area() != 50 {
			panic("p.Area() != 50")
		}
	}
	{
		// A failed assertion gives a nil pointer.
		var s Shape = &c
		p, ok := s.(*Rect)
		if ok {
			panic("want !ok")
		}
		if p != nil {
			panic("want nil p")
		}
	}
	{
		// Blank ok target.
		var s Shape = &r
		p, _ := s.(*Rect)
		if p.Area() != 50 {
			panic("p.Area() != 50")
		}
		var s2 Shape = &c
		q, _ := s2.(*Rect)
		if q != nil {
			panic("want nil q")
		}
	}
	{
		// Blank value target.
		var s Shape = &r
		_, ok := s.(*Rect)
		if !ok {
			panic("want ok")
		}
	}
	{
		// Plain assignment.
		var s Shape = &r
		var p *Rect
		var ok bool
		p, ok = s.(*Rect)
		if !ok || p.Area() != 50 {
			panic("want ok and p.Area() == 50")
		}
		s = &c
		p, ok = s.(*Rect)
		if ok || p != nil {
			panic("want !ok and nil p")
		}
	}
	{
		// A selector operand is read two times.
		canvas := Canvas{shape: &r}
		p, ok := canvas.shape.(*Rect)
		if !ok || p.Area() != 50 {
			panic("want ok and p.Area() == 50")
		}
	}
	{
		// A redeclared target keeps its type.
		var s Shape = &r
		ok := false
		p, ok := s.(*Rect)
		if !ok || p.Area() != 50 {
			panic("want ok and p.Area() == 50")
		}
	}
	{
		// A direct assertion casts without a check.
		var s Shape = &r
		p := s.(*Rect)
		if p.Area() != 50 {
			panic("p.Area() != 50")
		}
	}
	{
		// A blank value target with a plain assignment.
		var s Shape = &r
		var ok bool
		_, ok = s.(*Rect)
		if !ok {
			panic("want ok")
		}
	}
	{
		// An assertion in the init of an if statement.
		var s Shape = &c
		if _, ok := s.(*Rect); ok {
			panic("want !ok")
		}
	}
	{
		// A nil interface matches no concrete type.
		var s Shape
		p, ok := s.(*Rect)
		if ok || p != nil {
			panic("want !ok and nil p")
		}
	}
}
