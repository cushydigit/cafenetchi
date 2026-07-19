package otp

import "context"

// Store defines the persistence layer used to store OTP code.
// Implementation may use Redis, an in-memory store, or any other backend
type Store interface {
	// Set stores an OTP for a phone number.
	Set(ctx context.Context, phone, otp string) error

	// Get returns the stored OTP.
	//
	// It returns ErrNotFound if the OTP is not found.
	Get(ctx context.Context, phone string) (string, error)

	// Delete removes the stored OTP
	Del(ctx context.Context, phone string) error
}
