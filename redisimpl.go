package limitless

import (
	"context"
	"errors"
	"math"
	"strconv"
	"time"

	"github.com/redis/go-redis/v9"
)

// redis implementation is a better alternative

type redisTokenBucket struct {
	context    context.Context
	client     redis.UniversalClient
	key        string
	capacity   int64
	rate       int
	tokens     int64
	lastUpdate time.Time
}

func NewRedisTokenBucket(ctx context.Context, client redis.UniversalClient, key string, capacity int64, rate int) (*redisTokenBucket, error) {
	acquired, err := acquireLock(&ctx, client, key)
	if err != nil || !acquired {
		return nil, errors.New("error acquiring lock")
	}
	defer releaseLock(&ctx, client, key)

	state, err := client.HGetAll(ctx, key).Result()
	if isNonNilErr(err) {
		return nil, err
	}

	if err == redis.Nil {
		lastAccess := time.Now()
		_ = client.HSet(ctx, key,
			"last_access", strconv.FormatInt(lastAccess.Unix(), 10),
			"last_tokens", strconv.FormatInt(capacity, 10),
		)

		return &redisTokenBucket{
			context:    ctx,
			client:     client,
			key:        key,
			capacity:   capacity,
			rate:       rate,
			tokens:     capacity,
			lastUpdate: lastAccess,
		}, nil
	}

	lastAccess, err := strconv.ParseInt(state["last_access"], 10, 64)
	if err != nil {
		return nil, err
	}
	lastTokens, err := strconv.ParseInt(state["last_tokens"], 10, 64)
	if err != nil {
		return nil, err
	}

	return &redisTokenBucket{
		context:    ctx,
		client:     client,
		key:        key,
		capacity:   capacity,
		rate:       rate,
		lastUpdate: time.Unix(lastAccess, 0),
		tokens:     lastTokens,
	}, nil
}

func (tb *redisTokenBucket) allow() (*bool, error) {
	acquired, err := acquireLock(&tb.context, tb.client, tb.key)
	if err != nil || !acquired {
		return nil, errors.New("error acquiring lock")
	}
	defer releaseLock(&tb.context, tb.client, tb.key)

	err = tb.load()
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

func (tb *redisTokenBucket) load() error {
	state, err := tb.client.HGetAll(tb.context, tb.key).Result()
	if err != nil {
		return err
	}
	lastAccess, err := strconv.ParseInt(state["last_access"], 10, 64)
	if err != nil {
		return err
	}
	tb.lastUpdate = time.Unix(lastAccess, 0)
	lastTokens, err := strconv.ParseInt(state["last_tokens"], 10, 64)
	if err != nil {
		return err
	}
	tb.tokens = lastTokens
	return nil
}

func (tb *redisTokenBucket) refill() {
	elapsedTime := time.Since(tb.lastUpdate)
	tokensToAdd := math.Floor(elapsedTime.Seconds() * float64(tb.rate))
	if tokensToAdd > 0 {
		newTokens := float64(tb.tokens) + tokensToAdd
		tb.tokens = int64(math.Min(newTokens, float64(tb.capacity)))
	}
}

func (tb *redisTokenBucket) available() *bool {
	if tb.tokens < 1 {
		return boolPtr(false)
	}
	return nil
}

func (tb *redisTokenBucket) exhaust() error {
	if tb.tokens > 0 {
		tb.tokens--
	}
	now := time.Now()
	_ = tb.client.HSet(tb.context, tb.key,
		"last_access", strconv.FormatInt(now.Unix(), 10),
		"last_tokens", strconv.FormatFloat(float64(tb.tokens), 'f', -1, 64),
	)

	tb.lastUpdate = now

	return nil
}
