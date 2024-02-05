package limitless

import (
	"context"
	"math"
	"strconv"
	"time"

	"github.com/go-redis/redis/v8"
)

// redis implementation is a better alternative

type redisTokenBucket struct {
	pipe       redis.Pipeliner
	context    context.Context
	client     *redis.Client
	key        string
	capacity   int64
	rate       int
	tokens     int64
	lastUpdate time.Time
}

func NewRedisTokenBucket(ctx context.Context, client *redis.Client, key string, capacity int64, rate int) (*redisTokenBucket, error) {
	pipe := client.TxPipeline()
	defer pipe.Close()

	state, err := pipe.HGetAll(ctx, key).Result()
	if isNonNilErr(err) {
		return nil, err
	}

	if err == redis.Nil {
		lastAccess := time.Now()
		_ = pipe.HSet(ctx, key,
			"last_access", strconv.FormatInt(lastAccess.Unix(), 10),
			"last_tokens", strconv.FormatInt(capacity, 10),
		)

		_, err := pipe.Exec(ctx)
		if err != nil {
			return nil, err
		}

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
	tb.pipe = tb.client.TxPipeline()
	defer tb.pipe.Close()

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

func (tb *redisTokenBucket) load() error {
	state, err := tb.pipe.HGetAll(tb.context, tb.key).Result()
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
	_ = tb.pipe.HSet(tb.context, tb.key,
		"last_access", strconv.FormatInt(now.Unix(), 10),
		"last_tokens", strconv.FormatFloat(float64(tb.tokens), 'f', -1, 64),
	)
	_, err := tb.pipe.Exec(tb.context)
	if err != nil {
		return err
	}
	tb.lastUpdate = now

	return nil
}
