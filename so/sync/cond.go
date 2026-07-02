package sync

// Cond is a condition variable, a rendezvous point for threads waiting for or
// announcing the occurrence of an event.
//
// Each Cond is associated with a [Mutex], which must be held when changing the
// condition and when calling [Cond.Wait]. The zero value is not usable; call
// [Cond.Init] before use. A Cond must not be copied after Init.
type Cond struct {
	cond pthread_cond_t
	mu   *Mutex
}

// Init prepares c for use, associating it with mu. It must be called exactly
// once before any other method. A Cond must not be copied after Init.
func (c *Cond) Init(mu *Mutex) {
	c.mu = mu
	rc := condInitMono(&c.cond)
	if rc != 0 {
		panic("sync: Cond.Init failed")
	}
}

// Wait atomically unlocks the associated mutex and suspends the calling thread
// until the Cond is signaled, then re-locks the mutex before returning.
// The caller must hold the mutex when calling Wait.
func (c *Cond) Wait() {
	rc := pthread_cond_wait(&c.cond, &c.mu.mu)
	if rc != 0 {
		panic("sync: Cond.Wait failed")
	}
}

// WaitFor behaves like [Cond.Wait] but stops waiting once nsec nanoseconds
// have elapsed on the monotonic clock.
//
// Measuring against the monotonic clock keeps the timeout unaffected by
// wall-clock changes (NTP steps, manual resets). To honor a fixed deadline
// across spurious wakeups, recompute the remaining time on each call.
//
// WaitFor reports whether it timed out: true means the deadline passed, false
// means the Cond was signaled. A non-positive nsec times out without blocking.
//
// As with Wait, the caller must hold the mutex and it is re-locked before returning.
func (c *Cond) WaitFor(nsec int64) bool {
	rc := condWaitRel(&c.cond, &c.mu.mu, nsec)
	if rc == eTIMEDOUT {
		return true
	}
	if rc != 0 {
		panic("sync: Cond.WaitFor failed")
	}
	return false
}

// Signal wakes at least one thread waiting on c, if any.
func (c *Cond) Signal() {
	rc := pthread_cond_signal(&c.cond)
	if rc != 0 {
		panic("sync: Cond.Signal failed")
	}
}

// Broadcast wakes all threads waiting on c, if any.
func (c *Cond) Broadcast() {
	rc := pthread_cond_broadcast(&c.cond)
	if rc != 0 {
		panic("sync: Cond.Broadcast failed")
	}
}

// Free releases the resources held by c. The Cond is unusable afterward.
func (c *Cond) Free() {
	rc := pthread_cond_destroy(&c.cond)
	if rc != 0 {
		panic("sync: Cond.Free failed")
	}
}
