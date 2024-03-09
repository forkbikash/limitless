package limitless

import (
	"context"
	"time"

	"github.com/redis/go-redis/v9"
)

func isNonNilErr(err error) bool {
	if err != nil && err != redis.Nil {
		return true
	}

	return false
}

func boolPtr(value bool) *bool {
	return &value
}

const REDIS_LOCK_EXPIRY = 5 * time.Second

func acquireLock(ctx *context.Context, client redis.UniversalClient, key string) (bool, error) {
	result, err := client.SetNX(*ctx, key, "1", REDIS_LOCK_EXPIRY).Result()
	if err != nil {
		return false, err
	}

	return result, nil
}

func releaseLock(ctx *context.Context, client redis.UniversalClient, key string) error {
	_, err := client.Del(*ctx, key).Result()
	return err
}
