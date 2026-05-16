package sub

type Point struct {
	X, Y int
}

func MakePoint(x, y int) (Point, error) {
	return Point{X: x, Y: y}, nil
}
