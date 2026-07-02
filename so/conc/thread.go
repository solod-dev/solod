package conc

import "solod.dev/so/c"

// Thread is a handle to a running OS thread.
// Start a thread with [Go], then wait for completion with [Thread.Wait],
// or hand its resources to the runtime with [Thread.Detach].
type Thread struct {
	t pthread_t
}

// ThreadOpts holds optional Thread settings. A nil *ThreadOpts selects defaults.
type ThreadOpts struct {
	StackSize int // thread stack size in bytes; 0 = system default
}

// Go launches an OS thread that runs fn(arg) and returns a handle to it.
// fn is a named function and arg must point to storage that outlives the
// thread: until Wait returns, or until fn returns for a detached thread.
// opts may be nil for default settings.
//
//	var job Job
//	t := conc.Go(work, &job, nil)
//	// ... do other work concurrently ...
//	t.Wait() // job is complete once Wait returns
//
// Always either [Thread.Wait] or [Thread.Detach] a thread started with Go,
// otherwise its resources will leak.
//
// Unlike Go's goroutines, OS threads are not cheap to start. Prefer [Pool] for
// short-lived or numerous tasks; reserve Go for long-lived threads or a small,
// fixed number of threads you manage (join or detach) directly.
func Go(entry func(any) any, arg any, opts *ThreadOpts) Thread {
	stackSize := 0
	if opts != nil {
		c.Assert(opts.StackSize >= 0, "conc: stack size must be >= 0")
		stackSize = opts.StackSize
	}

	var ap *pthread_attr_t = nil
	var attr pthread_attr_t
	if stackSize > 0 {
		rc := pthread_attr_init(&attr)
		rc |= pthread_attr_setstacksize(&attr, uintptr(stackSize))
		if rc != 0 {
			panic("conc: thread attr setup failed")
		}
		ap = &attr
	}

	var th Thread
	rc := pthread_create(&th.t, ap, entry, arg)
	if rc != 0 {
		panic("conc: thread create failed")
	}

	if stackSize > 0 {
		pthread_attr_destroy(&attr)
	}
	return th
}

// Wait blocks until the thread terminates, then
// returns the thread's return value, if any.
//
// Call Wait exactly once per thread, and never on a detached thread.
func (th Thread) Wait() any {
	var ret any
	rc := pthread_join(th.t, &ret)
	if rc != 0 {
		panic("conc: thread join failed")
	}
	return ret
}

// Detach hands the thread's resources to the runtime, which reclaims them
// when the thread terminates. A detached thread must not be joined.
func (th Thread) Detach() {
	rc := pthread_detach(th.t)
	if rc != 0 {
		panic("conc: thread detach failed")
	}
}
