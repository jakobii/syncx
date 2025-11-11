package gatomic_test

import (
	"testing"

	"github.com/jakobii/syncx/gatomic"
)

func TestValueCompareAndSwap(t *testing.T) {
	var v gatomic.Value[int]
	if got := v.CompareAndSwap(0, 1); !got {
		t.Fatalf("failed to init zero value")
	}
	if got := v.CompareAndSwap(2, 3); got {
		t.Fatalf("old value should not have matched atomic state")
	}
	if got := v.CompareAndSwap(1, 2); !got {
		t.Fatalf("old value should have matched atomix state")
	}
}

func TestValueLoad(t *testing.T) {
	var v gatomic.Value[int]
	if v.Load() != 0 {
		t.Fatalf("failed to load zero value")
	}
	v.V.Store(123)
	if v.Load() != 123 {
		t.Fatalf("failed to load set int")
	}
	v.Store(456)
	if v.Load() != 456 {
		t.Fatalf("failed to load store int")
	}
}

func TestValueStore(t *testing.T) {
	var v gatomic.Value[int]
	v.Store(456)
	if v.Load() != 456 {
		t.Fatalf("failed to load store int")
	}
}

func TestValueSwap(t *testing.T) {
	var v gatomic.Value[int]
	if v.Swap(123) != 0 {
		t.Fatalf("failed to load zero value")
	}
	if v.Load() != 123 {
		t.Fatalf("failed to load store int")
	}
	if v.Swap(456) != 123 {
		t.Fatalf("failed to swap previous value")
	}
	if v.Load() != 456 {
		t.Fatalf("failed to load store int")
	}
}
