package syncx

import (
	"sync"
	"sync/atomic"
)

// Event marks event complete in a one-time, inspired by grpc,
// see https://github.com/grpc/grpc-go/blob/v1.36.0/internal/grpcsync/event.go
type Event struct {
	fired int32
	done  chan struct{}
	once  sync.Once
}

// NewEvent returns a ready-to-use instance
func NewEvent() *Event {
	return &Event{
		done: make(chan struct{}),
	}
}

// Done returns a read-only channel
func (e *Event) Done() <-chan struct{} {
	return e.done
}

// Fire fires event done and close the done channel
func (e *Event) Fire() {
	e.once.Do(func() {
		atomic.StoreInt32(&e.fired, 1)
		close(e.done)
	})
}

// HasFired returns a value whether the event has done or not
func (e *Event) HasFired() bool {
	return atomic.LoadInt32(&e.fired) == 1
}
