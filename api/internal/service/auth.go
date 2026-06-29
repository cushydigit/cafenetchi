package service

import (
	db "cafenetchi-api/internal/db/generated"
	"cafenetchi-api/internal/model"
	"cafenetchi-api/internal/utils"
	"context"
	"errors"
	"log/slog"
	"time"
)

// Auth or AuthService
type Auth struct {
	userQueries *db.Queries
	otpSvc      OTP
	smsSvc      SMS
	jwtSecret   string
	logger      *slog.Logger
}

func NewAuth(userQueries *db.Queries, otpSvc OTP, smsSvc SMS, jwtSecret string, logger *slog.Logger) *Auth {
	return &Auth{
		userQueries: userQueries,
		otpSvc:      otpSvc,
		smsSvc:      smsSvc,
		jwtSecret:   jwtSecret,
		logger:      logger,
	}

}

func (s *Auth) SendOTP(ctx context.Context, phone string) error {
	// TODO: Business rule: Maybe check rate limiting here later
	otpCode, err := s.otpSvc.GenerateOTP(ctx, phone)
	if err != nil {
		return err
	}

	s.logger.Info("code generated for opt", phone, otpCode)

	return s.smsSvc.SendOTP(phone, otpCode)
}

func (s *Auth) ValidateOTP(ctx context.Context, phone, code string) (*model.User, string, bool, error) {
	// verify OTP
	validated, err := s.otpSvc.ValidateOTP(ctx, phone, code)
	if err != nil {
		return nil, "", false, err
	}
	if !validated {
		return nil, "", false, errors.New("invalid otp code")
	}
	// TODO: Find or Create User

	token, err := utils.GenerateJWT(999, phone, code, s.jwtSecret, time.Second*3600*24)
	if err != nil {
		return nil, "", false, err
	}

	return &model.User{}, token, false, nil

}
