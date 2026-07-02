package sync

// Mutex is a mutual exclusion lock.
//
// The zero value is not usable; call [Mutex.Init] before use.
// A Mutex must not be copied after Init.
type Mutex struct {
	mu pthread_mutex_t
}

// Init prepares m for use, leaving it unlocked. It must be called exactly once
// before any other method. A Mutex must not be copied after Init.
func (m *Mutex) Init() {
	rc := pthread_mutex_init(&m.mu, nil)
	if rc != 0 {
		panic("sync: Mutex.Init failed")
	}
}

// Lock locks m, blocking until the lock is available.
func (m *Mutex) Lock() {
	rc := pthread_mutex_lock(&m.mu)
	if rc != 0 {
		panic("sync: Mutex.Lock failed")
	}
}

// TryLock tries to lock m and reports whether it succeeded.
// It returns false without blocking if the lock is already held.
func (m *Mutex) TryLock() bool {
	rc := pthread_mutex_trylock(&m.mu)
	if rc != 0 && rc != eBUSY {
		panic("sync: Mutex.TryLock failed")
	}
	return rc == 0
}

// Unlock unlocks m.
func (m *Mutex) Unlock() {
	rc := pthread_mutex_unlock(&m.mu)
	if rc != 0 {
		panic("sync: Mutex.Unlock failed")
	}
}

// Free releases the resources held by m. The Mutex is unusable afterward.
func (m *Mutex) Free() {
	rc := pthread_mutex_destroy(&m.mu)
	if rc != 0 {
		panic("sync: Mutex.Free failed")
	}
}
