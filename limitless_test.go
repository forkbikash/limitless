package limitless

import (
	"context"
	"fmt"
	"testing"
	"time"

	"github.com/redis/go-redis/v9"
)

func TestInMemoryImpl(t *testing.T) {
	limiter := NewInMemoryTokenBucket(5, 2)
	executeExample(t, limiter, 5, 2)
}

func TestRedisImpl(t *testing.T) {
	ctx := context.Background()
	client, err := newRedisClient(ctx)
	if err != nil {
		fmt.Println(err)
		return
	}
	defer client.Close()

	limiter, err := NewRedisTokenBucket(ctx, client, "myTokenBucket", 5, 2)
	if err != nil {
		fmt.Println(err)
		return
	}
	executeExample(t, limiter, 5, 2)
}

func executeExample(t *testing.T, limiter RateLimiter, capacity int64, rate int) {
	// rate limit does not exceed
	failed := false
	for i := 0; i < 10; i++ {
		allow, err := Allow(limiter)
		if err != nil {
			fmt.Println(err)
			return
		}

		if allow != nil && *allow {
			fmt.Printf("Operation allowed")
		} else {
			fmt.Printf("Rate limit exceeded")
			failed = true
		}

		time.Sleep(time.Duration(rate) * time.Second)
	}
	if failed {
		t.Fail()
	}

	// rate limit exceeds
	failed = false
	for i := 0; i < 10; i++ {
		allow, err := Allow(limiter)
		if err != nil {
			fmt.Println(err)
			return
		}

		if allow != nil && *allow {
			fmt.Printf("Operation allowed")
		} else {
			failed = true
			fmt.Printf("Rate limit exceeded")
		}

		time.Sleep(time.Duration(rate-1) * time.Second)
	}
	if !failed {
		t.Fail()
	}
}

func newRedisClient(ctx context.Context) (redis.UniversalClient, error) {
	client := redis.NewUniversalClient(&redis.UniversalOptions{
		Addrs:    []string{"127.0.0.1:6379"},
		Password: "XXXX",
		DB:       0,
	})
	_, err := client.Ping(ctx).Result()
	if err != nil {
		return nil, err
	}

	return client, nil
}
