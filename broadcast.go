package syncx

import "sync"

type Broadcast struct {
	mu sync.Mutex
	ch chan struct{}
}

// Recv will block until signaled by [Broadcast.Send].
func (b *Broadcast) Recv() <-chan struct{} {
	b.mu.Lock()
	defer b.mu.Unlock()
	if b.ch == nil {
		b.ch = make(chan struct{})
	}
	return b.ch
}

// Send will unblock all [Broadcast.Recv], if there is one.
func (b *Broadcast) Send() {
	b.mu.Lock()
	defer b.mu.Unlock()
	if b.ch == nil {
		return
	}
	close(b.ch)
	b.ch = nil
}
