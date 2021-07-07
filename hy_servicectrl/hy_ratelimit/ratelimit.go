package hy_ratelimit

import (
	"github.com/juju/ratelimit"
	"sync"
)

type RateLimit struct {
	qps int
	// mu guards the fields below it.
	mu        sync.Mutex
	bucketMap map[interface{}]*ratelimit.Bucket
}

func NewRateLimit(qps int) *RateLimit {
	if qps <= 0 {
		panic("invalid qps")
	}
	return &RateLimit{
		qps:       qps,
		bucketMap: map[interface{}]*ratelimit.Bucket{},
	}
}

func (p *RateLimit) Limited(key interface{}) bool {
	if key == nil || key == "" {
		return false
	}
	p.mu.Lock()
	defer p.mu.Unlock()
	bucket, ok := p.bucketMap[key]
	if !ok {
		bucket = ratelimit.NewBucketWithRate(float64(p.qps), int64(p.qps))
		p.bucketMap[key] = bucket
	}
	return bucket.TakeAvailable(1) <= 0
}
