package main

import (
	"fmt"
	"sync"
)

// Problem: deadlock because we're trying to acquire a RLock while still having a Lock in place
// Solution: Remove `defer`, Unlock before calling c.logState()
type coordinator struct {
	lock   sync.RWMutex
	leader string
}

func newCoordinator(leader string) *coordinator {
	return &coordinator{
		lock:   sync.RWMutex{},
		leader: leader,
	}
}

func (c *coordinator) logState() {
	c.lock.RLock()
	defer c.lock.RUnlock()

	fmt.Printf("leader = %q\n", c.leader)
}

func (c *coordinator) setLeader(leader string, shouldLog bool) {
	c.lock.Lock()
	c.leader = leader

	c.lock.Unlock()

	if shouldLog {
		c.logState()
	}
}

func main() {
	c := newCoordinator("us-east")
	c.logState()
	c.setLeader("us-west", true)
}
