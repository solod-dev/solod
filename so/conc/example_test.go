package conc_test

import (
	"solod.dev/so/conc"
	"solod.dev/so/fmt"
	"solod.dev/so/mem"
)

// Task carries one task's input and output.
type Task struct {
	in  int
	out int
}

// square is a task function which squares
// its input and stores the result.
func square(arg any) {
	task := arg.(*Task)
	task.out = task.in * task.in
}

func ExamplePool() {
	const n = 10
	tasks := make([]Task, n)

	// Process tasks in parallel with a pool of two workers.
	opts := conc.PoolOpts{NumThreads: 2}
	pool := conc.NewPool(mem.System, opts)
	defer pool.Free()
	for i := range tasks {
		tasks[i].in = i
		pool.Go(square, &tasks[i])
	}
	pool.Wait()

	// Print results after all tasks have finished.
	for i := range tasks {
		fmt.Printf("%d squared is %d\n", tasks[i].in, tasks[i].out)
	}
}

// producer carries the channel to send on, the number of
// values to produce, and the storage backing those values.
type producer struct {
	ch   conc.Chan[int]
	n    int
	vals []int
}

// produce sends the numbers 0..n-1 on the channel, then closes it.
func produce(arg any) any {
	p := arg.(*producer)
	for i := 0; i < p.n; i++ {
		p.vals[i] = i
		p.ch.Send(&p.vals[i])
	}
	p.ch.Close()
	return nil
}

func ExampleChan() {
	const n = 5
	ch := conc.NewChan[int](mem.System, 2)
	defer ch.Free()

	// Run the producer on a single worker: it sends n
	// values into the channel, then closes it.
	prod := producer{ch: ch, n: n, vals: make([]int, n)}
	thr := conc.Go(produce, &prod, nil)
	defer thr.Wait()

	// Consume in the main thread until the channel is closed and drained.
	for {
		v, ok := ch.Recv()
		if !ok {
			break
		}
		fmt.Printf("received %d\n", *v)
	}
}
