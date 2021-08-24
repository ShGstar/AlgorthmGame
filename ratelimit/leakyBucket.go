package ratelimit

import (
	"sync"
	"time"
)

//速率 = capacity / fullTime
type LeakyBucket struct {
	capacity    int
	interval    time.Duration // 漏出水滴的时间
	inDrops     int
	lastRunTime time.Time
	mutex       sync.Mutex
}

func NewLeakyBucket(interval time.Duration, capacity int) *LeakyBucket {
	return &LeakyBucket{
		capacity:    capacity,
		interval:    interval,
		inDrops:     0,
		lastRunTime: time.Time{},
	}
}

func (t *LeakyBucket) Access() bool {
	t.mutex.Lock()
	defer t.mutex.Unlock()

	now := time.Now()
	sub := now.Sub(t.lastRunTime)

	leaks := int(float64(sub) / float64(t.interval))

	if leaks > 0 {
		if t.inDrops <= leaks {
			t.inDrops = 0
		} else {
			t.inDrops = t.inDrops - leaks
		}
		t.lastRunTime = now
	}

	if t.inDrops < t.capacity {
		t.inDrops++
		return true
	}

	return false
}
