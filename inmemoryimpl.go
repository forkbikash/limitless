package limitless

import (
	"math"
	"sync"
	"time"
)

// Note:
// This in-memory implementation is not designed for user specific rate limiting
// as huge amount of RAM is required for storing rate limiting info for all users.
// Use case: third party api calls rate limiting
// Refer to redis implementation for flexible rate limiting
// Can't be used in distributed environment

type inMemoryTokenBucket struct {
	mu         sync.Mutex // lock
	capacity   int64      // max tokens
	rate       int        // request per second
	tokens     int64      // available tokens to be used
	lastUpdate time.Time  // tokens updated at
}

func NewInMemoryTokenBucket(capacity int64, rate int) *inMemoryTokenBucket {
	return &inMemoryTokenBucket{
		capacity:   capacity,
		rate:       rate,
		tokens:     capacity,
		lastUpdate: time.Now(),
	}
}

func (tb *inMemoryTokenBucket) allow() (*bool, error) {
	tb.mu.Lock()
	defer tb.mu.Unlock()

	err := tb.load()
	if err != nil {
		return nil, err
	}

	tb.refill()

	allow := tb.available()
	if allow != nil && !*allow {
		return boolPtr(false), nil
	}

	err = tb.exhaust()
	if err != nil {
		return nil, err
	}

	return boolPtr(true), nil
}

func (tb *inMemoryTokenBucket) exhaust() error {
	if tb.tokens > 0 {
		tb.lastUpdate = time.Now()
		tb.tokens--
	}

	return nil
}

func (tb *inMemoryTokenBucket) available() *bool {
	if tb.tokens < 1 {
		return boolPtr(false)
	}
	return nil
}

func (tb *inMemoryTokenBucket) refill() {
	elapsed := time.Since(tb.lastUpdate)
	tokensToAdd := math.Floor(elapsed.Seconds() * float64(tb.rate))
	if tokensToAdd > 0 {
		newTokens := tb.tokens + int64(tokensToAdd)
		tb.lastUpdate = time.Now()
		tb.tokens = int64(math.Min(float64(newTokens), float64(tb.capacity)))
	}
}

func (tb *inMemoryTokenBucket) load() error {
	return nil
}
