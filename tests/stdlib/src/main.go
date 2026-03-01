package main

import (
	"github.com/nalgeon/solod/so/errors"
	"github.com/nalgeon/solod/so/io"
	"github.com/nalgeon/solod/so/math"
)

var Err42 = errors.New("42")

func work(n int) (int, error) {
	if n == 42 {
		return 0, Err42
	}
	return 42, nil
}

func main() {
	x := math.Sqrt(4.0)
	_ = x

	r1, err := work(11)
	if err != nil {
		panic("unexpected error")
	}
	_ = r1

	r2, err := work(42)
	if err != Err42 {
		panic("expected Err42")
	}
	_ = r2

	var rdr io.Reader
	_ = rdr
}
