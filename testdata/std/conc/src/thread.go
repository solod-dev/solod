package main

import (
	"solod.dev/so/conc"
	"solod.dev/so/sync"
)

func testThread() {
	testThread_Wait()
	testThread_Detach()
}

func increment(arg any) any {
	n := arg.(*int)
	*n = *n + 1
	return arg
}

// Starts a thread per element, waits for them all, and checks every result.
func testThread_Wait() {
	print("- wait...")
	const n = 16
	nums := make([]int, n)
	threads := make([]conc.Thread, n)
	for i := range nums {
		nums[i] = i
		threads[i] = conc.Go(increment, &nums[i], nil)
	}
	for i := range threads {
		res := threads[i].Wait()
		if *(res.(*int)) != i+1 {
			panic("wrong increment result")
		}
	}

	for i := range nums {
		if nums[i] != i+1 {
			panic("wrong increment result")
		}
	}
	println("ok")
}

// latch lets a detached thread report completion, since it cannot be joined.
type latch struct {
	mu   sync.Mutex
	cond sync.Cond
	done bool
	out  int
}

// squareLatch squares l.out in place, then marks the latch done.
func squareLatch(arg any) any {
	l := arg.(*latch)
	l.mu.Lock()
	l.out = l.out * l.out
	l.done = true
	l.cond.Broadcast()
	l.mu.Unlock()
	return nil
}

// Runs a task on a detached thread and waits for it through a condition.
func testThread_Detach() {
	print("- detach...")
	var l latch
	l.mu.Init()
	l.cond.Init(&l.mu)
	l.out = 9

	th := conc.Go(squareLatch, &l, nil)
	th.Detach()

	l.mu.Lock()
	for !l.done {
		l.cond.Wait()
	}
	l.mu.Unlock()

	if l.out != 81 {
		panic("wrong detached result")
	}
	l.mu.Free()
	l.cond.Free()
	println("ok")
}
