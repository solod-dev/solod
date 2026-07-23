package sync

import "solod.dev/so/c"

//so:include <pthread.h>
//so:include.c <errno.h>
//so:link pthread

//so:embed sync.h
var sync_h string

// eBUSY is returned by pthread_mutex_trylock when the mutex is already locked.
//
//so:extern EBUSY
const eBUSY = 0

// eTIMEDOUT is returned by a timed cond wait when the deadline passes.
//
//so:extern ETIMEDOUT
const eTIMEDOUT = 0

// pthread_mutex_t is an opaque pthread mutex.
//
//so:extern
type pthread_mutex_t struct{}

// pthread_cond_t is an opaque pthread condition variable.
//
//so:extern
type pthread_cond_t struct{}

// pthread_mutex_init initializes a mutex with default attributes.
// Returns 0 on success.
//
//so:extern
func pthread_mutex_init(mu *pthread_mutex_t, attr any) c.Int {
	_, _ = mu, attr
	return 0
}

// pthread_mutex_destroy destroys a mutex.
// Returns 0 on success.
//
//so:extern
func pthread_mutex_destroy(mu *pthread_mutex_t) c.Int {
	_ = mu
	return 0
}

// pthread_mutex_lock locks a mutex, blocking until it is available.
// Returns 0 on success.
//
//so:extern
func pthread_mutex_lock(mu *pthread_mutex_t) c.Int {
	_ = mu
	return 0
}

// pthread_mutex_unlock unlocks a mutex.
// Returns 0 on success.
//
//so:extern
func pthread_mutex_unlock(mu *pthread_mutex_t) c.Int {
	_ = mu
	return 0
}

// pthread_mutex_trylock tries to lock a mutex without blocking.
// Returns 0 if the lock was acquired, EBUSY if it is already locked.
//
//so:extern
func pthread_mutex_trylock(mu *pthread_mutex_t) c.Int {
	_ = mu
	return 0
}

// pthread_cond_destroy destroys a condition variable.
// Returns 0 on success.
//
//so:extern
func pthread_cond_destroy(cond *pthread_cond_t) c.Int {
	_ = cond
	return 0
}

// pthread_cond_wait atomically unlocks m and blocks until cond is signaled,
// then re-locks m before returning. Returns 0 on success.
//
//so:extern
func pthread_cond_wait(cond *pthread_cond_t, mu *pthread_mutex_t) c.Int {
	_, _ = cond, mu
	return 0
}

// pthread_cond_signal wakes at least one thread blocked on cond.
// Returns 0 on success.
//
//so:extern
func pthread_cond_signal(cond *pthread_cond_t) c.Int {
	_ = cond
	return 0
}

// pthread_cond_broadcast wakes all threads blocked on cond.
// Returns 0 on success.
//
//so:extern
func pthread_cond_broadcast(cond *pthread_cond_t) c.Int {
	_ = cond
	return 0
}

// condInitMono initializes cond so its timed waits use the monotonic clock,
// immune to wall-clock changes. Returns 0 on success.
//
//so:extern sync_condInitMono
func condInitMono(cond *pthread_cond_t) c.Int {
	_ = cond
	return 0
}

// condWaitRel atomically unlocks mu and blocks on cond until it is signaled or
// nsec nanoseconds elapse on the monotonic clock, then re-locks mu. A
// non-positive nsec polls without blocking. Returns 0 if signaled,
// ETIMEDOUT on timeout.
//
//so:extern sync_condWaitRel
func condWaitRel(cond *pthread_cond_t, mu *pthread_mutex_t, nsec int64) c.Int {
	_, _, _ = cond, mu, nsec
	return 0
}
