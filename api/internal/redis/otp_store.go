package redis

import (
	"cafenetchi-api/internal/otp"
	"context"
	"errors"
	"fmt"
	"time"

	rds "github.com/redis/go-redis/v9"
)

type RedisOTPStore struct {
	client *Client
	ttl    time.Duration
}

func NewOTPStore(client *Client, ttl time.Duration) *RedisOTPStore {
	return &RedisOTPStore{
		client: client,
		ttl:    ttl,
	}
}

// otpKey generates a key for storing OTPs in Redis.
func otpKey(phone string) string {
	return fmt.Sprintf("otp:%s", phone)
}

// SetOTP sets the OTP for a given phone number in Redis.
func (c *RedisOTPStore) Set(ctx context.Context, phone, otp string) error {
	return c.client.rdsClient.Set(ctx, otpKey(phone), otp, c.ttl).Err()
}

// GetOTP retrieves the OTP (One-Time Password) associated with a given phone number from Redis.
func (c *RedisOTPStore) Get(ctx context.Context, phone string) (string, error) {
	code, err := c.client.rdsClient.Get(ctx, otpKey(phone)).Result()
	if err != nil {
		if errors.Is(err, rds.Nil) {
			return "", otp.ErrNotFound
		}
		return "", err
	}
	return code, nil
}

// DeleteOTP deletes the OTP (One-Time Password) associated with a given phone number from Redis.
func (c *RedisOTPStore) Del(ctx context.Context, phone string) error {
	return c.client.rdsClient.Del(ctx, otpKey(phone)).Err()
}
