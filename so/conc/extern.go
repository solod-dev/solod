package conc

import "solod.dev/so/c"

//so:include <pthread.h>

//so:embed conc.h
var conc_h string

// pthread_t is an opaque pthread thread handle.
//
//so:extern
type pthread_t struct{}

// pthread_attr_t is an opaque pthread thread-creation attribute set.
//
//so:extern
type pthread_attr_t struct{}

// pthread_create starts a new thread running start(arg).
// The thread handle is stored in *t. attr may be nil for default attributes.
// Returns 0 on success.
//
//so:extern
func pthread_create(t *pthread_t, attr *pthread_attr_t,
	start func(any) any, arg any) c.Int {
	_, _, _, _ = t, attr, start, arg
	return -1
}

// pthread_join blocks until the thread t terminates.
// retval may be nil. Returns 0 on success.
//
//so:extern
func pthread_join(t pthread_t, retval *any) c.Int {
	_, _ = t, retval
	return -1
}

// pthread_detach marks thread t so its resources are released automatically
// when it terminates. A detached thread cannot be joined.
// Returns 0 on success.
//
//so:extern
func pthread_detach(t pthread_t) c.Int {
	_ = t
	return -1
}

// pthread_attr_init initializes a thread attribute set.
// Returns 0 on success.
//
//so:extern
func pthread_attr_init(attr *pthread_attr_t) c.Int {
	_ = attr
	return -1
}

// pthread_attr_destroy destroys a thread attribute set.
// Returns 0 on success.
//
//so:extern
func pthread_attr_destroy(attr *pthread_attr_t) c.Int {
	_ = attr
	return -1
}

// pthread_attr_setstacksize sets the stack size for threads created with a.
// Returns 0 on success.
//
//so:extern
func pthread_attr_setstacksize(attr *pthread_attr_t, size uintptr) c.Int {
	_, _ = attr, size
	return -1
}
