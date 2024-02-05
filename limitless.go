package limitless

// implements token bucket algorithm

type RateLimiter interface {
	allow() (*bool, error)
	load() error
	refill()
	available() *bool
	exhaust() error
}

func Allow(rateLimiter RateLimiter) (*bool, error) {
	return rateLimiter.allow()
}
