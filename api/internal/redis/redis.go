package redis

import (
	"context"
	"fmt"
	"log"

	rds "github.com/redis/go-redis/v9"
)

var client *rds.Client

// addr: e.g 127.0.0.1:6379 || localhost:6379 || redis:6379
// db:0 is default
func Init(ctx context.Context, addr, pass string, db int) {
	// init client
	c := rds.NewClient(&rds.Options{
		Addr:     addr,
		Password: pass,
		DB:       db,
	})

	// Test connection
	_, err := c.Ping(context.Background()).Result()
	if err != nil {
		log.Fatalf("Failed to connect to Redis: %v", err)
		return
	}

	client = c
	fmt.Println("Redis connected successfully")
}

func Close() error {
	if client != nil {
		return client.Close()
	}
	return nil
}
