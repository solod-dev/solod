#include "main.h"

// -- Types --

typedef struct sumTask sumTask;
typedef struct seqTask seqTask;
typedef struct rangeTask rangeTask;
typedef struct latch latch;

// sumTask carries a channel and the resulting sum between threads.
typedef struct sumTask {
    conc_Chan ch;
    so_int sum;
} sumTask;

// seqTask for sending a sequence of values to a channel.
typedef struct seqTask {
    conc_Chan ch;
    so_Slice nums;
} seqTask;

// rangeTask for sending a range of values to a channel.
typedef struct rangeTask {
    conc_Chan ch;
    so_int base;
    so_int n;
    so_Slice vals;
} rangeTask;

// latch lets a detached thread report completion, since it cannot be joined.
typedef struct latch {
    sync_Mutex mu;
    sync_Cond cond;
    bool done;
    so_int out;
} latch;

// -- Forward declarations --
static void testChan(void);
static void testChan_Buffered(void);
static void* consume(void* arg);
static void testChan_ProducerConsumer(void);
static void* produceSeq(void* arg);
static void testChan_Unbuffered(void);
static void produceRange(void* arg);
static void testChan_UnbufferedMultiProducer(void);
static void testChan_CloseDrain(void);
static void testChan_TimeoutBuffered(void);
static void testChan_TimeoutExpires(void);
static void testChan_TimeoutHandoff(void);
static void testChan_TimeoutSend(void);
static void testPool(void);
static void square(void* arg);
static void testPool_ParallelMap(void);
static void testPool_BackPressure(void);
static void testPool_QueueLarge(void);
static void testPool_QueueOne(void);
static void checkEven(void* arg);
static void testPool_Error(void);
static void testThread(void);
static void* increment(void* arg);
static void testThread_Wait(void);
static void* squareLatch(void* arg);
static void testThread_Detach(void);

// -- Variables and constants --
static so_Error errOddInput = errors_New("odd input");

// -- chan.go --

static void testChan(void) {
    testChan_Buffered();
    testChan_ProducerConsumer();
    testChan_Unbuffered();
    testChan_UnbufferedMultiProducer();
    testChan_CloseDrain();
    testChan_TimeoutBuffered();
    testChan_TimeoutExpires();
    testChan_TimeoutHandoff();
    testChan_TimeoutSend();
}

// Fills a buffered channel without blocking
// and checks that pointers come back in FIFO order.
static void testChan_Buffered(void) {
    so_print("%s", "- chan buffered...");
    so_Slice vals = so_make_slice(so_int, 4, 4);
    conc_Chan ch = conc_NewChan(so_int, (mem_System), (4));
    for (so_int i = 0; i < so_len(vals); i++) {
        so_at(so_int, vals, i) = i * 10;
        conc_Chan_Send(so_int, (&ch), (&so_at(so_int, vals, i)));
    }
    for (so_int i = 0; i < so_len(vals); i++) {
        so_R_ptr_bool _res1 = conc_Chan_Recv(so_int, (&ch));
        so_int* v = _res1.val;
        bool ok = _res1.val2;
        if (!ok || *v != i * 10) {
            so_panic("wrong buffered value");
        }
    }
    conc_Chan_Free(so_int, (&ch));
    so_println("%s", "ok");
}

// consume receives pointers until the channel is closed and accumulates them.
static void* consume(void* arg) {
    sumTask* task = (sumTask*)arg;
    for (;;) {
        so_R_ptr_bool _res1 = conc_Chan_Recv(so_int, (&task->ch));
        so_int* v = _res1.val;
        bool ok = _res1.val2;
        if (!ok) {
            break;
        }
        task->sum += *v;
    }
    return NULL;
}

// Sends 0..n-1 from the main thread through a small buffered channel
// while a worker thread sums them, exercising back-pressure.
static void testChan_ProducerConsumer(void) {
    so_print("%s", "- chan producer/consumer...");
    const int64_t n = 1000;
    so_Slice nums = so_make_slice(so_int, n, n);
    sumTask task = (sumTask){.ch = conc_NewChan(so_int, (mem_System), (8)), .sum = 0};
    conc_Thread thr = conc_Go(consume, &task, NULL);
    for (so_int i = 0; i < so_len(nums); i++) {
        so_at(so_int, nums, i) = i;
        conc_Chan_Send(so_int, (&task.ch), (&so_at(so_int, nums, i)));
    }
    conc_Chan_Close(so_int, (&task.ch));
    conc_Thread_Wait(thr);
    // Sum of 0..999.
    if (task.sum != 499500) {
        so_panic("wrong producer/consumer sum");
    }
    conc_Chan_Free(so_int, (&task.ch));
    so_println("%s", "ok");
}

// produceSeq sends 0..9 to the channel and then closes it.
static void* produceSeq(void* arg) {
    seqTask* task = (seqTask*)arg;
    for (so_int i = 0; i < so_len(task->nums); i++) {
        so_at(so_int, task->nums, i) = i;
        conc_Chan_Send(so_int, (&task->ch), (&so_at(so_int, task->nums, i)));
    }
    conc_Chan_Close(so_int, (&task->ch));
    return NULL;
}

// Receives from an unbuffered channel fed by a worker thread
// and checks the handoff order.
static void testChan_Unbuffered(void) {
    so_print("%s", "- chan unbuffered...");
    seqTask task = (seqTask){.ch = conc_NewChan(so_int, (mem_System), (0)), .nums = so_make_slice(so_int, 10, 10)};
    so_int want = 0;
    conc_Thread thr = conc_Go(produceSeq, &task, NULL);
    for (;;) {
        so_R_ptr_bool _res1 = conc_Chan_Recv(so_int, (&task.ch));
        so_int* v = _res1.val;
        bool ok = _res1.val2;
        if (!ok) {
            break;
        }
        if (*v != want) {
            so_panic("wrong unbuffered handoff order");
        }
        want++;
    }
    conc_Thread_Wait(thr);
    if (want != 10) {
        so_panic("missing unbuffered values");
    }
    conc_Chan_Free(so_int, (&task.ch));
    so_println("%s", "ok");
}

// produceRange sends base..base+n-1 to the channel.
static void produceRange(void* arg) {
    rangeTask* task = (rangeTask*)arg;
    for (so_int i = 0; i < task->n; i++) {
        so_at(so_int, task->vals, i) = task->base + i;
        conc_Chan_Send(so_int, (&task->ch), (&so_at(so_int, task->vals, i)));
    }
}

// Runs several producer threads sending on a single unbuffered channel while
// the main thread receives. Each value 0..N-1 is sent exactly once across
// producers; the receiver checks none is lost or duplicated. This exercises
// the rendezvous handshake with concurrent senders.
static void testChan_UnbufferedMultiProducer(void) {
    so_print("%s", "- chan unbuffered multi-producer...");
    const int64_t producers = 4;
    const int64_t perProducer = 250;
    const int64_t total = producers * perProducer;
    conc_Chan ch = conc_NewChan(so_int, (mem_System), (0));
    conc_PoolOpts opts = (conc_PoolOpts){.NumThreads = producers};
    conc_Pool* p = conc_NewPool(mem_System, opts);
    so_Slice tasks = so_make_slice(rangeTask, producers, producers);
    for (so_int i = 0; i < so_len(tasks); i++) {
        so_at(rangeTask, tasks, i) = (rangeTask){.ch = ch, .base = i * perProducer, .n = perProducer, .vals = so_make_slice(so_int, perProducer, perProducer)};
        conc_Pool_Go(p, produceRange, &so_at(rangeTask, tasks, i));
    }
    so_Slice seen = so_make_slice(bool, total, total);
    for (so_int _i = 0; _i < total; _i++) {
        so_R_ptr_bool _res1 = conc_Chan_Recv(so_int, (&ch));
        so_int* v = _res1.val;
        bool ok = _res1.val2;
        if (!ok) {
            so_panic("unexpected close");
        }
        if (*v < 0 || *v >= total || so_at(bool, seen, *v)) {
            so_panic("lost or duplicated unbuffered value");
        }
        so_at(bool, seen, *v) = true;
    }
    conc_Pool_Free(p);
    conc_Chan_Free(so_int, (&ch));
    so_println("%s", "ok");
}

// Checks that buffered values survive Close and are drained in order
// before Recv reports the channel closed.
static void testChan_CloseDrain(void) {
    so_print("%s", "- chan close drain...");
    so_Slice vals = (so_Slice){(so_int[3]){1, 2, 3}, 3, 3};
    conc_Chan ch = conc_NewChan(so_int, (mem_System), (4));
    for (so_int i = 0; i < so_len(vals); i++) {
        conc_Chan_Send(so_int, (&ch), (&so_at(so_int, vals, i)));
    }
    conc_Chan_Close(so_int, (&ch));
    so_int seen = 0;
    so_int want = 1;
    for (;;) {
        so_R_ptr_bool _res1 = conc_Chan_Recv(so_int, (&ch));
        so_int* v = _res1.val;
        bool ok = _res1.val2;
        if (!ok) {
            break;
        }
        if (*v != want) {
            so_panic("wrong drained value");
        }
        want++;
        seen++;
    }
    if (seen != 3) {
        so_panic("did not drain all buffered values");
    }
    conc_Chan_Free(so_int, (&ch));
    so_println("%s", "ok");
}

// Exercises non-blocking SendTimeout/RecvTimeout (d == 0) on a buffered channel
// from a single thread, where the outcomes are fully deterministic: sends fail
// once full, receives fail once empty, and a drained closed channel reports
// Closed.
static void testChan_TimeoutBuffered(void) {
    so_print("%s", "- chan timeout buffered...");
    so_Slice vals = (so_Slice){(so_int[3]){10, 20, 30}, 3, 3};
    conc_Chan ch = conc_NewChan(so_int, (mem_System), (2));
    // The buffer holds 2; the third non-blocking send must time out.
    if (conc_Chan_SendTimeout(so_int, (&ch), (&so_at(so_int, vals, 0)), (0)) != conc_Ok || conc_Chan_SendTimeout(so_int, (&ch), (&so_at(so_int, vals, 1)), (0)) != conc_Ok) {
        so_panic("SendTimeout should succeed with room");
    }
    if (conc_Chan_SendTimeout(so_int, (&ch), (&so_at(so_int, vals, 2)), (0)) != conc_Timeout) {
        so_panic("SendTimeout should time out when full");
    }
    // Drain in FIFO order, then a non-blocking receive must time out.
    so_R_ptr_int _res1 = conc_Chan_RecvTimeout(so_int, (&ch), (0));
    so_int* v = _res1.val;
    conc_Status st = _res1.val2;
    if (st != conc_Ok || *v != 10) {
        so_panic("wrong first RecvTimeout value");
    }
    so_R_ptr_int _res2 = conc_Chan_RecvTimeout(so_int, (&ch), (0));
    v = _res2.val;
    st = _res2.val2;
    if (st != conc_Ok || *v != 20) {
        so_panic("wrong second RecvTimeout value");
    }
    {
        so_R_ptr_int _res3 = conc_Chan_RecvTimeout(so_int, (&ch), (0));
        st = _res3.val2;
        if (st != conc_Timeout) {
            so_panic("RecvTimeout should time out when empty");
        }
    }
    // After close with no buffered values, a receive reports Closed.
    conc_Chan_Close(so_int, (&ch));
    {
        so_R_ptr_int _res4 = conc_Chan_RecvTimeout(so_int, (&ch), (0));
        st = _res4.val2;
        if (st != conc_Closed) {
            so_panic("RecvTimeout should report Closed");
        }
    }
    conc_Chan_Free(so_int, (&ch));
    so_println("%s", "ok");
}

// Checks that timed operations actually give up at the deadline when no peer
// ever appears: both a send and a receive on an idle unbuffered channel must
// return Timeout rather than block forever.
static void testChan_TimeoutExpires(void) {
    so_print("%s", "- chan timeout expires...");
    conc_Chan ch = conc_NewChan(so_int, (mem_System), (0));
    so_int x = 1;
    if (conc_Chan_SendTimeout(so_int, (&ch), (&x), (10 * time_Millisecond)) != conc_Timeout) {
        so_panic("SendTimeout should time out with no receiver");
    }
    {
        so_R_ptr_int _res1 = conc_Chan_RecvTimeout(so_int, (&ch), (10 * time_Millisecond));
        conc_Status st = _res1.val2;
        if (st != conc_Timeout) {
            so_panic("RecvTimeout should time out with no sender");
        }
    }
    conc_Chan_Free(so_int, (&ch));
    so_println("%s", "ok");
}

// Receives from an unbuffered channel with a deadline while a worker thread
// feeds it with blocking sends. The loop tolerates timeouts and stops on
// Closed, checking the handoff order.
static void testChan_TimeoutHandoff(void) {
    so_print("%s", "- chan timeout handoff...");
    seqTask task = (seqTask){.ch = conc_NewChan(so_int, (mem_System), (0)), .nums = so_make_slice(so_int, 10, 10)};
    conc_Thread thr = conc_Go(produceSeq, &task, NULL);
    so_int want = 0;
    for (;;) {
        so_R_ptr_int _res1 = conc_Chan_RecvTimeout(so_int, (&task.ch), (50 * time_Millisecond));
        so_int* v = _res1.val;
        conc_Status st = _res1.val2;
        if (st == conc_Closed) {
            break;
        }
        if (st == conc_Timeout) {
            // no sender ready yet; keep polling
            continue;
        }
        if (*v != want) {
            so_panic("wrong timeout handoff order");
        }
        want++;
    }
    conc_Thread_Wait(thr);
    if (want != 10) {
        so_panic("missing timeout handoff values");
    }
    conc_Chan_Free(so_int, (&task.ch));
    so_println("%s", "ok");
}

// Sends on an unbuffered channel with a deadline while a worker thread drains
// it with blocking receives. Each send retries until a receiver takes it.
static void testChan_TimeoutSend(void) {
    so_print("%s", "- chan timeout send...");
    const int64_t n = 100;
    so_Slice nums = so_make_slice(so_int, n, n);
    sumTask task = (sumTask){.ch = conc_NewChan(so_int, (mem_System), (0)), .sum = 0};
    conc_Thread thr = conc_Go(consume, &task, NULL);
    for (so_int i = 0; i < so_len(nums); i++) {
        so_at(so_int, nums, i) = i;
        for (; conc_Chan_SendTimeout(so_int, (&task.ch), (&so_at(so_int, nums, i)), (50 * time_Millisecond)) != conc_Ok;) {
        }
    }
    // No receiver ready yet; keep retrying.
    conc_Chan_Close(so_int, (&task.ch));
    conc_Thread_Wait(thr);
    // Sum of 0..99.
    if (task.sum != 4950) {
        so_panic("wrong timeout send sum");
    }
    conc_Chan_Free(so_int, (&task.ch));
    so_println("%s", "ok");
}

// -- main.go --

int main(void) {
    so_println("%s", "solod.dev/so/conc");
    testThread();
    testPool();
    testChan();
    return 0;
}

// -- pool.go --

static void testPool(void) {
    testPool_ParallelMap();
    testPool_BackPressure();
    testPool_QueueLarge();
    testPool_QueueOne();
    testPool_Error();
}

static void square(void* arg) {
    main_Task* task = (main_Task*)arg;
    task->out = task->in * task->in;
}

// Squares 0..99 in parallel and checks every result.
static void testPool_ParallelMap(void) {
    so_print("%s", "- parallel map...");
    const int64_t n = 100;
    so_Slice tasks = so_make_slice(main_Task, n, n);
    conc_PoolOpts opts = (conc_PoolOpts){.NumThreads = 8};
    conc_Pool* p = conc_NewPool(mem_System, opts);
    for (so_int i = 0; i < so_len(tasks); i++) {
        so_at(main_Task, tasks, i).in = i;
        conc_Pool_Go(p, square, &so_at(main_Task, tasks, i));
    }
    conc_Pool_Wait(p);
    for (so_int i = 0; i < so_len(tasks); i++) {
        if (so_at(main_Task, tasks, i).out != i * i) {
            conc_Pool_Free(p);
            so_panic("wrong square result");
        }
    }
    so_println("%s", "ok");
    conc_Pool_Free(p);
}

// Submits far more tasks than workers, exercising the queue-full wait.
static void testPool_BackPressure(void) {
    so_print("%s", "- back-pressure...");
    const int64_t n = 1000;
    so_Slice tasks = so_make_slice(main_Task, n, n);
    conc_PoolOpts opts = (conc_PoolOpts){.NumThreads = 2};
    conc_Pool* p = conc_NewPool(mem_System, opts);
    for (so_int i = 0; i < so_len(tasks); i++) {
        so_at(main_Task, tasks, i).in = i;
        conc_Pool_Go(p, square, &so_at(main_Task, tasks, i));
    }
    conc_Pool_Wait(p);
    so_int sum = 0;
    for (so_int i = 0; i < so_len(tasks); i++) {
        sum += so_at(main_Task, tasks, i).out;
    }
    // Sum of i*i for i in 0..999.
    if (sum != 332833500) {
        conc_Pool_Free(p);
        so_panic("wrong sum");
    }
    so_println("%s", "ok");
    conc_Pool_Free(p);
}

// Uses a queue far larger than the worker limit, so most submissions
// enqueue without blocking. All results must still be correct.
static void testPool_QueueLarge(void) {
    so_print("%s", "- queue larger than workers...");
    const int64_t n = 200;
    so_Slice tasks = so_make_slice(main_Task, n, n);
    conc_PoolOpts opts = (conc_PoolOpts){.NumThreads = 2, .QueueSize = 128};
    conc_Pool* p = conc_NewPool(mem_System, opts);
    for (so_int i = 0; i < so_len(tasks); i++) {
        so_at(main_Task, tasks, i).in = i;
        conc_Pool_Go(p, square, &so_at(main_Task, tasks, i));
    }
    conc_Pool_Wait(p);
    for (so_int i = 0; i < so_len(tasks); i++) {
        if (so_at(main_Task, tasks, i).out != i * i) {
            conc_Pool_Free(p);
            so_panic("wrong square result");
        }
    }
    so_println("%s", "ok");
    conc_Pool_Free(p);
}

// Uses the smallest possible queue, so each submission past the first must
// wait for a worker to drain a slot. This stresses the queue-full
// back-pressure path with an explicit queue size.
static void testPool_QueueOne(void) {
    so_print("%s", "- queue of size one...");
    const int64_t n = 50;
    so_Slice tasks = so_make_slice(main_Task, n, n);
    conc_PoolOpts opts = (conc_PoolOpts){.NumThreads = 4, .QueueSize = 1};
    conc_Pool* p = conc_NewPool(mem_System, opts);
    for (so_int i = 0; i < so_len(tasks); i++) {
        so_at(main_Task, tasks, i).in = i;
        conc_Pool_Go(p, square, &so_at(main_Task, tasks, i));
    }
    conc_Pool_Wait(p);
    for (so_int i = 0; i < so_len(tasks); i++) {
        if (so_at(main_Task, tasks, i).out != i * i) {
            conc_Pool_Free(p);
            so_panic("wrong square result");
        }
    }
    so_println("%s", "ok");
    conc_Pool_Free(p);
}

static void checkEven(void* arg) {
    main_Task* task = (main_Task*)arg;
    if (task->in % 2 != 0) {
        task->err = errOddInput;
        return;
    }
    task->out = task->in;
}

// Checks that a task can report an error through its argument struct.
static void testPool_Error(void) {
    so_print("%s", "- error field...");
    const int64_t n = 10;
    so_Slice tasks = so_make_slice(main_Task, n, n);
    conc_PoolOpts opts = (conc_PoolOpts){.NumThreads = 4};
    conc_Pool* p = conc_NewPool(mem_System, opts);
    for (so_int i = 0; i < so_len(tasks); i++) {
        so_at(main_Task, tasks, i).in = i;
        conc_Pool_Go(p, checkEven, &so_at(main_Task, tasks, i));
    }
    conc_Pool_Wait(p);
    for (so_int i = 0; i < so_len(tasks); i++) {
        if (i % 2 != 0 && so_at(main_Task, tasks, i).err.self != errOddInput.self) {
            conc_Pool_Free(p);
            so_panic("expected error for odd input");
        }
        if (i % 2 == 0 && so_at(main_Task, tasks, i).err.self != NULL) {
            conc_Pool_Free(p);
            so_panic("unexpected error for even input");
        }
    }
    so_println("%s", "ok");
    conc_Pool_Free(p);
}

// -- thread.go --

static void testThread(void) {
    testThread_Wait();
    testThread_Detach();
}

static void* increment(void* arg) {
    so_int* n = (so_int*)arg;
    *n = *n + 1;
    return arg;
}

// Starts a thread per element, waits for them all, and checks every result.
static void testThread_Wait(void) {
    so_print("%s", "- wait...");
    const int64_t n = 16;
    so_Slice nums = so_make_slice(so_int, n, n);
    so_Slice threads = so_make_slice(conc_Thread, n, n);
    for (so_int i = 0; i < so_len(nums); i++) {
        so_at(so_int, nums, i) = i;
        so_at(conc_Thread, threads, i) = conc_Go(increment, &so_at(so_int, nums, i), NULL);
    }
    for (so_int i = 0; i < so_len(threads); i++) {
        void* res = conc_Thread_Wait(so_at(conc_Thread, threads, i));
        if (*((so_int*)res) != i + 1) {
            so_panic("wrong increment result");
        }
    }
    for (so_int i = 0; i < so_len(nums); i++) {
        if (so_at(so_int, nums, i) != i + 1) {
            so_panic("wrong increment result");
        }
    }
    so_println("%s", "ok");
}

// squareLatch squares l.out in place, then marks the latch done.
static void* squareLatch(void* arg) {
    latch* l = (latch*)arg;
    sync_Mutex_Lock(&l->mu);
    l->out = l->out * l->out;
    l->done = true;
    sync_Cond_Broadcast(&l->cond);
    sync_Mutex_Unlock(&l->mu);
    return NULL;
}

// Runs a task on a detached thread and waits for it through a condition.
static void testThread_Detach(void) {
    so_print("%s", "- detach...");
    latch l = {0};
    sync_Mutex_Init(&l.mu);
    sync_Cond_Init(&l.cond, &l.mu);
    l.out = 9;
    conc_Thread th = conc_Go(squareLatch, &l, NULL);
    conc_Thread_Detach(th);
    sync_Mutex_Lock(&l.mu);
    for (; !l.done;) {
        sync_Cond_Wait(&l.cond);
    }
    sync_Mutex_Unlock(&l.mu);
    if (l.out != 81) {
        so_panic("wrong detached result");
    }
    sync_Mutex_Free(&l.mu);
    sync_Cond_Free(&l.cond);
    so_println("%s", "ok");
}
