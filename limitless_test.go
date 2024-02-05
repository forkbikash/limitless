package limitless

import (
	"context"
	"fmt"
	"testing"
	"time"

	"github.com/go-redis/redis/v8"
)

func TestInMemoryImpl(t *testing.T) {
	limiter := NewInMemoryTokenBucket(5, 2)
	doExample(t, limiter, 5, 2)
}

func TestRedisImpl(t *testing.T) {
	ctx := context.Background()
	client, err := newRedisClient(ctx)
	if err != nil {
		fmt.Println(err)
		return
	}
	defer client.Close()
	// environment.InitializeEnvs(environment.VARIANT_TEST)
	// logger.InitializeLogger()
	// piceredis.InitializeRedisMap([]piceredis.EnumRedisDb{piceredis.REDIS_CREDIT_DB})

	limiter, err := NewRedisTokenBucket(ctx, client, "myTokenBucket", 5, 2)
	// limiter, err := NewRedisTokenBucket(ctx, piceredis.DefaultClient(&ctx), "myTokenBucket", 5, 2)
	if err != nil {
		fmt.Println(err)
		return
	}
	doExample(t, limiter, 5, 2)
}

func doExample(t *testing.T, limiter RateLimiter, capacity int64, rate int) {
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

func newRedisClient(ctx context.Context) (*redis.Client, error) {
	client := redis.NewClient(&redis.Options{
		Addr:     "127.0.0.1:6379",
		Password: "XXXX",
		DB:       0,
	})

	_, err := client.Ping(ctx).Result()
	if err != nil {
		return nil, err
	}

	return client, nil
}
