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

func NewEvent() *Event {
	return &Event{
		done: make(chan struct{}),
	}
}

func (e *Event) Done() <-chan struct{} {
	return e.done
}

func (e *Event) Fire() {
	e.once.Do(func() {
		atomic.StoreInt32(&e.fired, 1)
		close(e.done)
	})
}

func (e *Event) HasFired() bool {
	return atomic.CompareAndSwapInt32(&e.fired, 1, 1)
}
