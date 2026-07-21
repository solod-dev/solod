package main

// counter is unexported, but so:promote emits it in the header
// so the Stats can reference it.
//
//so:promote
type counter struct {
	val int
}

// newCounter is called from an inline (header) NewStats,
// so it needs to be promoted to the header.
//
//so:promote
func newCounter() counter {
	return counter{val: 0}
}

// inc is called from an inline (header) Stats.Inc method,
// so it needs to be promoted to the header.
//
//so:promote
func (c *counter) inc() {
	c.val++
}

// Alias renames an so:promote type.
//
//so:promote
type alias counter

//so:promote
var defaultCap int = 8

//so:promote
const version = 3

// Stats is exported and has a field of the unexported so:promote type.
// Its constructor and method are emitted in the header because of so:inline.
type Stats struct {
	c counter
}

//so:inline
func NewStats() Stats {
	return Stats{c: newCounter()}
}

//so:inline
func (w *Stats) Inc() {
	w.c.inc()
}

// GetCounter is a non-inline function that returns a promoted type.
func GetCounter() counter {
	return newCounter()
}

func main() {
	w := NewStats()
	w.Inc()
	_ = alias{val: 1}
	_ = defaultCap
	_ = version
	_ = GetCounter()
}
