package service

import (
	"cafenetchi-api/internal/model"
	"cafenetchi-api/internal/repository"
	"cafenetchi-api/internal/utils"
	"errors"
	"log"
	"time"
)

// Auth or AuthService
type Auth struct {
	userRepo   repository.UserRepository
	otpSvc     OTP
	smsService SMS
	jwtSecret  string
}

func NewAuth(userRep repository.UserRepository, otpSvc OTP, smsSvc SMS, jwtSecret string) *Auth {
	return &Auth{
		userRepo:   userRep,
		smsService: smsSvc,
		jwtSecret:  jwtSecret,
	}

}

func (s *Auth) SendOTP(phone string) error {
	// TODO: Business rule: Maybe check rate limiting here later
	otpCode, err := s.otpSvc.GenerateOTP(phone)
	if err != nil {
		return err
	}

	log.Printf("generate code for %s is %s", phone, otpCode)

	return s.smsService.SendOTP(phone, otpCode)
}

func (s *Auth) ValidateOTP(phone, code string) (*model.User, string, bool, error) {
	// verify OTP
	if !s.otpSvc.ValidateOTP(phone, code) {
		return nil, "", false, errors.New("Invalid otp code: " + code)
	}
	// TODO: Find or Create User

	token, err := utils.GenerateJWT(999, phone, code, s.jwtSecret, time.Second*3600*24)
	if err != nil {
		return nil, "", false, err
	}

	return &model.User{}, token, false, nil

}
