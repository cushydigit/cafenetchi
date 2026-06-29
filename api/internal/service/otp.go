package service

import (
	"cafenetchi-api/internal/redis"
	"context"
	"errors"
	"fmt"
	"math/rand"
)

type OTP interface {
	GenerateOTP(ctx context.Context, phone string) (string, error)
	ValidateOTP(ctx context.Context, phone, otp string) (bool, error)
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

func (s *InRedisOTP) ValidateOTP(ctx context.Context, phone, otp string) (bool, error) {
	storedOTP, err := redis.GetOTP(ctx, phone)
	if err != nil {
		return false, errors.New("Failed to retrieve OTP")
	}
	if storedOTP == otp {
		if err := redis.DeleteOTP(ctx, phone); err != nil {
			return false, errors.New("Failed to delete OTP")
		}
		return true, nil
	}
	return false, nil
}
