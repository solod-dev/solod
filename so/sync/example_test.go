package sync_test

import (
	"solod.dev/so/conc"
	"solod.dev/so/mem"
	"solod.dev/so/sync"
)

type counter struct {
	mu  *sync.Mutex
	val int
}

func increment(arg any) {
	c := arg.(*counter)
	c.mu.Lock()
	c.val++
	c.mu.Unlock()
}

func ExampleMutex() {
	var mu sync.Mutex
	mu.Init()
	defer mu.Free()

	cnt := counter{mu: &mu, val: 0}
	opts := conc.PoolOpts{NumThreads: 4}
	pool := conc.NewPool(mem.System, opts)
	defer pool.Free()
	for range 100 {
		pool.Go(increment, &cnt)
	}
	pool.Wait()

	println(cnt.val)
	// 100
}

type gate struct {
	mu    sync.Mutex
	cond  sync.Cond
	ready bool
}

func await(arg any) any {
	g := arg.(*gate)
	g.mu.Lock()
	for !g.ready {
		g.cond.Wait()
	}
	println("go!")
	g.mu.Unlock()
	return nil
}

func ExampleCond() {
	var g gate
	g.mu.Init()
	g.cond.Init(&g.mu)
	defer g.mu.Free()
	defer g.cond.Free()

	thr := conc.Go(await, &g, nil)
	defer thr.Wait()

	g.mu.Lock()
	g.ready = true
	g.cond.Signal()
	g.mu.Unlock()
	// go!
}

func hello() {
	println("Hello, world!")
}

func onceHello(arg any) {
	once := arg.(*sync.Once)
	once.Do(hello)
}

func ExampleOnce() {
	var once sync.Once
	once.Init()
	defer once.Free()

	opts := conc.PoolOpts{NumThreads: 4}
	pool := conc.NewPool(mem.System, opts)
	for range 10 {
		pool.Go(onceHello, &once)
	}
	pool.Free()
	// Hello, world!
}
