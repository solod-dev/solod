#include "main.h"

// -- Types --

typedef struct gate gate;
typedef struct counter counter;
typedef struct onceTask onceTask;

// gate coordinates a single waiter with the main thread through a condition
// variable and a shared ready flag.
typedef struct gate {
    sync_Mutex* mu;
    sync_Cond* cond;
    bool* ready;
    bool* woke;
} gate;

// counter is a shared count guarded by a mutex.
typedef struct counter {
    sync_Mutex* mu;
    so_int* val;
} counter;

// onceTask carries the shared Once and a slot for the value
// each worker observes right after its Do returns.
typedef struct onceTask {
    sync_Once* once;
    so_int* seen;
} onceTask;

// -- Forward declarations --
static void* waiter(void* arg);
static void testCond(void);
static void testMutex(void);
static void bump(void* arg);
static void testMutex_LockUnlock(void);
static void testMutex_TryLock(void);
static void onceInit(void);
static void callOnce(void* arg);
static void testOnce(void);

// -- Variables and constants --

// onceVal is set by onceInit; onceRuns counts how many times onceInit ran.
static so_int onceVal = 0;
static so_int onceRuns = 0;

// -- cond.go --

static void* waiter(void* arg) {
    gate* g = (gate*)arg;
    sync_Mutex_Lock(g->mu);
    for (; !*g->ready;) {
        sync_Cond_Wait(g->cond);
    }
    *g->woke = true;
    sync_Mutex_Unlock(g->mu);
    return NULL;
}

// Starts a worker that waits on a condition variable until main sets
// a ready flag and broadcasts, then checks the worker observed the signal.
static void testCond(void) {
    so_print("%s", "- cond...");
    sync_Mutex mu = {0};
    sync_Mutex_Init(&mu);
    sync_Cond cond = {0};
    sync_Cond_Init(&cond, &mu);
    bool ready = false;
    bool woke = false;
    gate g = (gate){.mu = &mu, .cond = &cond, .ready = &ready, .woke = &woke};
    conc_Thread thr = conc_Go(waiter, &g, NULL);
    sync_Mutex_Lock(&mu);
    ready = true;
    sync_Cond_Broadcast(&cond);
    sync_Mutex_Unlock(&mu);
    conc_Thread_Wait(thr);
    if (!woke) {
        so_panic("waiter did not observe signal");
    }
    sync_Cond_Free(&cond);
    sync_Mutex_Free(&mu);
    so_println("%s", "ok");
}

// -- main.go --

int main(void) {
    so_println("%s", "solod.dev/so/sync");
    testMutex();
    testCond();
    testOnce();
    return 0;
}

// -- mutex.go --

static void testMutex(void) {
    testMutex_LockUnlock();
    testMutex_TryLock();
}

static void bump(void* arg) {
    counter* c = (counter*)arg;
    sync_Mutex_Lock(c->mu);
    *c->val = *c->val + 1;
    sync_Mutex_Unlock(c->mu);
}

// Checks that no updates are lost when many workers
// concurrently increment a shared counter under a mutex.
static void testMutex_LockUnlock(void) {
    so_print("%s", "- mutex...");
    const int64_t n = 1000;
    sync_Mutex mu = {0};
    sync_Mutex_Init(&mu);
    so_int val = 0;
    so_Slice jobs = so_make_slice(counter, n, n);
    conc_PoolOpts opts = (conc_PoolOpts){.NumThreads = 8};
    conc_Pool* p = conc_NewPool(mem_System, opts);
    for (so_int i = 0; i < so_len(jobs); i++) {
        so_at(counter, jobs, i).mu = &mu;
        so_at(counter, jobs, i).val = &val;
        conc_Pool_Go(p, bump, &so_at(counter, jobs, i));
    }
    conc_Pool_Free(p);
    if (val != n) {
        so_panic("lost updates under mutex");
    }
    sync_Mutex_Free(&mu);
    so_println("%s", "ok");
}

// Checks that TryLock acquires a free mutex and refuses
// to acquire one that is already held.
static void testMutex_TryLock(void) {
    so_print("%s", "- trylock...");
    sync_Mutex mu = {0};
    sync_Mutex_Init(&mu);
    if (!sync_Mutex_TryLock(&mu)) {
        so_panic("TryLock failed on free mutex");
    }
    if (sync_Mutex_TryLock(&mu)) {
        so_panic("TryLock succeeded on held mutex");
    }
    sync_Mutex_Unlock(&mu);
    if (!sync_Mutex_TryLock(&mu)) {
        so_panic("TryLock failed after unlock");
    }
    sync_Mutex_Unlock(&mu);
    sync_Mutex_Free(&mu);
    so_println("%s", "ok");
}

// -- once.go --

// onceInit is the one-time initialization run through sync.Once.
static void onceInit(void) {
    onceVal = 42;
    onceRuns++;
}

static void callOnce(void* arg) {
    onceTask* task = (onceTask*)arg;
    sync_Once_Do(task->once, onceInit);
    *task->seen = onceVal;
}

// Has many workers race on a single Once and checks that the
// initializer ran exactly once and that every Do returned only after
// it completed (each worker observes the initialized value).
static void testOnce(void) {
    so_print("%s", "- once...");
    const int64_t n = 1000;
    sync_Once once = {0};
    sync_Once_Init(&once);
    onceVal = 0;
    onceRuns = 0;
    so_Slice tasks = so_make_slice(onceTask, n, n);
    so_Slice seen = so_make_slice(so_int, n, n);
    conc_PoolOpts opts = (conc_PoolOpts){.NumThreads = 8};
    conc_Pool* p = conc_NewPool(mem_System, opts);
    for (so_int i = 0; i < so_len(tasks); i++) {
        so_at(onceTask, tasks, i).once = &once;
        so_at(onceTask, tasks, i).seen = &so_at(so_int, seen, i);
        conc_Pool_Go(p, callOnce, &so_at(onceTask, tasks, i));
    }
    conc_Pool_Free(p);
    if (onceRuns != 1) {
        so_panic("once ran the initializer more than once");
    }
    for (so_int i = 0; i < so_len(seen); i++) {
        if (so_at(so_int, seen, i) != 42) {
            so_panic("Do returned before the initializer completed");
        }
    }
    sync_Once_Free(&once);
    so_println("%s", "ok");
}
