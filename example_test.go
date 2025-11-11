package syncx_test

import (
	"context"
	"fmt"
	"time"

	"github.com/jakobii/syncx"
)

// Demonstrates a few ways that Mutex can be used.
func ExampleMutex() {
	var mu syncx.Mutex

	// Lock
	mu.Lock()
	fmt.Println("A drop-in replacement for sync.Mutex.")
	mu.Unlock()

	// TryLock
	if mu.TryLock() {
		fmt.Println("Why not give it a try?")
		mu.Unlock()
	}

	// Acquire
	select {
	case mu.Acquire() <- syncx.Lock:
		fmt.Println("This mutex plays nice with select statements.")
		mu.Unlock()
	default:
		return
	}

	// LockContext
	if err := mu.LockContext(context.Background()); err != nil {
		return
	}
	fmt.Println("Which means it can easily respect contexts.")
	mu.Unlock()

	// Output:
	// A drop-in replacement for sync.Mutex.
	// Why not give it a try?
	// This mutex plays nice with select statements.
	// Which means it can easily respect contexts.
}

func ExampleMutex_LockContext() {
	var mu syncx.Mutex

	// Imagine some other process has the lock.
	mu.Lock()

	// Attempt to acquire the lock with a context. This timeout is short for the
	// sake of the completion of the example.
	ctx, cancel := context.WithTimeout(context.Background(), time.Millisecond)
	defer cancel()

	if err := mu.LockContext(ctx); err != nil {
		fmt.Println("Failed to acquire lock:", err)
		return
	}
	defer mu.Unlock()

	// Output:
	// Failed to acquire lock: context deadline exceeded
}

func ExampleMutex_Acquire() {
	var mu syncx.Mutex

	// The process that sends on the mutex channel is the one that obtains the
	// lock.
	select {
	case <-time.After(time.Millisecond):
		fmt.Println("Failed to acquire lock: timeout")
		return
	case mu.Acquire() <- struct{}{}:
		defer mu.Unlock()
		// ..do things while holding the lock...
		fmt.Println("Successfully acquired lock.")
	}

	// Output:
	// Successfully acquired lock.
}

func ExampleWaitGroup() {
	// Works like sync.WaitGroup.
	var wg syncx.WaitGroup
	for range 3 {
		wg.Go(func() {
			fmt.Println("Worker completed")
		})
	}
	wg.Wait()
	// Output:
	// Worker completed
	// Worker completed
	// Worker completed
}

func ExampleWaitGroup_Await() {
	var wg syncx.WaitGroup

	// Start a worker
	wg.Go(func() {
		time.Sleep(10 * time.Millisecond)
		fmt.Println("Work completed")
	})

	// Use Await in a select statement
	select {
	case <-wg.Await():
		fmt.Println("All work finished")
	case <-time.After(time.Second * 10):
		fmt.Println("Timeout waiting for work")
	}

	// Output:
	// Work completed
	// All work finished
}

func ExampleWaitGroup_WaitContext() {
	var wg syncx.WaitGroup

	// Start a long-running task
	wg.Go(func() {
		time.Sleep(100 * time.Millisecond)
		fmt.Println("Long task completed")
	})

	// Wait with a short timeout
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	if err := wg.WaitContext(ctx); err != nil {
		fmt.Println("Waiting cancelled:", err)
		return
	}
	fmt.Println("All tasks completed")

	// Output:
	// Waiting cancelled: context deadline exceeded
}
