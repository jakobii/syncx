# syncx

[![Go Reference](https://pkg.go.dev/badge/github.com/jakobii/syncx.svg)](https://pkg.go.dev/github.com/jakobii/syncx)
[![Test](https://github.com/jakobii/syncx/actions/workflows/test.yml/badge.svg?branch=main)](https://github.com/jakobii/syncx/actions/workflows/test.yml)
[![Go Report Card](https://goreportcard.com/badge/github.com/jakobii/syncx)](https://goreportcard.com/report/github.com/jakobii/syncx)

Enhanced synchronization primitives for Go. A collection of context-aware, channel-based sync utilities that extend the standard library's `sync` package.

```
go get github.com/jakobii/syncx
```

## Features

- 🔄 **Drop-in replacement** for `sync.Mutex`
- ⏰ **Context-aware locking** with `WaitLock(ctx)`
- 🎯 **Channel-based API** that works with `select` statements
- 👀 **Zero dependencies** beyond Go standard library

## Usage

The `syncx.Mutex` can work just like `sync.Mutex`.

```go
var mu syncx.Mutex
mu.Lock()
defer mu.Unlock()
```

The standard library's `sync.Mutex` does not offer a way to cancel its `Lock()`
method while it is blocking to acquire the lock. This is usually fine when the
mutex is guarding a resource that does not take much time to access, like a
struct field. But when synchronizing longer-running processes, the need to
cancel work is frequently an issue.

If all you need is to cancel acquiring a lock with a context, this module's
`syncx.Mutex` has a convenient method for this.

```go
var mu syncx.Mutex
func Work(ctx context.Context) error {
	if err := mu.WaitLock(ctx); err != nil {
		return fmt.Errorf("context ended before work started: %w", err)
	}
	defer mu.Unlock()
	fmt.Println("acquired lock, now we can work")
}
```

The `syncx.Mutex` can also work with Go's `select` statement. Sending
`struct{}{}` is a common way of signaling, and here we use it with `Acquire()`
to acquire the lock.

```go
var mu syncx.Mutex
func MyWork(ctx context.Context) error {
	select {
	case <-time.After(someMaxDuration):
		return errTimeout
	case mu.Acquire() <- struct{}{}:
		defer mu.Unlock()
		fmt.Println("acquired lock")
	}
}
```

See more [examples](./example_test.go).