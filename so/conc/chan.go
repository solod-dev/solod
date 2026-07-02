package conc

import (
	"solod.dev/so/c"
	"solod.dev/so/mem"
	"solod.dev/so/time"
)

// Chan is a thread-safe FIFO channel, similar to Go's built-in `chan T`.
// It carries pointers to T: senders hand off ownership of an allocated value
// and receivers take it.
//
// It supports buffered mode (created with n > 0) and unbuffered rendezvous mode
// (n == 0), where each send blocks until a receiver takes the value. Exactly one
// of the two backing engines is non-nil, chosen at creation time.
type Chan[T any] struct {
	buf *Buffer     // non-nil for buffered channels (n > 0)
	rdv *Rendezvous // non-nil for unbuffered channels (n == 0)
}

// NewChan creates a channel of *T backed by alloc. n is the buffer size:
// n > 0 makes it buffered, n == 0 makes it an unbuffered rendezvous channel.
// Call [Chan.Free] exactly once when done.
//
//so:inline
func NewChan[T any](alloc mem.Allocator, n int) Chan[T] {
	_n := n
	var _ch Chan[T]
	c.Assert(_n >= 0, "conc: chan size must be >= 0")
	if _n > 0 {
		_ch.buf = NewBuffer(alloc, _n)
	} else {
		_ch.rdv = NewRendezvous(alloc)
	}
	return _ch
}

// Send sends v on the channel, blocking until there is room (buffered) or a
// receiver takes it (unbuffered). The value pointed to by v must outlive the
// handoff. Sending on a closed channel panics.
//
// Send is thread-safe.
//
//so:inline
func (ch *Chan[T]) Send(v *T) {
	if ch.buf != nil {
		ch.buf.Send(any(v))
	} else {
		ch.rdv.Send(any(v))
	}
}

// SendTimeout sends v, waiting up to d for room (buffered) or a receiver
// (unbuffered). A zero or negative d makes it non-blocking. The value pointed
// to by v must outlive the handoff.
//
// Returns [Ok] if the value was sent, [Timeout] if the deadline passed first,
// or [Closed] if the channel is closed. Unlike [Chan.Send], it does not panic
// on a closed channel.
//
// A non-blocking send on an unbuffered channel reports [Timeout] unless a
// receiver is already parked and takes the value in the brief window it is offered.
//
// SendTimeout is thread-safe.
//
//so:inline
func (ch *Chan[T]) SendTimeout(v *T, d time.Duration) Status {
	var _st Status
	if ch.buf != nil {
		_st = ch.buf.SendTimeout(any(v), d)
	} else {
		_st = ch.rdv.SendTimeout(any(v), d)
	}
	return _st
}

// Recv receives a pointer from the channel. The bool is false when the channel
// is closed and no buffered values remain, in which case the pointer is nil.
//
// Recv is thread-safe.
//
//so:inline
func (ch *Chan[T]) Recv() (*T, bool) {
	var _v any
	var _ok bool
	if ch.buf != nil {
		_v, _ok = ch.buf.Recv()
	} else {
		_v, _ok = ch.rdv.Recv()
	}
	return _v.(*T), _ok
}

// RecvTimeout receives a pointer, waiting up to d for a value. A zero or
// negative d makes it non-blocking. Returns the pointer with [Ok], nil with
// [Timeout] if the deadline passed first, or nil with [Closed] if the channel
// is closed and drained.
//
// RecvTimeout is thread-safe.
//
//so:inline
func (ch *Chan[T]) RecvTimeout(d time.Duration) (*T, Status) {
	var _v any
	var _st Status
	if ch.buf != nil {
		_v, _st = ch.buf.RecvTimeout(d)
	} else {
		_v, _st = ch.rdv.RecvTimeout(d)
	}
	return _v.(*T), _st
}

// Close closes the channel. Subsequent sends panic; receivers drain remaining
// buffered values and then return (nil, false). Closing a closed channel panics.
//
// Close is thread-safe and may run concurrently with Send and Recv but it
// must be called exactly once; a repeated or concurrent Close panics.
//
//so:inline
func (ch *Chan[T]) Close() {
	if ch.buf != nil {
		ch.buf.Close()
	} else {
		ch.rdv.Close()
	}
}

// Free releases the channel's resources. The channel is unusable afterward.
// Free should only be called once; it's not thread-safe.
//
//so:inline
func (ch *Chan[T]) Free() {
	if ch.buf != nil {
		ch.buf.Free()
	} else {
		ch.rdv.Free()
	}
}
