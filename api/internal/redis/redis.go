package redis

import (
	"context"
	"log"

	rds "github.com/redis/go-redis/v9"
)

type Client struct {
	*rds.Client
}

// addr: e.g 127.0.0.1:6379 || localhost:6379 || redis:6379
// New function will create a new redis client and ping to test connection before return it.
// context will use as context for ping
func New(ctx context.Context, addr, pass string, db int) (*Client, error) {
	// create client
	client := rds.NewClient(&rds.Options{
		Addr:     addr,
		Password: pass,
		DB:       db,
	})

	// test connection
	_, err := client.Ping(ctx).Result()
	if err != nil {
		log.Fatalf("Failed to connect to Redis: %v", err)
		return nil, err
	}

	return &Client{client}, nil
}
