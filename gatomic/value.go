package gatomic

import "sync/atomic"

// Value wraps [sync/atomic.Value] with generics to make it a bit more
// convenient to use in modern Go.
type Value[T any] struct {
	V atomic.Value
}

func (v *Value[T]) CompareAndSwap(old T, new T) (swapped bool) {
	return v.V.CompareAndSwap(old, new)
}

// Load returns the zero value of T or the stored value.
func (v *Value[T]) Load() (val T) {
	loaded := v.V.Load()
	if loaded == nil {
		return val
	}
	return loaded.(T)
}

func (v *Value[T]) Store(val T) {
	v.V.Store(val)
}

func (v *Value[T]) Swap(new T) (old T) {
	return v.V.Swap(new).(T)
}
