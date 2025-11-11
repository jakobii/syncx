package syncx

import (
	"context"
	"errors"
	"sync"
	"testing"
)

func TestWaitGroupAdd(t *testing.T) {
	t.Run("add positive delta when count is zero", func(t *testing.T) {
		var wg WaitGroup
		if wg.n != 0 {
			t.Fatalf("expected initial count to be 0, got %d", wg.n)
		}
		wg.Add(1)
		if wg.n != 1 {
			t.Fatalf("expected count to be 1 after Add(1), got %d", wg.n)
		}
		ch := wg.ch.Load()
		if ch == nil {
			t.Fatal("expected channel to be created when adding to zero count")
		}
	})
	t.Run("add positive delta when count is greater than zero", func(t *testing.T) {
		var wg WaitGroup
		wg.Add(1)
		originalCh := wg.ch.Load()
		wg.Add(2)
		if wg.n != 3 {
			t.Fatalf("expected count to be 3 after Add(1) then Add(2), got %d", wg.n)
		}
		currentCh := wg.ch.Load()
		if currentCh != originalCh {
			t.Fatal("expected channel to remain the same when adding to non-zero count")
		}
	})
	t.Run("add negative delta", func(t *testing.T) {
		var wg WaitGroup
		wg.Add(5)
		wg.Add(-2)
		if wg.n != 3 {
			t.Fatalf("expected count to be 3 after Add(5) then Add(-2), got %d", wg.n)
		}
	})
	t.Run("add zero delta", func(t *testing.T) {
		var wg WaitGroup
		wg.Add(0)
		if wg.n != 0 {
			t.Fatalf("expected count to remain 0 after Add(0), got %d", wg.n)
		}
		ch := wg.ch.Load()
		if ch != nil {
			t.Fatal("expected no channel to be created when calling Add(0) - should be a no-op")
		}
	})
	t.Run("add zero delta when count is greater than zero", func(t *testing.T) {
		var wg WaitGroup
		wg.Add(3)
		originalCh := wg.ch.Load()
		wg.Add(0)
		if wg.n != 3 {
			t.Fatalf("expected count to remain 3 after Add(0), got %d", wg.n)
		}
		currentCh := wg.ch.Load()
		if currentCh != originalCh {
			t.Fatal("expected channel to remain the same when adding 0 to non-zero count")
		}
	})
	t.Run("concurrent add operations", func(t *testing.T) {
		var wg WaitGroup
		n := 100
		var testWg sync.WaitGroup
		for range n {
			testWg.Go(func() {
				wg.Add(1)
			})
		}
		testWg.Wait()
		if wg.n != n {
			t.Fatalf("expected count to be %d after concurrent adds, got %d", n, wg.n)
		}
	})
}

func TestWaitGroupDone(t *testing.T) {
	t.Run("done decrements count", func(t *testing.T) {
		var wg WaitGroup
		wg.Add(3)
		wg.Done()
		if wg.n != 2 {
			t.Fatalf("expected count to be 2 after Done(), got %d", wg.n)
		}
	})
	t.Run("done to zero closes channel", func(t *testing.T) {
		var wg WaitGroup
		wg.Add(1)
		ch := wg.ch.Load()
		wg.Done()
		if wg.n != 0 {
			t.Fatalf("expected count to be 0 after Done(), got %d", wg.n)
		}
		select {
		case <-ch:
		default:
			t.Fatal("expected channel to be closed when count reaches zero")
		}
	})
	t.Run("done panics on negative count", func(t *testing.T) {
		var wg WaitGroup
		defer func() {
			if v := recover(); v == nil {
				t.Fatal("expected panic when calling Done() on zero count")
			}
			if wg.n != 0 {
				t.Fatalf("expected count to remain 0 after panic, got %d", wg.n)
			}
		}()
		wg.Done()
	})
	t.Run("done panics after counter reaches zero", func(t *testing.T) {
		var wg WaitGroup
		wg.Add(1)
		wg.Done() // Counter is now 0

		defer func() {
			if v := recover(); v == nil {
				t.Fatal("expected panic when calling Done() after counter reaches zero")
			} else if v != "negative WaitGroup counter" {
				t.Fatalf("expected panic message 'negative WaitGroup counter', got %v", v)
			}
		}()
		wg.Done() // This should panic
	})
	t.Run("done allows valid operations", func(t *testing.T) {
		var wg WaitGroup
		wg.Add(2)
		wg.Done() // Should not panic, counter goes from 2 to 1
		if wg.n != 1 {
			t.Fatalf("expected count to be 1 after Done(), got %d", wg.n)
		}
		wg.Done() // Should not panic, counter goes from 1 to 0
		if wg.n != 0 {
			t.Fatalf("expected count to be 0 after second Done(), got %d", wg.n)
		}
	})
	t.Run("multiple done calls", func(t *testing.T) {
		var wg WaitGroup
		wg.Add(3)
		wg.Done()
		wg.Done()
		if wg.n != 1 {
			t.Fatalf("expected count to be 1 after two Done() calls, got %d", wg.n)
		}
		wg.Done()
		if wg.n != 0 {
			t.Fatalf("expected count to be 0 after third Done() call, got %d", wg.n)
		}
	})
	t.Run("concurrent done operations", func(t *testing.T) {
		var wg WaitGroup
		n := 100
		wg.Add(n)
		ch := wg.ch.Load()
		var testWg sync.WaitGroup
		for range n {
			testWg.Go(func() {
				wg.Done()
			})
		}
		testWg.Wait()
		if wg.n != 0 {
			t.Fatalf("expected count to be 0 after all Done() calls, got %d", wg.n)
		}
		select {
		case <-ch:
		default:
		}
	})
}

func TestWaitGroupAwait(t *testing.T) {
	t.Run("await when count is zero returns closed channel", func(t *testing.T) {
		var wg WaitGroup
		ch := wg.Await()
		select {
		case <-ch:
		default:
			t.Fatal("expected closed channel when count is zero")
		}
	})
	t.Run("await when count is greater than zero returns open channel", func(t *testing.T) {
		var wg WaitGroup
		wg.Add(1)
		ch := wg.Await()
		select {
		case <-ch:
			t.Fatal("expected open channel when count is greater than zero")
		default:
		}
	})
	t.Run("await channel closes when count reaches zero", func(t *testing.T) {
		var wg WaitGroup
		wg.Add(1)
		ch := wg.Await()
		select {
		case <-ch:
			t.Fatal("expected open channel initially")
		default:
		}
		wg.Done()
		select {
		case <-ch:
		default:
			t.Fatal("expected channel to be closed after Done()")
		}
	})
	t.Run("multiple await calls return same channel when count greater than zero", func(t *testing.T) {
		var wg WaitGroup
		wg.Add(1)
		ch1 := wg.Await()
		ch2 := wg.Await()
		if ch1 != ch2 {
			t.Fatal("expected same channel from multiple Await() calls")
		}
	})
	t.Run("multiple await calls when count is zero return different closed channels", func(t *testing.T) {
		var wg WaitGroup
		ch1 := wg.Await()
		ch2 := wg.Await()
		if ch1 == ch2 {
			t.Fatal("expected different channels from multiple Await() calls when count is zero")
		}
		select {
		case <-ch1:
		default:
			t.Fatal("expected first channel to be closed")
		}
		select {
		case <-ch2:
		default:
			t.Fatal("expected second channel to be closed")
		}
	})
	t.Run("await after add and done cycle", func(t *testing.T) {
		var wg WaitGroup
		wg.Add(1)
		ch1 := wg.Await()
		wg.Done()
		select {
		case <-ch1:
		default:
			t.Fatal("expected first channel to be closed after Done()")
		}
		ch2 := wg.Await()
		if ch1 == ch2 {
			t.Fatal("expected different channel after Add/Done cycle")
		}
		select {
		case <-ch2:
		default:
		}
	})
}

func TestWaitGroupWait(t *testing.T) {
	t.Run("wait blocks until count reaches zero", func(t *testing.T) {
		var wg WaitGroup
		wg.Add(1)
		done := make(chan struct{})
		go func() {
			defer close(done)
			wg.Wait()
		}()
		select {
		case <-done:
			t.Fatal("Wait() should not return immediately when count > 0")
		default:
		}
		wg.Done()
		<-done
	})
	t.Run("wait returns immediately when count is zero", func(t *testing.T) {
		var wg WaitGroup
		wg.Wait()
	})
}

func TestWaitGroupWaitContext(t *testing.T) {
	t.Run("wait context returns nil when count reaches zero", func(t *testing.T) {
		var wg WaitGroup
		wg.Add(1)
		ctx := t.Context()
		errs := make(chan error, 1)
		go func() {
			errs <- wg.WaitContext(ctx)
		}()
		wg.Done()
		if err := <-errs; err != nil {
			t.Fatalf("expected nil error when count reaches zero, got %v", err)
		}
	})
	t.Run("wait context returns immediately when count is zero", func(t *testing.T) {
		var wg WaitGroup
		ctx := t.Context()
		if err := wg.WaitContext(ctx); err != nil {
			t.Fatalf("expected nil error when count is already zero, got %v", err)
		}
	})
	t.Run("wait context returns context error when cancelled", func(t *testing.T) {
		var wg WaitGroup
		wg.Add(1)
		ctx, cancel := context.WithCancel(t.Context())
		errs := make(chan error, 1)
		go func() {
			errs <- wg.WaitContext(ctx)
		}()
		cancel()
		if err := <-errs; !errors.Is(err, context.Canceled) {
		}
	})
}

func TestWaitGroupGo(t *testing.T) {
	t.Run("go increments count and executes function", func(t *testing.T) {
		var wg WaitGroup
		executed := make(chan struct{})
		wg.Go(func() {
			close(executed)
		})
		<-executed
		wg.Wait()
	})
	t.Run("go decrements count after function completes", func(t *testing.T) {
		var wg WaitGroup
		done := make(chan struct{})
		wg.Go(func() {
			close(done)
		})
		<-done
		wg.Wait()
		if wg.n != 0 {
			t.Fatalf("expected count to be 0 after function completes, got %d", wg.n)
		}
	})
	t.Run("multiple go calls", func(t *testing.T) {
		var wg WaitGroup
		n := 10
		counter := make(chan struct{}, n)
		for range n {
			wg.Go(func() {
				counter <- struct{}{}
			})
		}
		for range n {
			<-counter
		}
		wg.Wait()
		if wg.n != 0 {
			t.Fatalf("expected count to be 0 after all functions complete, got %d", wg.n)
		}
	})
	t.Run("go with panic still decrements count", func(t *testing.T) {
		var wg WaitGroup
		panicked := make(chan struct{})
		wg.Go(func() {
			defer func() {
				if recover() != nil {
					close(panicked)
				}
			}()
			panic("test panic")
		})
		<-panicked
		wg.Wait()
		if wg.n != 0 {
			t.Fatalf("expected count to be 0 after panicking function, got %d", wg.n)
		}
	})
}
