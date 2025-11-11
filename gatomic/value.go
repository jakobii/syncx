package gatomic

import (
	"sync/atomic"
)

// Value wraps [sync/atomic.Value] with generics to make it a bit more
// convenient to use in modern Go. It inits the value to the Zero value rather
// then deal with existance of nil.
type Value[T any] struct {
	V atomic.Value
}

func (v *Value[T]) init() {
	var x T
	_ = v.V.CompareAndSwap(nil, x)
}

func (v *Value[T]) CompareAndSwap(old T, new T) (swapped bool) {
	v.init()
	return v.V.CompareAndSwap(old, new)
}

func (v *Value[T]) Load() (val T) {
	v.init()
	return v.V.Load().(T)
}

func (v *Value[T]) Store(val T) {
	v.init()
	v.V.Store(val)
}

func (v *Value[T]) Swap(new T) (old T) {
	v.init()
	return v.V.Swap(new).(T)
}
