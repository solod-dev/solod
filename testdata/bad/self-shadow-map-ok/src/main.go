package main

func main() {
	m := map[string]int{"a": 1}
	v := "a"
	if len(v) > 0 {
		v, ok := m[v]
		_ = v
		_ = ok
	}
}
