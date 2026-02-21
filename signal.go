package syncx

import (
	"context"
	"sync"
	"time"
)

// Signal is a reusable notifier. It wakes one or more goroutines.
//
// A Signal can be used to nudge workers to re-check some shared state:
//
//	var sig syncx.Signal
//	go func() {
//	    for {
//	        select {
//	        case <-ctx.Done():
//	            return
//	        case <-sig.Recv():
//	            // check shared state or retry work
//	        }
//	    }
//	}()
//	sig.Send()
type Signal struct {
	mu     Mutex
	once   sync.Once
	sends  Cond
	sendMu Mutex
	x      chan struct{}
}

func (s *Signal) init() {
	s.once.Do(func() {
		s.sends.L = &s.mu
	})
}

// Recv will block until signaled by [Signal.Send] or [Signal.Cast].
func (s *Signal) Recv() <-chan struct{} {
	s.init()
	s.mu.Lock()
	defer s.mu.Unlock()
	if s.x == nil {
		s.x = make(chan struct{})
		s.sends.Signal()
	}
	return s.x
}

// Cast will unblock all previously waiting [Signal.Recv]'s, if there were any.
func (s *Signal) Cast() {
	_ = s.CastContext(context.Background())
}

// Cast will unblock all previously waiting [Signal.Recv]'s, if there were any.
func (s *Signal) CastContext(ctx context.Context) error {
	if err := s.mu.LockContext(ctx); err != nil {
		return err
	}
	defer s.mu.Unlock()
	if s.x != nil {
		close(s.x)
		s.x = nil
	}
	return nil
}

// Send waits until at least one or more [Signal.Recv] is unblocked.
func (s *Signal) Send() {
	_ = s.SendContext(context.Background())
}

func (s *Signal) SendTimeout(d time.Duration) error {
	ctx, cancel := context.WithTimeout(context.Background(), d)
	defer cancel()
	return s.SendContext(ctx)
}

func (s *Signal) SendContext(ctx context.Context) error {
	s.init()
	if err := s.sendMu.LockContext(ctx); err != nil {
		return err
	}
	defer s.sendMu.Unlock()
	if err := s.mu.LockContext(ctx); err != nil {
		return err
	}
	for s.x == nil {
		s.sends.Signal() // clear previous
		if err := s.sends.WaitContext(ctx); err != nil {
			return err
		}
	}
	defer s.mu.Unlock()
	close(s.x)
	s.x = nil
	return nil
}
