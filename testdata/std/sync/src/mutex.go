package main

import (
	"solod.dev/so/conc"
	"solod.dev/so/mem"
	"solod.dev/so/sync"
)

func testMutex() {
	testMutex_LockUnlock()
	testMutex_TryLock()
}

// counter is a shared count guarded by a mutex.
type counter struct {
	mu  *sync.Mutex
	val *int
}

func bump(arg any) {
	c := arg.(*counter)
	c.mu.Lock()
	*c.val = *c.val + 1
	c.mu.Unlock()
}

// Checks that no updates are lost when many workers
// concurrently increment a shared counter under a mutex.
func testMutex_LockUnlock() {
	print("- mutex...")
	const n = 1000
	var mu sync.Mutex
	mu.Init()
	val := 0
	jobs := make([]counter, n)
	opts := conc.PoolOpts{NumThreads: 8}
	p := conc.NewPool(mem.System, opts)
	for i := range jobs {
		jobs[i].mu = &mu
		jobs[i].val = &val
		p.Go(bump, &jobs[i])
	}
	p.Free()

	if val != n {
		panic("lost updates under mutex")
	}
	mu.Free()
	println("ok")
}

// Checks that TryLock acquires a free mutex and refuses
// to acquire one that is already held.
func testMutex_TryLock() {
	print("- trylock...")
	var mu sync.Mutex
	mu.Init()

	if !mu.TryLock() {
		panic("TryLock failed on free mutex")
	}
	if mu.TryLock() {
		panic("TryLock succeeded on held mutex")
	}
	mu.Unlock()

	if !mu.TryLock() {
		panic("TryLock failed after unlock")
	}
	mu.Unlock()

	mu.Free()
	println("ok")
}
