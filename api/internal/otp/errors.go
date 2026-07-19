package otp

import "errors"

var (

	// ErrNotFound indicates that no OTP exists for the specified phone number.
	ErrNotFound = errors.New("otp not found")

	// ErrInvalid indicated that the provided OTP does not match the stored value.
	ErrInvalid = errors.New("otp is invalid")
)
