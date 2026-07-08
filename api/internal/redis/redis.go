package redis

import (
	"context"

	rds "github.com/redis/go-redis/v9"
)

type Client struct {
	rdsClient *rds.Client
}

// addr: e.g 127.0.0.1:6379 || localhost:6379 || redis:6379
// NewClient function will create a new redis client and ping to test connection before return it.
// context will use as context for ping
func New(ctx context.Context, addr, pass string, db int) (*Client, error) {
	// create client
	client := rds.NewClient(&rds.Options{
		Addr:     addr,
		Password: pass,
		DB:       db,
	})

	// test connection
	if _, err := client.Ping(ctx).Result(); err != nil {
		_ = client.Close()
		return nil, err
	}

	return &Client{rdsClient: client}, nil
}

func (c *Client) Close(ctx context.Context) error {
	return c.Close(ctx)
}
