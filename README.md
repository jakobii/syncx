# syncx

[![Go Reference](https://pkg.go.dev/badge/github.com/jakobii/syncx.svg)](https://pkg.go.dev/github.com/jakobii/syncx)
[![Test](https://github.com/jakobii/syncx/actions/workflows/test.yml/badge.svg?branch=main)](https://github.com/jakobii/syncx/actions/workflows/test.yml)
[![Go Report Card](https://goreportcard.com/badge/github.com/jakobii/syncx)](https://goreportcard.com/report/github.com/jakobii/syncx)

Enhanced synchronization primitives for Go. A collection of context-aware, channel-based sync utilities that extend the standard library's `sync` package.

```
go get github.com/jakobii/syncx
```

## Features

- 📦 **Drop-in replacement** for `sync.Mutex`, `sync.WaitGroup`
- 👀 **Channel-based API** that works with `select` statements
- 🫡 **Zero dependencies** beyond Go standard library

## Usage

The standard library's `sync` package is excellent for most use cases and provides the foundation for Go's concurrency model. However, there are scenarios where additional capabilities become valuable. When coordinating longer-running operations or building systems that need to respect contexts and timeouts, the ability to cancel blocking operations can be essential.

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

```go
var wg syncx.WaitGroup
wg.Go(func(){
	fmt.Printf("doing some work")
})
select {
case <-wg.Await():
	fmt.Println("work completed")
case <-time.After(someDuration):
	return errTimeout
}
```

See more [examples](./example_test.go).