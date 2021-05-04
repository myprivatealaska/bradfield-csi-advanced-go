package main

import (
	"fmt"
	"sync"
	"sync/atomic"
)

type counterService interface {
	// Returns values in ascending order; it should be safe to call
	// getNext() concurrently without any additional synchronization.
	getNext() uint64
}

// WRONG - 1. Don’t perform any synchronization
type CounterWrong struct {
	counter uint64
}

func (c *CounterWrong) getNext() uint64 {
	c.counter += 1
	return c.counter
}

// 2. Atomically increment a counter value using sync/atomic
type CounterAtomic struct {
	counter uint64
}

func (c *CounterAtomic) getNext() uint64 {
	return atomic.AddUint64(&c.counter, 1)
}

// 3. Use a sync.Mutex to guard access to a shared counter value
type CounterMu struct {
	counter uint64
	mu      sync.Mutex
}

func (c *CounterMu) getNext() uint64 {
	c.mu.Lock()
	c.counter += 1
	newC := c.counter
	c.mu.Unlock()
	return newC
}

// 4. Launch a separate goroutine with exclusive access to a private counter value;
//    handle getNext() calls by making “requests” and receiving “responses” on two separate channels
type CounterMonitor struct {
	req  chan struct{}
	resp chan uint64
}

func monitor(req chan struct{}, resp chan uint64) {
	var counter uint64
	for {
		<-req
		counter += 1
		resp <- counter
	}
}

func NewCounterMonitor() *CounterMonitor {
	s := CounterMonitor{
		req:  make(chan struct{}),
		resp: make(chan uint64),
	}
	go monitor(s.req, s.resp)
	return &s
}

func (e *CounterMonitor) getNext() uint64 {
	e.req <- struct{}{} // put a 'request' on the channel
	return <-e.resp     // wait for the 'response'
}

func main() {
	// 1 - Data race detected
	//e := CounterWrong{counter: 0}
	//for i := 0; i < 10; i++ {
	//	go e.getNext()
	//}
	//fmt.Println(e.counter)

	// 2 - use sync.Atomic
	//counterA := CounterAtomic{}
	//for i := 0; i < 10; i++ {
	//	go counterA.getNext()
	//}
	//fmt.Println(counterA.getNext())

	// 3 - use sync.Mutex
	//counterMu := CounterMu{}
	//for i := 0; i < 10; i++ {
	//	go counterMu.getNext()
	//}
	//fmt.Println(counterMu.getNext())

	// 4 - use a separate goroutine with exclusive access to a private counter value
	counterM := NewCounterMonitor()
	for i := 0; i < 10; i++ {
		go counterM.getNext()
	}
	fmt.Println(counterM.getNext())
}
