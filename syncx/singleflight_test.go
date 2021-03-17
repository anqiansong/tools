package syncx

import (
	"fmt"
	"sync"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

func TestSingleFlight_Do(t *testing.T) {
	var s SingleFlight
	fn := func() (interface{}, error) {
		time.Sleep(time.Second)
		fmt.Println("---")
		return 1, nil
	}
	var wg sync.WaitGroup
	for i := 0; i < 10; i++ {
		wg.Add(1)
		go func() {
			ret, err := s.Do("key", fn)
			assert.Nil(t, err)
			fmt.Println(ret)
			assert.Equal(t, 1, ret)
			wg.Done()
		}()
	}

	wg.Wait()
}
