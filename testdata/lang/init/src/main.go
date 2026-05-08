package main

var state int

func init() {
	state = 42
}

type value struct{ x int }

func (v *value) init(x int) {
	v.x = x
}

func main() {
	{
		// Init function.
		if state != 42 {
			panic("init() did not run")
		}
		println("ok")
	}
	{
		// Method named init (just a regular method).
		var v value
		v.init(123)
		if v.x != 123 {
			panic("v.x != 123")
		}
	}
}
