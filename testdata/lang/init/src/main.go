package main

var state int

func init() {
	state = 42
}

func main() {
	if state != 42 {
		panic("init() did not run")
	}
	println("ok")
}
