package syncx_test

import (
	"fmt"
	"log"
	"sync"
	"time"

	"github.com/anqiansong/tools/syncx"
)

func ExampleSingleFlight() {
	var s syncx.SingleFlight
	fn := func() (interface{}, error) {
		time.Sleep(200 * time.Millisecond)
		fmt.Println("do")
		return "test", nil
	}

	var wg sync.WaitGroup
	for i := 0; i < 5; i++ {
		wg.Add(1)
		go func() {
			ret, err := s.Do("test", fn)
			if err != nil {
				log.Fatal(err)
			}

			fmt.Println(ret)
			wg.Done()
		}()
	}
	wg.Wait()
	// Output:
	// do
	// test
	// test
	// test
	// test
	// test
}
