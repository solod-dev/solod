package main

import (
	"solod.dev/so/conc"
	"solod.dev/so/errors"
	"solod.dev/so/mem"
)

func testPool() {
	testPool_ParallelMap()
	testPool_BackPressure()
	testPool_QueueLarge()
	testPool_QueueOne()
	testPool_Error()
}

// Task carries one task's input, output and error through a *Task.
type Task struct {
	in  int
	out int
	err error
}

func square(arg any) {
	task := arg.(*Task)
	task.out = task.in * task.in
}

// Squares 0..99 in parallel and checks every result.
func testPool_ParallelMap() {
	print("- parallel map...")
	const n = 100
	tasks := make([]Task, n)
	opts := conc.PoolOpts{NumThreads: 8}
	p := conc.NewPool(mem.System, opts)
	defer p.Free()
	for i := range tasks {
		tasks[i].in = i
		p.Go(square, &tasks[i])
	}
	p.Wait()

	for i := range tasks {
		if tasks[i].out != i*i {
			panic("wrong square result")
		}
	}
	println("ok")
}

// Submits far more tasks than workers, exercising the queue-full wait.
func testPool_BackPressure() {
	print("- back-pressure...")
	const n = 1000
	tasks := make([]Task, n)
	opts := conc.PoolOpts{NumThreads: 2}
	p := conc.NewPool(mem.System, opts)
	defer p.Free()
	for i := range tasks {
		tasks[i].in = i
		p.Go(square, &tasks[i])
	}
	p.Wait()

	sum := 0
	for i := range tasks {
		sum += tasks[i].out
	}
	// Sum of i*i for i in 0..999.
	if sum != 332833500 {
		panic("wrong sum")
	}
	println("ok")
}

// Uses a queue far larger than the worker limit, so most submissions
// enqueue without blocking. All results must still be correct.
func testPool_QueueLarge() {
	print("- queue larger than workers...")
	const n = 200
	tasks := make([]Task, n)
	opts := conc.PoolOpts{NumThreads: 2, QueueSize: 128}
	p := conc.NewPool(mem.System, opts)
	defer p.Free()
	for i := range tasks {
		tasks[i].in = i
		p.Go(square, &tasks[i])
	}
	p.Wait()

	for i := range tasks {
		if tasks[i].out != i*i {
			panic("wrong square result")
		}
	}
	println("ok")
}

// Uses the smallest possible queue, so each submission past the first must
// wait for a worker to drain a slot. This stresses the queue-full
// back-pressure path with an explicit queue size.
func testPool_QueueOne() {
	print("- queue of size one...")
	const n = 50
	tasks := make([]Task, n)
	opts := conc.PoolOpts{NumThreads: 4, QueueSize: 1}
	p := conc.NewPool(mem.System, opts)
	defer p.Free()
	for i := range tasks {
		tasks[i].in = i
		p.Go(square, &tasks[i])
	}
	p.Wait()

	for i := range tasks {
		if tasks[i].out != i*i {
			panic("wrong square result")
		}
	}
	println("ok")
}

var errOddInput = errors.New("odd input")

func checkEven(arg any) {
	task := arg.(*Task)
	if task.in%2 != 0 {
		task.err = errOddInput
		return
	}
	task.out = task.in
}

// Checks that a task can report an error through its argument struct.
func testPool_Error() {
	print("- error field...")
	const n = 10
	tasks := make([]Task, n)
	opts := conc.PoolOpts{NumThreads: 4}
	p := conc.NewPool(mem.System, opts)
	defer p.Free()
	for i := range tasks {
		tasks[i].in = i
		p.Go(checkEven, &tasks[i])
	}
	p.Wait()

	for i := range tasks {
		if i%2 != 0 && tasks[i].err != errOddInput {
			panic("expected error for odd input")
		}
		if i%2 == 0 && tasks[i].err != nil {
			panic("unexpected error for even input")
		}
	}
	println("ok")
}
