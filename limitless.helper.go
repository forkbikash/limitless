package limitless

import "github.com/go-redis/redis/v8"

func isNonNilErr(err error) bool {
	if err != nil && err != redis.Nil {
		return true
	}

	return false
}

func boolPtr(value bool) *bool {
	return &value
}
