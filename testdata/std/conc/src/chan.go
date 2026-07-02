package main

import (
	"solod.dev/so/conc"
	"solod.dev/so/mem"
	"solod.dev/so/time"
)

func testChan() {
	testChan_Buffered()
	testChan_ProducerConsumer()
	testChan_Unbuffered()
	testChan_UnbufferedMultiProducer()
	testChan_CloseDrain()
	testChan_TimeoutBuffered()
	testChan_TimeoutExpires()
	testChan_TimeoutHandoff()
	testChan_TimeoutSend()
}

// Fills a buffered channel without blocking
// and checks that pointers come back in FIFO order.
func testChan_Buffered() {
	print("- chan buffered...")
	vals := make([]int, 4)
	ch := conc.NewChan[int](mem.System, 4)
	for i := range vals {
		vals[i] = i * 10
		ch.Send(&vals[i])
	}
	for i := range vals {
		v, ok := ch.Recv()
		if !ok || *v != i*10 {
			panic("wrong buffered value")
		}
	}
	ch.Free()
	println("ok")
}

// sumTask carries a channel and the resulting sum between threads.
type sumTask struct {
	ch  conc.Chan[int]
	sum int
}

// consume receives pointers until the channel is closed and accumulates them.
func consume(arg any) any {
	task := arg.(*sumTask)
	for {
		v, ok := task.ch.Recv()
		if !ok {
			break
		}
		task.sum += *v
	}
	return nil
}

// Sends 0..n-1 from the main thread through a small buffered channel
// while a worker thread sums them, exercising back-pressure.
func testChan_ProducerConsumer() {
	print("- chan producer/consumer...")
	const n = 1000
	nums := make([]int, n)
	task := sumTask{ch: conc.NewChan[int](mem.System, 8), sum: 0}

	thr := conc.Go(consume, &task, nil)
	for i := range nums {
		nums[i] = i
		task.ch.Send(&nums[i])
	}
	task.ch.Close()
	thr.Wait()

	// Sum of 0..999.
	if task.sum != 499500 {
		panic("wrong producer/consumer sum")
	}
	task.ch.Free()
	println("ok")
}

// seqTask for sending a sequence of values to a channel.
type seqTask struct {
	ch   conc.Chan[int]
	nums []int
}

// produceSeq sends 0..9 to the channel and then closes it.
func produceSeq(arg any) any {
	task := arg.(*seqTask)
	for i := range task.nums {
		task.nums[i] = i
		task.ch.Send(&task.nums[i])
	}
	task.ch.Close()
	return nil
}

// Receives from an unbuffered channel fed by a worker thread
// and checks the handoff order.
func testChan_Unbuffered() {
	print("- chan unbuffered...")
	task := seqTask{ch: conc.NewChan[int](mem.System, 0), nums: make([]int, 10)}

	want := 0
	thr := conc.Go(produceSeq, &task, nil)
	for {
		v, ok := task.ch.Recv()
		if !ok {
			break
		}
		if *v != want {
			panic("wrong unbuffered handoff order")
		}
		want++
	}
	thr.Wait()

	if want != 10 {
		panic("missing unbuffered values")
	}
	task.ch.Free()
	println("ok")
}

// rangeTask for sending a range of values to a channel.
type rangeTask struct {
	ch   conc.Chan[int]
	base int
	n    int
	vals []int
}

// produceRange sends base..base+n-1 to the channel.
func produceRange(arg any) {
	task := arg.(*rangeTask)
	for i := 0; i < task.n; i++ {
		task.vals[i] = task.base + i
		task.ch.Send(&task.vals[i])
	}
}

// Runs several producer threads sending on a single unbuffered channel while
// the main thread receives. Each value 0..N-1 is sent exactly once across
// producers; the receiver checks none is lost or duplicated. This exercises
// the rendezvous handshake with concurrent senders.
func testChan_UnbufferedMultiProducer() {
	print("- chan unbuffered multi-producer...")
	const producers = 4
	const perProducer = 250
	const total = producers * perProducer

	ch := conc.NewChan[int](mem.System, 0)
	opts := conc.PoolOpts{NumThreads: producers}
	p := conc.NewPool(mem.System, opts)

	tasks := make([]rangeTask, producers)
	for i := range tasks {
		tasks[i] = rangeTask{ch: ch, base: i * perProducer, n: perProducer, vals: make([]int, perProducer)}
		p.Go(produceRange, &tasks[i])
	}

	seen := make([]bool, total)
	for range total {
		v, ok := ch.Recv()
		if !ok {
			panic("unexpected close")
		}
		if *v < 0 || *v >= total || seen[*v] {
			panic("lost or duplicated unbuffered value")
		}
		seen[*v] = true
	}
	p.Free()
	ch.Free()
	println("ok")
}

// Checks that buffered values survive Close and are drained in order
// before Recv reports the channel closed.
func testChan_CloseDrain() {
	print("- chan close drain...")
	vals := []int{1, 2, 3}
	ch := conc.NewChan[int](mem.System, 4)
	for i := range vals {
		ch.Send(&vals[i])
	}
	ch.Close()

	seen := 0
	want := 1
	for {
		v, ok := ch.Recv()
		if !ok {
			break
		}
		if *v != want {
			panic("wrong drained value")
		}
		want++
		seen++
	}
	if seen != 3 {
		panic("did not drain all buffered values")
	}
	ch.Free()
	println("ok")
}

// Exercises non-blocking SendTimeout/RecvTimeout (d == 0) on a buffered channel
// from a single thread, where the outcomes are fully deterministic: sends fail
// once full, receives fail once empty, and a drained closed channel reports
// Closed.
func testChan_TimeoutBuffered() {
	print("- chan timeout buffered...")
	vals := []int{10, 20, 30}
	ch := conc.NewChan[int](mem.System, 2)

	// The buffer holds 2; the third non-blocking send must time out.
	if ch.SendTimeout(&vals[0], 0) != conc.Ok || ch.SendTimeout(&vals[1], 0) != conc.Ok {
		panic("SendTimeout should succeed with room")
	}
	if ch.SendTimeout(&vals[2], 0) != conc.Timeout {
		panic("SendTimeout should time out when full")
	}

	// Drain in FIFO order, then a non-blocking receive must time out.
	v, st := ch.RecvTimeout(0)
	if st != conc.Ok || *v != 10 {
		panic("wrong first RecvTimeout value")
	}
	v, st = ch.RecvTimeout(0)
	if st != conc.Ok || *v != 20 {
		panic("wrong second RecvTimeout value")
	}
	if _, st = ch.RecvTimeout(0); st != conc.Timeout {
		panic("RecvTimeout should time out when empty")
	}

	// After close with no buffered values, a receive reports Closed.
	ch.Close()
	if _, st = ch.RecvTimeout(0); st != conc.Closed {
		panic("RecvTimeout should report Closed")
	}
	ch.Free()
	println("ok")
}

// Checks that timed operations actually give up at the deadline when no peer
// ever appears: both a send and a receive on an idle unbuffered channel must
// return Timeout rather than block forever.
func testChan_TimeoutExpires() {
	print("- chan timeout expires...")
	ch := conc.NewChan[int](mem.System, 0)

	x := 1
	if ch.SendTimeout(&x, 10*time.Millisecond) != conc.Timeout {
		panic("SendTimeout should time out with no receiver")
	}
	if _, st := ch.RecvTimeout(10 * time.Millisecond); st != conc.Timeout {
		panic("RecvTimeout should time out with no sender")
	}
	ch.Free()
	println("ok")
}

// Receives from an unbuffered channel with a deadline while a worker thread
// feeds it with blocking sends. The loop tolerates timeouts and stops on
// Closed, checking the handoff order.
func testChan_TimeoutHandoff() {
	print("- chan timeout handoff...")
	task := seqTask{ch: conc.NewChan[int](mem.System, 0), nums: make([]int, 10)}

	thr := conc.Go(produceSeq, &task, nil)
	want := 0
	for {
		v, st := task.ch.RecvTimeout(50 * time.Millisecond)
		if st == conc.Closed {
			break
		}
		if st == conc.Timeout {
			continue // no sender ready yet; keep polling
		}
		if *v != want {
			panic("wrong timeout handoff order")
		}
		want++
	}
	thr.Wait()

	if want != 10 {
		panic("missing timeout handoff values")
	}
	task.ch.Free()
	println("ok")
}

// Sends on an unbuffered channel with a deadline while a worker thread drains
// it with blocking receives. Each send retries until a receiver takes it.
func testChan_TimeoutSend() {
	print("- chan timeout send...")
	const n = 100
	nums := make([]int, n)
	task := sumTask{ch: conc.NewChan[int](mem.System, 0), sum: 0}

	thr := conc.Go(consume, &task, nil)
	for i := range nums {
		nums[i] = i
		for task.ch.SendTimeout(&nums[i], 50*time.Millisecond) != conc.Ok {
			// No receiver ready yet; keep retrying.
		}
	}
	task.ch.Close()
	thr.Wait()

	// Sum of 0..99.
	if task.sum != 4950 {
		panic("wrong timeout send sum")
	}
	task.ch.Free()
	println("ok")
}
