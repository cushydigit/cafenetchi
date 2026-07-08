package service

import (
	"context"
	"crypto/rand"
	"encoding/binary"
	"errors"
	"fmt"
)

type OTPStore interface {
	SetOTP(ctx context.Context, phone, otp string) error
	GetOTP(ctx context.Context, phone string) (string, error)
	DelOTP(ctx context.Context, phone string) error
}

type OTP interface {
	GenerateOTP(ctx context.Context, phone string) (string, error)
	ValidateOTP(ctx context.Context, phone, otp string) (bool, error)
}

type RedisOTP struct {
	store OTPStore
}

func NewRedisOTP(store OTPStore) *RedisOTP {
	return &RedisOTP{store: store}
}

func (s *RedisOTP) GenerateOTP(ctx context.Context, phone string) (string, error) {
	var b [4]byte
	if _, err := rand.Read(b[:]); err != nil {
		return "", err
	}

	n := binary.BigEndian.Uint32(b[:]) % 1000000

	return fmt.Sprintf("%06d", n), nil
}

func (s *RedisOTP) ValidateOTP(ctx context.Context, phone, otp string) (bool, error) {
	storedOTP, err := s.store.GetOTP(ctx, phone)
	if err != nil {
		return false, errors.New("Failed to retrieve OTP")
	}
	if storedOTP == otp {
		if err := s.store.DelOTP(ctx, phone); err != nil {
			return false, errors.New("Failed to delete OTP")
		}
		return true, nil
	}
	return false, nil
}
