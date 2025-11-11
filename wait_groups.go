package syncx

import (
	"context"

	"github.com/jakobii/syncx/gatomic"
)

// WaitGroup is a synchronization primitive that waits for a collection of
// goroutines to finish. It is a drop-in replacement for [sync.WaitGroup] with
// additional methods for enhanced functionality.
//
// The zero value of WaitGroup is ready to use without initialization. A
// WaitGroup must not be copied after first use.
//
// Usage example:
//
//	var wg WaitGroup
//	wg.Go(func() {
//	    // do work
//	})
//
//	// Use Await in select statements
//	select {
//	case <-wg.Await():
//	    fmt.Println("all work completed")
//	case <-time.After(5 * time.Second):
//	    fmt.Println("timeout")
//	}
type WaitGroup struct {
	mu Mutex
	ch gatomic.Value[chan struct{}]
	n  int
}

// Add adds delta to the WaitGroup counter.
func (wg *WaitGroup) Add(delta int) {
	wg.mu.Lock()
	defer wg.mu.Unlock()
	if delta == 0 {
		return
	}
	if wg.n == 0 {
		wg.ch.Store(make(chan struct{}))
	}
	wg.n += delta
}

// Done decrements the WaitGroup counter by one.
func (wg *WaitGroup) Done() {
	wg.mu.Lock()
	defer wg.mu.Unlock()
	if wg.n <= 0 {
		panic("negative WaitGroup counter")
	}
	wg.n--
	if wg.n == 0 {
		ch := wg.ch.Load()
		close(ch)
	}
}

// Go runs f in a new goroutine and adds it to the WaitGroup.
func (wg *WaitGroup) Go(f func()) {
	wg.Add(1)
	go func() {
		defer wg.Done()
		f()
	}()
}

// Wait blocks until the WaitGroup counter reaches zero.
func (wg *WaitGroup) Wait() {
	<-wg.Await()
}

// Await returns a channel that will be closed when the WaitGroup counter
// reaches zero. This allows the wait operation to be used in select statements
// for non-blocking waits or in combination with other channels.
//
// The returned channel should not be closed by the caller. The channel is
// managed internally and will be closed automatically when all goroutines have
// finished.
//
// If the counter is already zero when Await is called, it returns a closed
// channel that will immediately unblock any receive operation.
func (wg *WaitGroup) Await() <-chan struct{} {
	wg.mu.Lock()
	defer wg.mu.Unlock()
	if wg.n == 0 {
		ch := make(chan struct{})
		close(ch)
		return ch
	}
	return wg.ch.Load()
}

// WaitContext waits for the WaitGroup counter to reach zero or for the context
// to be cancelled, whichever happens first. It returns nil if the counter
// reaches zero, or the context's error if the context is cancelled.
func (wg *WaitGroup) WaitContext(ctx context.Context) error {
	select {
	case <-ctx.Done():
		return ctx.Err()
	case <-wg.Await():
		return nil
	}
}
