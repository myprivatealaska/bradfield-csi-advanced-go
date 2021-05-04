package main

import (
	"fmt"
)

// Problem: The channel was not used correctly
// Solution: put a value on the channel inside of the goroutine, read off of it outside (line 16)
func main() {
	done := make(chan struct{}, 1)
	go func() {
		fmt.Println("performing initialization...")
		done <- struct{}{}
	}()

	<-done
	fmt.Println("initialization done, continuing with rest of program")
}
