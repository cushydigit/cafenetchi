package redis

import (
	"context"
	"fmt"
	"time"
)

const (
	OTP_TTL = 2 * time.Minute
)

func otpKey(phone string) string {
	return fmt.Sprintf("otp:%s", phone)
}

func SetOTP(ctx context.Context, phone, otp string) error {
	return client.Set(ctx, otpKey(phone), otp, OTP_TTL).Err()
}

func GetOTP(ctx context.Context, phone string) (string, error) {
	return client.Get(ctx, otpKey(phone)).Result()
}
