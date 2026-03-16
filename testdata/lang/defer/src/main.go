package main

var state int = 0

func xopen(x *int) {
	(*x)++
}

func xclose(a any) {
	x := a.(*int)
	(*x)--
}

func funcScope() {
	xopen(&state)
	defer xclose(&state)
	if state != 1 {
		panic("unexpected state")
	}
}

func funcWithReturn() int {
	xopen(&state)
	defer xclose(&state)
	if state != 1 {
		panic("unexpected state")
	}
	return 42
}

func blockScope() {
	{
		xopen(&state)
		defer xclose(&state)
		if state != 1 {
			panic("unexpected state")
		}
	}
	if state != 0 {
		panic("unexpected state")
	}
	{
		xopen(&state)
		defer xclose(&state)
		if state != 1 {
			panic("unexpected state")
		}
	}
	if state != 0 {
		panic("unexpected state")
	}
}

func main() {
	funcScope()
	if state != 0 {
		panic("unexpected state")
	}
	funcWithReturn()
	if state != 0 {
		panic("unexpected state")
	}
	blockScope()
	if state != 0 {
		panic("unexpected state")
	}
}
