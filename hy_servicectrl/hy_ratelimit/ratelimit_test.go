package hy_ratelimit

import (
	"github.com/stretchr/testify/assert"
	"sync"
	"testing"
	"time"
)

func TestRateLimit(t *testing.T) {

	limit := NewRateLimit(100)

	var wg sync.WaitGroup

	wg.Add(1)
	go func() {
		defer wg.Done()
		for i := 0; i < 100; i++ {
			assert.False(t, limit.Limited(nil))
		}
	}()

	wg.Add(1)
	go func() {
		defer wg.Done()
		for i := 0; i < 100; i++ {
			assert.False(t, limit.Limited(""))
		}
	}()

	wg.Add(1)
	go func() {
		defer wg.Done()
		for i := 0; i < 100; i++ {
			assert.False(t, limit.Limited(123))
		}
		assert.True(t, limit.Limited(123))
	}()

	wg.Add(1)
	go func() {
		defer wg.Done()
		for i := 0; i < 100; i++ {
			assert.False(t, limit.Limited("key"))
		}
		assert.True(t, limit.Limited("key"))
		time.Sleep(time.Second)
		assert.False(t, limit.Limited("key"))
	}()

	wg.Wait()
}
