package ratelimiter

import (
	"log"
	"sync"
	"time"
)

type RateLimiter struct {
	mu          sync.Mutex
	lastRequest time.Time
}

func (rL *RateLimiter) Allow() bool {
	now := time.Now()
	if now.Sub(rL.lastRequest) < 2*time.Second {
		log.Println("Please re-try after 2 seconds")
		return false
	}

	rL.mu.Lock()
	rL.lastRequest = now
	rL.mu.Unlock()
	return true
}
