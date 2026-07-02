package conc

import (
	"solod.dev/so/mem"
	"solod.dev/so/sync"
	"solod.dev/so/time"
)

// Rendezvous is the non-generic engine behind an unbuffered [Chan]:
// a thread-safe handoff with no buffer. Each send blocks until a receiver
// takes the value; the pointer is handed straight from the sender to the
// receiver with no staging buffer. In most cases, using [Chan] is more convenient.
//
// The single handoff slot is owned by the sender for its whole lifetime: the
// sender sets full on publish and clears it only after observing claimed, so no
// other sender can enter (and reset claimed) while a sender is mid-handoff.
// This keeps the handshake unambiguous with any number of senders.
type Rendezvous struct {
	alloc mem.Allocator

	mu   sync.Mutex
	cond sync.Cond // broadcast on every slot state change

	src any // the sender's published value (valid while full)

	full    bool // a sender has published src and not yet freed the slot
	claimed bool // the current value has been taken by a receiver
	closed  bool // true after Close
}

// NewRendezvous creates an unbuffered channel handing off pointers.
func NewRendezvous(alloc mem.Allocator) *Rendezvous {
	ch := mem.Alloc[Rendezvous](alloc)
	ch.alloc = alloc
	ch.src = nil
	ch.full, ch.claimed, ch.closed = false, false, false

	ch.mu.Init()
	ch.cond.Init(&ch.mu)
	return ch
}

// SendTimeout publishes v and waits up to d for a receiver to take it. A zero
// or negative d makes it non-blocking: the value is offered only for the instant
// before the deadline, so the send reports Timeout unless a receiver is already
// parked and takes it in that window.
//
// Returns Ok if a receiver took the value, Timeout if the deadline passed first,
// or Closed if the channel is closed.
func (ch *Rendezvous) SendTimeout(v any, d time.Duration) Status {
	deadline := time.Now().Add(d)
	ch.mu.Lock()
	// Wait until the handoff slot is free (the previous sender finished).
	timedOut := false
	for ch.full && !ch.closed && !timedOut {
		dur := int64(time.Until(deadline))
		timedOut = ch.cond.WaitFor(dur)
	}
	if ch.closed {
		ch.mu.Unlock()
		return Closed
	}
	if ch.full {
		// Still occupied: the deadline passed before the slot freed.
		ch.mu.Unlock()
		return Timeout
	}
	// Publish the value and wake a receiver.
	ch.src = v
	ch.full = true
	ch.claimed = false
	ch.cond.Broadcast()
	// Wait for a receiver to take the value, up to the deadline.
	timedOut = false
	for !ch.claimed && !ch.closed && !timedOut {
		dur := int64(time.Until(deadline))
		timedOut = ch.cond.WaitFor(dur)
	}
	// A receiver may have claimed right at the deadline, so a claim wins over a
	// timeout or close. Free the slot whichever way the handoff ended.
	claimed, closed := ch.claimed, ch.closed
	ch.full = false
	ch.src = nil
	ch.cond.Broadcast()
	ch.mu.Unlock()
	if claimed {
		return Ok
	}
	if closed {
		return Closed
	}
	return Timeout
}

// Send publishes v and blocks until a receiver takes it. Panics if the
// channel is closed before the handoff completes.
func (ch *Rendezvous) Send(v any) {
	ch.mu.Lock()
	// Wait until the handoff slot is free (the previous sender finished).
	for ch.full && !ch.closed {
		ch.cond.Wait()
	}
	if ch.closed {
		ch.mu.Unlock()
		panic("conc: send on closed channel")
	}
	// Publish the value and wake a receiver.
	ch.src = v
	ch.full = true
	ch.claimed = false
	ch.cond.Broadcast()
	// Wait for a receiver to take the value.
	for !ch.claimed && !ch.closed {
		ch.cond.Wait()
	}
	// Free the slot whether the handoff completed or the channel closed.
	done := ch.claimed
	ch.full = false
	ch.src = nil
	ch.cond.Broadcast()
	ch.mu.Unlock()
	if !done {
		// Closed before any receiver took the value.
		panic("conc: send on closed channel")
	}
}

// Recv takes the handed-off value. It reports whether a value was received:
// false means the channel is closed with no pending value.
func (ch *Rendezvous) Recv() (any, bool) {
	ch.mu.Lock()
	// Wait for a published, not-yet-claimed value.
	for (!ch.full || ch.claimed) && !ch.closed {
		ch.cond.Wait()
	}
	if !ch.full || ch.claimed {
		// Closed with no value to take.
		ch.mu.Unlock()
		return nil, false
	}
	v := ch.src
	ch.claimed = true
	ch.cond.Broadcast()
	ch.mu.Unlock()
	return v, true
}

// RecvTimeout takes a published value, waiting up to d for a sender to publish
// one. A zero or negative d makes it non-blocking: it succeeds only if a sender
// is already parked on the handoff.
//
// Returns the value with Ok, nil with Timeout if the deadline passed first,
// or nil with Closed if the channel is closed with no pending value.
func (ch *Rendezvous) RecvTimeout(d time.Duration) (any, Status) {
	deadline := time.Now().Add(d)
	ch.mu.Lock()
	timedOut := false
	for (!ch.full || ch.claimed) && !ch.closed && !timedOut {
		dur := int64(time.Until(deadline))
		timedOut = ch.cond.WaitFor(dur)
	}
	if ch.full && !ch.claimed {
		// A value is available (possibly published right at the deadline), so
		// it wins over both close and timeout.
		v := ch.src
		ch.claimed = true
		ch.cond.Broadcast()
		ch.mu.Unlock()
		return v, Ok
	}
	closed := ch.closed
	ch.mu.Unlock()
	if closed {
		return nil, Closed
	}
	return nil, Timeout
}

// Close marks the channel closed. A sender blocked on the handoff panics; a
// receiver with no pending value returns false. Closing a closed channel panics.
func (ch *Rendezvous) Close() {
	ch.mu.Lock()
	if ch.closed {
		ch.mu.Unlock()
		panic("conc: close of closed channel")
	}
	ch.closed = true
	ch.cond.Broadcast()
	ch.mu.Unlock()
}

// Free releases the channel's resources. The channel is unusable afterward.
func (ch *Rendezvous) Free() {
	ch.mu.Free()
	ch.cond.Free()
	mem.Free(ch.alloc, ch)
}
