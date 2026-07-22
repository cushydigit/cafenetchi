package otp

import (
	"context"
	"crypto/rand"
	"encoding/binary"
	"errors"
	"fmt"
)

// OTP provides operations for generating and validating one-time passwords
type Service interface {
	Generate(ctx context.Context, phone string) (string, error)
	Validate(ctx context.Context, phone, code string) error
	Consume(ctx context.Context, phone string) error
}

type service struct {
	store Store
}

// New creates a new OTP service using the given store.
//
// The store is responsible for persisting OTP codes and may be backed by Redis, in-memory implementation, or any other storage backend.
func New(s Store) Service {
	return &service{
		store: s,
	}
}

// Generate generates a random 6-digit OTP, stores it for the
// specified phone number, and returns the generated code.
func (s *service) Generate(ctx context.Context, phone string) (string, error) {
	var b [4]byte
	if _, err := rand.Read(b[:]); err != nil {
		return "", fmt.Errorf("generation otp: %w", err)
	}

	code := fmt.Sprintf("%06d", binary.BigEndian.Uint32(b[:])%1000000)

	if err := s.store.Set(ctx, phone, code); err != nil {
		return "", fmt.Errorf("%w: %w", ErrStoreSet, err)
	}
	return code, nil
}

// Validate verifies the provided OTP against the stored value.
func (s *service) Validate(ctx context.Context, phone, code string) error {
	storedOTP, err := s.store.Get(ctx, phone)
	if err != nil {
		if errors.Is(err, ErrNotFound) {
			return ErrNotFound
		}
		return fmt.Errorf("%w: %w", ErrStoreGet, err)
	}
	if storedOTP != code {
		return ErrInvalid
	}
	return nil
}

// Consume will delete stored OTP and cannot be used again.
func (s *service) Consume(ctx context.Context, phone string) error {
	if err := s.store.Del(ctx, phone); err != nil {
		return fmt.Errorf("%w: %w", ErrStoreDel, err)
	}
	return nil
}
