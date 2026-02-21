package syncx

import (
	"context"
	"sync"
)

type ContextLocker interface {
	sync.Locker
	LockContext(ctx context.Context) error
}

// Cond is similar to sync.Cond and it respects contexts.  it’s safe to use the
// zero value.
type Cond struct {
	once sync.Once
	L    ContextLocker
	C    chan chan struct{}
}

func NewCond(l ContextLocker) *Cond {
	return &Cond{
		L: l,
		C: make(chan chan struct{}, 1),
	}
}

func (c *Cond) init() {
	c.once.Do(func() {
		if c.C == nil {
			c.C = make(chan chan struct{}, 1)
		}
		if c.L == nil {
			c.L = &Mutex{}
		}
	})
}

func (c *Cond) Signal() {
	c.init()
	select {
	case wake := <-c.C:
		close(wake)
	default:
	}
}

func (c *Cond) Broadcast() {
	c.init()
	for {
		select {
		case wake := <-c.C:
			close(wake)
		default:
			return
		}
	}
}

// Wait is short for [Cond.WaitContext].
func (c *Cond) Wait() {
	_ = c.WaitContext(context.Background())
}

// WaitContext waits for a signal and unblock so that caller can check a
// condition.  Callers must hold c.L when calling; on return, L is again held
// (or ctx.Err() is returned).
func (c *Cond) WaitContext(ctx context.Context) error {
	c.init()
	wake := make(chan struct{})
	select {
	case <-ctx.Done():
		return ctx.Err()
	case c.C <- wake:
	}
	c.L.Unlock()
	select {
	case <-ctx.Done():
		return ctx.Err()
	case <-wake:
	}
	if err := c.L.LockContext(ctx); err != nil {
		return err
	}
	return nil
}
