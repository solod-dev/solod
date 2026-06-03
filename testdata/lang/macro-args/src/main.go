package main

type point struct {
	x, y int
}

func main() {
	// append a composite-literal value.
	pts := make([]point, 0, 2)
	pts = append(pts, point{1, 2})
	if pts[0].y != 2 {
		panic("append value")
	}

	// map with a composite-literal value.
	mv := make(map[int]point, 1)
	mv[0] = point{3, 4}
	if mv[0].x != 3 {
		panic("map value")
	}

	// map with a composite-literal key.
	mk := make(map[point]int, 1)
	mk[point{1, 2}] = 42
	if mk[point{1, 2}] != 42 {
		panic("map key")
	}
	v, ok := mk[point{1, 2}]
	if !ok || v != 42 {
		panic("map key comma-ok")
	}
}
