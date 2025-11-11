package syncx

import (
	"context"
	"sync"
)

// Mutex is a drop-in replacement for the standard library's [sync.Mutex]. It
// offers the ability to cancel acquiring a lock by making use of Go's select
// statement. The zero value is safe to use. Satisfies [sync.Locker].
type Mutex struct {
	// x must be a buffered channel with a length of 1. It is considered
	// locked when it has a length of 1. It is considered unlocked when it has a
	// length of 0.
	x    chan struct{}
	once sync.Once
}

// Unlocks m. Panics if m is not locked. Calling Unlock on an unlocked mutex
// usually indicates a race condition.
func (m *Mutex) Unlock() {
	select {
	case <-m.state():
	default:
		panic("unlock of unlocked mutex")
	}
}

// Acquire locks m when sending to its returned channel. Do not close it.
//
//	select {
//	case mu.Acquire() <- syncx.Lock:
//	case <-time.After(time.Minute):
//	}
func (m *Mutex) Acquire() chan<- struct{} {
	return m.state()
}

// Lock locks m. If the lock is already in use, the calling goroutine blocks
// until the mutex is available.
func (m *Mutex) Lock() {
	m.Acquire() <- Lock
}

// TryLock tries to lock m and reports whether it succeeded.
//
// Note that while correct uses of TryLock do exist, they are rare, and use of
// TryLock is often a sign of a deeper problem in a particular use of mutexes.
func (m *Mutex) TryLock() bool {
	select {
	case m.Acquire() <- Lock:
		return true
	default:
		return false
	}
}

// WaitLock locks m or returns ctx's error.
func (m *Mutex) WaitLock(ctx context.Context) error {
	select {
	case <-ctx.Done():
		return ctx.Err()
	case m.Acquire() <- Lock:
		return nil
	}
}

// Get the raw chan. Init it if not done so yet.
func (m *Mutex) state() chan struct{} {
	m.once.Do(func() {
		m.x = make(chan struct{}, 1)
	})
	return m.x
}

var Lock struct{}
