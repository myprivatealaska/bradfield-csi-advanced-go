package main

import (
	"fmt"
	"time"
)

// Problem: Loop variable capture
// Solution: i is shared by all goroutines so we need to supply it to the inner function implicitly!
func main() {
	for i := 0; i < 10; i++ {
		go func(i int) {
			fmt.Printf("launched goroutine %d\n", i)
		}(i)
	}
	// Wait for goroutines to finish
	time.Sleep(time.Second)
}
