package syncx

import "sync"

type call struct {
	wg  sync.WaitGroup
	v   interface{}
	err error
}

// SingleFlight allows one function to execute in case one more
// goroutines calling in the same time
type SingleFlight struct {
	calls map[string]*call
	mu    sync.Mutex
}

// Do executes and returns the results of the given function, making sure that
// only one function could be executed at a time, the another functions waits
// and shares the same results.
func (s *SingleFlight) Do(key string, fn func() (interface{}, error)) (interface{}, error) {
	s.mu.Lock()
	if s.calls == nil {
		s.calls = make(map[string]*call)
	}

	if c, ok := s.calls[key]; ok {
		s.mu.Unlock()
		c.wg.Wait()
		return c.v, c.err
	}

	c := new(call)
	c.wg.Add(1)
	s.calls[key] = c
	s.mu.Unlock()

	v, err := fn()
	c.v = v
	c.err = err
	s.mu.Lock()
	defer func() {
		c.wg.Done()
		delete(s.calls, key)
		s.mu.Unlock()
	}()

	return v, err
}
