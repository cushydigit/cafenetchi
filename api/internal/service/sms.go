package service

type SMS interface {
	SendOTP(phone, otp string) error
	SendCustom(phone, message string) error
}
