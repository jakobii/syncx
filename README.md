# syncx

[![Go Reference](https://pkg.go.dev/badge/github.com/jakobii/syncx.svg)](https://pkg.go.dev/github.com/jakobii/syncx)
[![Test](https://github.com/jakobii/syncx/actions/workflows/test.yml/badge.svg?branch=main)](https://github.com/jakobii/syncx/actions/workflows/test.yml)
[![Go Report Card](https://goreportcard.com/badge/github.com/jakobii/syncx)](https://goreportcard.com/report/github.com/jakobii/syncx)

Enhanced synchronization primitives for Go. A collection of context-aware, channel-based sync utilities that extend the standard library's `sync` package.

```
go get github.com/jakobii/syncx
```

## Features

- 📦 **Drop-in replacement** for `sync.Mutex`
- 👀 **Channel-based API** that works with `select` statements
- 🫡 **Zero dependencies** beyond Go standard library

## Usage

The standard library's `sync.Mutex` is excellent for most use cases and provides the foundation for Go's concurrency model. However, there are scenarios where additional capabilities become valuable. When coordinating longer-running operations or building systems that need to respect contexts and timeouts, the ability to cancel lock acquisition can be essential.

Use `syncx.Mutex` as a drop-in replacement for `sync.Mutex` and it satisfies the `sync.Locker` interface:

```go
var mu syncx.Mutex
mu.Lock()
defer mu.Unlock()
```

For context-aware locking:

```go
var mu syncx.Mutex
if err := mu.LockContext(ctx); err != nil {
	return fmt.Errorf("context ended before work started: %w", err)
}
defer mu.Unlock()
fmt.Println("acquired lock")
```

For integration with `select` statements:

```go
var mu syncx.Mutex
select {
case mu.Acquire() <- syncx.Lock:
	defer mu.Unlock()
	fmt.Println("acquired lock")
case <-time.After(someDuration):
	return errTimeout
}
```

See more [examples](./example_test.go).