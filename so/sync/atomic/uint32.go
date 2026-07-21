package atomic

// Uint32 is an atomic uint32. The zero value is zero.
// Uint32 must not be copied after first use.
type Uint32 struct {
	v uint32
}

// Load atomically loads and returns the value stored in x.
func (x *Uint32) Load() uint32 {
	return load(&x.v)
}

// Store atomically stores val into x.
func (x *Uint32) Store(val uint32) {
	store(&x.v, val)
}

// Add atomically adds delta to x and returns the new value.
func (x *Uint32) Add(delta uint32) uint32 {
	return add(&x.v, delta)
}

// Sub atomically subtracts delta from x and returns the new value.
func (x *Uint32) Sub(delta uint32) uint32 {
	return add(&x.v, ^(delta - 1))
}

// Swap atomically stores new into x and returns the previous value.
func (x *Uint32) Swap(new uint32) uint32 {
	return swap(&x.v, new)
}

// CompareAndSwap atomically sets x to new if it currently holds old,
// reporting whether the swap happened.
func (x *Uint32) CompareAndSwap(old, new uint32) bool {
	return cas(&x.v, old, new)
}
