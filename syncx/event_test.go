package syncx

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestEvent(t *testing.T) {
	t.Run("HasFired", func(t *testing.T) {
		e := NewEvent()
		assert.False(t, e.HasFired())
		e.Fire()
		assert.True(t, e.HasFired())
	})

	t.Run("Fire", func(t *testing.T) {
		e := NewEvent()
		for i := 0; i < 3; i++ {
			e.Fire()
		}
		assert.True(t, e.HasFired())
	})

	t.Run("Done", func(t *testing.T) {
		e := NewEvent()
		for i := 0; i < 100; i++ {
			go func(ev *Event) {
				ev.Fire()
			}(e)
		}
		select {
		case <-e.Done():
			assert.True(t, e.HasFired())
		}
	})
}
