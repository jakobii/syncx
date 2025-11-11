package syncx

import (
	"fmt"
	"testing"
)

func TestDebugAdd(t *testing.T) {
	var wg WaitGroup
	
	fmt.Printf("Initial state: n=%d, ch=%v\n", wg.n, wg.ch.Load())
	
	wg.Add(0)
	
	fmt.Printf("After Add(0): n=%d, ch=%v\n", wg.n, wg.ch.Load())
	
	ch := wg.ch.Load()
	if ch == nil {
		t.Fatal("Channel is nil!")
	}
}