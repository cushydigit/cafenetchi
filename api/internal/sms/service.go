package sms

import (
	"fmt"

	"github.com/kavenegar/kavenegar-go"
)

type SMSService interface {
	SendOTP(phone, otp string) error
	SendCustom(phone, message string) error
}

// KavenegarService implements SMSService
type KavenegarService struct {
	client *kavenegar.Kavenegar
	sender string
}

func NewKavenegarService(apiKey, sender string) SMSService {
	return &KavenegarService{
		client: kavenegar.New(apiKey),
		sender: sender,
	}
}

func (s *KavenegarService) SendOTP(phone, otp string) error {
	message := fmt.Sprintf("کد ورود شما به کافه نت: %s\n\nاین کد فقط ۲ دقیقه معتبر است.", otp)

	_, err := s.client.Message.Send(s.sender, []string{phone}, message, nil)
	if err != nil {
		return fmt.Errorf("failed to send SMS: %w", err)
	}
	return nil
}

func (s *KavenegarService) SendCustom(phone, message string) error {
	_, err := s.client.Message.Send(s.sender, []string{phone}, message, nil)
	return err
}
