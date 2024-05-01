package limitless

import (
	"testing"
	"time"
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
