package ratelimit

import (
	"sync"
	"time"
)

//速率 = capacity / fullTime
type TokenBucket struct {
	capacity    int
	remain      int
	fullTime    time.Duration
	lastRunTime time.Time
	mutex       sync.Mutex
}

func NewTokenBucket(fullTime time.Duration, capacity int) *TokenBucket {
	return &TokenBucket{
		capacity:    capacity,
		remain:      capacity,
		fullTime:    fullTime,
		lastRunTime: time.Now(),
	}
}

func (t *TokenBucket) Access() bool {
	t.mutex.Lock()
	defer t.mutex.Unlock()

	now := time.Now()
	sub := now.Sub(t.lastRunTime)

	//update
	if sub >= t.fullTime { //full
		t.remain = t.capacity
	} else {
		add := int(float64(sub) / float64(t.remain) * float64(t.capacity))

		if add+t.remain >= t.capacity {
			t.remain = t.capacity
		} else {
			t.remain = add + t.remain
		}
	}

	//不够了
	if t.remain <= 0 {
		return false
	}

	t.remain--
	return true
}
