package sms

import (
	"cafenetchi-api/internal/service"
	"fmt"

	"github.com/kavenegar/kavenegar-go"
)

type Kavenegar struct {
	client *kavenegar.Kavenegar
	sender string
}

// implements SMS service
func NewKavenegar(apiKey, sender string) service.SMS {
	return &Kavenegar{
		client: kavenegar.New(apiKey),
		sender: sender,
	}
}

func (s *Kavenegar) SendOTP(phone, otp string) error {
	message := fmt.Sprintf("کد ورود شما به کافه نت: %s\n\nاین کد فقط ۲ دقیقه معتبر است.", otp)

	_, err := s.client.Message.Send(s.sender, []string{phone}, message, nil)
	if err != nil {
		return fmt.Errorf("failed to send SMS: %w", err)
	}
	return nil
}

func (s *Kavenegar) SendCustom(phone, message string) error {
	_, err := s.client.Message.Send(s.sender, []string{phone}, message, nil)
	return err
}
