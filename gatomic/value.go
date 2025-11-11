package gatomic

import (
	"sync/atomic"
)

// Value wraps [sync/atomic.Value] with generics to make it a bit more
// convenient to use in modern Go. It inits the value to the zero value rather
// than deal with existence of nil.
type Value[T any] struct {
	V atomic.Value
}

// init replaces [sync/atomic.Value]'s initial value of nil with a zero value T.
func (v *Value[T]) init() {
	var x T
	_ = v.V.CompareAndSwap(nil, x)
}

// CompareAndSwap swaps old with new if old matches current value. Initial
// values are zero values of T.
func (v *Value[T]) CompareAndSwap(old T, new T) (swapped bool) {
	v.init()
	return v.V.CompareAndSwap(old, new)
}

// Load returns the current value. Initial values are zero values of T.
func (v *Value[T]) Load() (val T) {
	v.init()
	return v.V.Load().(T)
}

// Store stores a value, replacing any current value.
func (v *Value[T]) Store(val T) {
	v.V.Store(val)
}

// Swap replaces the current value with a new one. The old value is returned.
// Initial values are zero values of T.
func (v *Value[T]) Swap(new T) (old T) {
	v.init()
	return v.V.Swap(new).(T)
}
