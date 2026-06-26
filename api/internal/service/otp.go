package service

import (
	"cafenetchi-api/internal/redis"
	"context"
	"fmt"
	"math/rand"
)

type OTP interface {
	GenerateOTP(ctx context.Context, phone string) (string, error)
	ValidateOTP(ctx context.Context, phone, otp string) bool
}

type InRedisOTP struct{}

func NewInRedisOTP() *InRedisOTP {
	return &InRedisOTP{}
}

func (s *InRedisOTP) GenerateOTP(ctx context.Context, phone string) (string, error) {
	otp := fmt.Sprintf("%06d", rand.Intn(1000000))
	if err := redis.SetOTP(ctx, phone, otp); err != nil {
		return "", err
	}

	return otp, nil
}

func (s *InRedisOTP) ValidateOTP(ctx context.Context, phone, otp string) bool {
	storedOTP, err := redis.GetOTP(ctx, phone)
	if err != nil {
		return false
	}
	if storedOTP == otp {
		redis.DeleteOTP(ctx, phone)
		return true
	}
	return false
}
