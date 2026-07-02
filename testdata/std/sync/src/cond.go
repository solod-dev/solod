package main

import (
	"solod.dev/so/conc"
	"solod.dev/so/sync"
)

// gate coordinates a single waiter with the main thread through a condition
// variable and a shared ready flag.
type gate struct {
	mu    *sync.Mutex
	cond  *sync.Cond
	ready *bool
	woke  *bool
}

func waiter(arg any) any {
	g := arg.(*gate)
	g.mu.Lock()
	for !*g.ready {
		g.cond.Wait()
	}
	*g.woke = true
	g.mu.Unlock()
	return nil
}

// Starts a worker that waits on a condition variable until main sets
// a ready flag and broadcasts, then checks the worker observed the signal.
func testCond() {
	print("- cond...")
	var mu sync.Mutex
	mu.Init()
	var cond sync.Cond
	cond.Init(&mu)
	ready := false
	woke := false

	g := gate{mu: &mu, cond: &cond, ready: &ready, woke: &woke}
	thr := conc.Go(waiter, &g, nil)

	mu.Lock()
	ready = true
	cond.Broadcast()
	mu.Unlock()

	thr.Wait()

	if !woke {
		panic("waiter did not observe signal")
	}
	cond.Free()
	mu.Free()
	println("ok")
}
