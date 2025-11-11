package syncx

import (
	"context"
	"errors"
	"sync"
	"testing"
)

func TestLocker(t *testing.T) {
	var l sync.Locker = &Mutex{}
	l.Lock()
	defer l.Unlock()
}

func TestMutexLockContext(t *testing.T) {
	var mu Mutex
	mu.Acquire() <- struct{}{}
	defer mu.Unlock()
	if len(mu.state()) != 1 {
		t.Fatal("failed to set lock state")
	}
}

func TestMutexLock(t *testing.T) {
	var mu Mutex
	mu.Lock()
	defer mu.Unlock()
	if len(mu.state()) != 1 {
		t.Fatal("failed to set lock state")
	}
}

func TestMutexTryLock(t *testing.T) {
	var mu Mutex
	if !mu.TryLock() {
		t.Fatal("failed to obtain lock")
	}
	defer mu.Unlock()
	if len(mu.state()) != 1 {
		t.Fatal("failed to set lock state")
	}
}

func TestMutexTryLock_already_locked(t *testing.T) {
	var mu Mutex
	mu.state() <- struct{}{}
	if mu.TryLock() {
		t.Fatal("obtain lock")
	}
	defer mu.Unlock()
	if len(mu.state()) != 1 {
		t.Fatal("failed to set lock state")
	}
}

func TestMutexLockCtx(t *testing.T) {
	var mu Mutex
	mu.LockContext(t.Context())
	defer mu.Unlock()
	if len(mu.state()) != 1 {
		t.Fatal("failed to set lock state")
	}
}

func TestMutexLockCtx_cancels(t *testing.T) {
	var mu Mutex
	mu.state() <- struct{}{}
	ctx, cancel := context.WithCancel(t.Context())
	go cancel()
	if err := mu.LockContext(ctx); !errors.Is(err, context.Canceled) {
		t.Fatal("did not receive context cancel error")
	}
}

func TestMutexUnlock(t *testing.T) {
	var mu Mutex
	mu.state() <- struct{}{}
	mu.Unlock()
	if len(mu.state()) != 0 {
		t.Fatal("failed to set unlock state")
	}
}

func TestMutexUnlock_panics_when_already_unlocked(t *testing.T) {
	var mu Mutex
	defer func() {
		if v := recover(); v == nil {
			t.Fatal("failed to panic when unlocking an unlocked mutex")
		}
		if len(mu.state()) != 0 {
			t.Fatal("mutated state of unlocked mutex")
		}
	}()
	mu.Unlock()
}

// if Mutex does not work this is a race condition.
// must be tested with "-race"
func TestMutexLock_race(t *testing.T) {
	var mu Mutex
	var i int // some non atomic value to mutate.
	n := 100
	var wg sync.WaitGroup
	wg.Add(n)
	for range n {
		go func() {
			defer wg.Done()
			mu.Lock()
			defer mu.Unlock()
			i++
		}()
	}
	wg.Wait()
	if i != n {
		t.Fatalf("expected %d locks, got %d", n, i)
	}
}
