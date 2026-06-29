package service

import (
	db "cafenetchi-api/internal/db/generated"
	"cafenetchi-api/internal/model"
	"cafenetchi-api/internal/utils"
	"context"
	"errors"
	"log"
	"time"
)

// Auth or AuthService
type Auth struct {
	userQueries *db.Queries
	otpSvc      OTP
	smsSvc      SMS
	jwtSecret   string
}

func NewAuth(userQueries *db.Queries, otpSvc OTP, smsSvc SMS, jwtSecret string) *Auth {
	return &Auth{
		userQueries: userQueries,
		otpSvc:      otpSvc,
		smsSvc:      smsSvc,
		jwtSecret:   jwtSecret,
	}

}

func (s *Auth) SendOTP(ctx context.Context, phone string) error {
	// TODO: Business rule: Maybe check rate limiting here later
	otpCode, err := s.otpSvc.GenerateOTP(ctx, phone)
	if err != nil {
		return err
	}

	log.Printf("generate code for %s is %s", phone, otpCode)

	return s.smsSvc.SendOTP(phone, otpCode)
}

func (s *Auth) ValidateOTP(ctx context.Context, phone, code string) (*model.User, string, bool, error) {
	// verify OTP
	if !s.otpSvc.ValidateOTP(ctx, phone, code) {
		return nil, "", false, errors.New("Invalid otp code: " + code)
	}
	// TODO: Find or Create User

	token, err := utils.GenerateJWT(999, phone, code, s.jwtSecret, time.Second*3600*24)
	if err != nil {
		return nil, "", false, err
	}

	return &model.User{}, token, false, nil

}
