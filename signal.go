package syncx

import "sync"

// Signal is a reusable notifier. It wake one or more goroutines.
//
// A Signal can be used to nudge workers to re-check some shared state:
//
//	var sig syncx.Signal
//	go func() {
//		for {
//			select {
//			case <-ctx.Done():
//				return
//			case <-sig.Recv():
//				// check shared state or retry work
//			}
//		}
//	}()
//	sig.Send()
type Signal struct {
	mu    sync.Mutex
	once  sync.Once
	sends chan chan struct{}
	x     chan struct{}
}

func (s *Signal) init() {
	s.once.Do(func() {
		s.sends = make(chan chan struct{}, 1)
	})
}

// Recv will block until signaled by [Signal.Send].
func (s *Signal) Recv() <-chan struct{} {
	s.init()
	s.mu.Lock()
	defer s.mu.Unlock()
	if s.x == nil {
		s.x = make(chan struct{})
		select {
		case wake := <-s.sends:
			close(wake)
		default:
		}
	}
	return s.x
}

// Send blocks until one or more [Signal.Recv] is unblocked.
func (s *Signal) Send() {
	s.init()
	s.mu.Lock()
	for s.x == nil {
		wake := make(chan struct{})
		s.mu.Unlock()
		s.sends <- wake
		<-wake
		s.mu.Lock()
	}
	defer s.mu.Unlock()
	close(s.x)
	s.x = nil
}

// Cast will unblock all [Signal.Recv], if there is one.
func (s *Signal) Cast() {
	s.mu.Lock()
	defer s.mu.Unlock()
	if s.x == nil {
		return
	}
	close(s.x)
	s.x = nil
}
