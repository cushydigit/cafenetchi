package service

import (
	"cafenetchi-api/internal/redis"
	"context"
	"fmt"
	"math/rand"
)

type OTP interface {
	GenerateOTP(phone string) (string, error)
	ValidateOTP(phone, otp string) bool
}

type InRedisOTP struct{}

func NewInRedisOTP() *InRedisOTP {
	return &InRedisOTP{}
}

func (s *InRedisOTP) GenerateOTP(phone string) (string, error) {
	otp := fmt.Sprintf("%06d", rand.Intn(1000000))
	if err := redis.SetOTP(context.Background(), phone, otp); err != nil {
		return "", err
	}

	return otp, nil
}

func (s *InRedisOTP) ValidateOTP(phone, otp string) bool {
	storedOTP, err := redis.GetOTP(context.Background(), phone)
	if err != nil {
		return false
	}
	if storedOTP == otp {
		return true
	}
	return false
}
