package service

import (
	"cafenetchi-api/internal/logger"
	"cafenetchi-api/internal/model"
	"cafenetchi-api/internal/otp"
	"cafenetchi-api/internal/repository"
	"cafenetchi-api/internal/sms"
	"cafenetchi-api/internal/types"
	"cafenetchi-api/internal/utils"
	"context"
	"errors"
	"time"
)

type Auth interface {
	SendOTP(ctx context.Context, phone string) error
	ValidateOTP(ctx context.Context, phone, code string) (*AuthResult, error)
}

type AuthResult struct {
	User      *model.User
	Token     string
	IsNewUser bool
}

// Auth or AuthService
type auth struct {
	userRepo  repository.User
	otpSvc    otp.Service
	smsSvc    sms.Service
	jwtSecret string
	logger    *logger.Logger
}

func NewAuth(ur repository.User, o otp.Service, s sms.Service, jwtSecret string, l *logger.Logger) Auth {
	return &auth{
		userRepo:  ur,
		otpSvc:    o,
		smsSvc:    s,
		jwtSecret: jwtSecret,
		logger:    l,
	}

}

func (s *auth) SendOTP(ctx context.Context, phone string) error {
	// TODO: Business rule: Maybe check rate limiting here later

	// the validation check on the handler but the auth service
	// is application service so we are to protect it.
	if phone == "" {
		return types.ErrPhoneRequired
	}
	code, err := s.otpSvc.Generate(ctx, phone)
	if err != nil {
		s.logger.Error(
			"failed to generate otp",
			"phone", phone,
			"error", err,
		)
		return types.ErrInternalServer
	}

	if err := s.smsSvc.Send(phone, code); err != nil {
		s.logger.Error(
			"failed to send otp",
			"phone", phone,
			"error", err,
		)
		return types.ErrInternalServer
	}

	return nil

}

func (s *auth) ValidateOTP(ctx context.Context, phone, code string) (*AuthResult, error) {
	// verify OTP
	if err := s.otpSvc.Validate(ctx, phone, code); err != nil {
		switch {
		case errors.Is(err, otp.ErrInvalid):
			return nil, types.ErrInvalidOTP
		case errors.Is(err, otp.ErrNotFound):
			return nil, types.ErrInvalidOTP
		default:
			s.logger.Error(
				"otp validation failed",
				"phone", phone,
				"error", err,
			)
		}
		return nil, types.ErrInternalServer
	}

	var (
		user      *model.User
		isNewUser bool
		err       error
	)

	user, err = s.userRepo.GetByPhone(ctx, phone)

	if err == nil {
		// existing user
	} else if errors.Is(err, repository.ErrUserNotFound) {

		user, err = s.userRepo.Create(ctx, phone)
		if err != nil {
			s.logger.Error(
				"failed to create user",
				"phone",
				phone,
				"error",
				err,
			)

			return nil, types.ErrInternalServer
		}

		isNewUser = true

		s.logger.Info(
			"user registered",
			"user_id", user.ID,
		)

	} else {
		s.logger.Error(
			"failed to get user",
			"phone", phone,
			"error", err,
		)

		return nil, types.ErrInternalServer
	}

	token, err := utils.GenerateJWT(
		user.ID,
		user.Phone,
		string(model.RoleCustomer),
		s.jwtSecret,
		time.Hour*24,
	)
	if err != nil {
		s.logger.Error(
			"failed to generate jwt",
			"user_id", user.ID,
			"error", err,
		)
		return nil, types.ErrInternalServer
	}

	s.logger.Info(
		"user authenticated",
		"user_id", user.ID,
		"is_new", isNewUser,
	)

	// consume or delete OTP
	if err := s.otpSvc.Consume(ctx, phone); err != nil {
		s.logger.Error(
			"otp error delete",
			"phone", phone,
			"error", err,
		)
		return nil, types.ErrInternalServer

	}

	return &AuthResult{
		User:      user,
		Token:     token,
		IsNewUser: isNewUser,
	}, nil

}
