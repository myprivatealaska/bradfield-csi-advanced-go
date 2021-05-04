package main

import (
	"fmt"
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

func (e *CounterWrong) getNext() uint64 {
	e.counter += 1
	return e.counter
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
	//e := CounterWrong{counter: 0}
	////go e.getNext()
	//e.getNext()
	//e.getNext()
	//e.getNext()
	//fmt.Println(e.counter)

	counterM := NewCounterMonitor()
	go counterM.getNext()
	counterM.getNext()
	counterM.getNext()
	fmt.Println(counterM.getNext())
}
