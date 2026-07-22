package sms

import (
	"fmt"

	kvn "github.com/kavenegar/kavenegar-go"
)

type Service interface {
	Send(phone, otp string) error
	SendCustom(phone, message string) error
}

type kavenegar struct {
	client *kvn.Kavenegar
	sender string
}

// implements SMS service
func NewKavenegar(apiKey, sender string) Service {
	return &kavenegar{
		client: kvn.New(apiKey),
		sender: sender,
	}
}

func (s *kavenegar) Send(phone, code string) error {
	message := fmt.Sprintf("کد ورود شما به کافه نت: %s\n\nاین کد فقط ۲ دقیقه معتبر است.", code)
	return s.send(phone, message)
}

func (s *kavenegar) SendCustom(phone, message string) error {
	return s.send(phone, message)
}

func (s *kavenegar) send(phone, message string) error {
	if _, err := s.client.Message.Send(s.sender, []string{phone}, message, nil); err != nil {
		return fmt.Errorf("%w: %v", ErrSend, err)
	}
	return nil
}
