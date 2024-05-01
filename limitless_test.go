package limitless

import (
	"context"
	"testing"
	"time"

	"github.com/redis/go-redis/v9"
)

func TestInMemoryTokenBucketAllow(t *testing.T) {
	// Test case: Allow with sufficient tokens
	tb := NewInMemoryTokenBucket(10, 2)
	allowed, err := Allow(tb)
	if err != nil {
		t.Errorf("Unexpected error: %v", err)
	}
	if !*allowed {
		t.Error("Expected to allow with sufficient tokens")
	}

	// Test case: Disallow when no tokens are available
	for i := 0; i < 10; i++ {
		_, _ = Allow(tb)
	}
	allowed, err = Allow(tb)
	if err != nil {
		t.Errorf("Unexpected error: %v", err)
	}
	if *allowed {
		t.Error("Expected to disallow when no tokens are available")
	}

	// Test case: Allow after refilling tokens
	time.Sleep(5 * time.Second)
	allowed, err = Allow(tb)
	if err != nil {
		t.Errorf("Unexpected error: %v", err)
	}
	if !*allowed {
		t.Error("Expected to allow after refilling tokens")
	}
}

func TestRedisTokenBucketAllow(t *testing.T) {
	client := redis.NewClient(&redis.Options{
		Addr: "localhost:6379",
	})
	defer client.Close()

	ctx := context.Background()
	key := "test_key"

	// Test case: Allow with sufficient tokens
	tb, err := NewRedisTokenBucket(ctx, client, key, 10, 2)
	if err != nil {
		t.Errorf("Unexpected error: %v", err)
	}
	allowed, err := Allow(tb)
	if err != nil {
		t.Errorf("Unexpected error: %v", err)
	}
	if !*allowed {
		t.Error("Expected to allow with sufficient tokens")
	}

	// Test case: Disallow when no tokens are available
	for i := 0; i < 10; i++ {
		_, _ = Allow(tb)
	}
	allowed, err = Allow(tb)
	if err != nil {
		t.Errorf("Unexpected error: %v", err)
	}
	if *allowed {
		t.Error("Expected to disallow when no tokens are available")
	}

	// Test case: Allow after refilling tokens
	time.Sleep(5 * time.Second)
	allowed, err = Allow(tb)
	if err != nil {
		t.Errorf("Unexpected error: %v", err)
	}
	if !*allowed {
		t.Error("Expected to allow after refilling tokens")
	}

	// Clean up Redis
	err = client.Del(ctx, key).Err()
	if err != nil {
		t.Errorf("Unexpected error cleaning up Redis: %v", err)
	}
}
