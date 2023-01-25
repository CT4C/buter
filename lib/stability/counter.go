package stability

import "sync"

type Counter struct {
	counter int
	mu      sync.Mutex
}

func (c *Counter) Increment() {
	c.mu.Lock()
	c.counter++
	c.mu.Unlock()
}

func (c *Counter) Decrement() {
	c.mu.Lock()
	c.counter--
	c.mu.Unlock()
}

func (c *Counter) Number() int {
	defer c.mu.Unlock()
	c.mu.Lock()
	return c.counter
}

func NewCounter() *Counter {
	return &Counter{
		counter: 0,
		mu:      sync.Mutex{},
	}
}
