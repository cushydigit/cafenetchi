package redis

import (
	"context"
	"fmt"
	"time"
)

const (
	OTP_TTL = 2 * time.Minute
)

// otpKey generates a key for storing OTPs in Redis.
//
// Parameters:
// - phone: The phone number associated with the OTP.
//
// Return:
// - string: The generated OTP key.
func otpKey(phone string) string {
	return fmt.Sprintf("otp:%s", phone)
}

// SetOTP sets the OTP for a given phone number in Redis.
//
// Parameters:
// - ctx: The context.Context object for the function.
// - phone: The phone number to set the OTP for.
// - otp: The OTP to set for the phone number.
//
// Return:
// - error: An error if there was a problem setting the OTP.
func (c *Client) SetOTP(ctx context.Context, phone, otp string) error {
	return c.Set(ctx, otpKey(phone), otp, OTP_TTL).Err()
}

// GetOTP retrieves the OTP (One-Time Password) associated with a given phone number from Redis.
//
// Parameters:
// - ctx: The context.Context object for the function.
// - phone: The phone number to retrieve the OTP for.
//
// Returns:
// - string: The OTP associated with the phone number.
// - error: An error if there was a problem retrieving the OTP.
func (c *Client) GetOTP(ctx context.Context, phone string) (string, error) {
	return c.Get(ctx, otpKey(phone)).Result()
}

// DeleteOTP deletes the OTP (One-Time Password) associated with a given phone number from Redis.
//
// Parameters:
// - ctx: The context.Context object for the function.
// - phone: The phone number to delete the OTP for.
//
// Returns:
// - error: An error if there was a problem deleting the OTP.
func (c *Client) DeleteOTP(ctx context.Context, phone string) error {
	return c.Del(ctx, otpKey(phone)).Err()
}
