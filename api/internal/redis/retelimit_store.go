package redis

import (
	"context"
	"fmt"
	"time"
)

const rateLimitPrefix = "ratelimit"

type RedisLimiter struct {
	client *Client
	scope  string
	limit  int64
	window time.Duration
}

func NewRedisLimiter(
	client *Client,
	scope string,
	limit int64,
	window time.Duration,
) *RedisLimiter {
	return &RedisLimiter{
		client: client,
		scope:  scope,
		limit:  limit,
		window: window,
	}
}

func rateLimitKey(scope, key string) string {
	return fmt.Sprintf("%s:%s:%s", rateLimitPrefix, scope, key)
}

func (l *RedisLimiter) Allow(ctx context.Context, key string) (bool, time.Duration, error) {
	keyString := rateLimitKey(l.scope, key)

	count, err := l.client.rdsClient.Incr(ctx, keyString).Result()
	if err != nil {
		return false, 0, err
	}

	// First request: start the expiration window
	if count == 1 {
		if err := l.client.rdsClient.Expire(ctx, key, l.window).Err(); err != nil {
			return false, 0, err
		}
	}

	// Still within the limit
	if count <= l.limit {
		return true, 0, nil
	}

	// Over the limit: tell the caller how long to wait
	ttl, err := l.client.rdsClient.TTL(ctx, key).Result()
	if err != nil {
		return false, 0, err
	}

	return false, ttl, nil
}
