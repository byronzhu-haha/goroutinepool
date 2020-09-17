package goroutinepool

import (
	"fmt"
	"testing"
)

func TestPool(t *testing.T) {
	pool := NewPool(3)
	go func() {
		for i := 0; i < 10; i++ {
			pool.Put(NewJob(func(arg interface{}) error {
				fmt.Printf("%+v\n", arg)
				return nil
			}, i))
		}
		pool.PutDone()
	}()
	pool.Run()
}
