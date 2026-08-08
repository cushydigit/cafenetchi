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
	"strings"
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
	log := s.logger.WithContext(ctx)
	// TODO: Business rule: Maybe check rate limiting here later

	// the validation check on the handler but the auth service
	// is application service so we are to protect it.
	if strings.TrimSpace(phone) == "" {
		return types.ErrPhoneRequired
	}
	code, err := s.otpSvc.Generate(ctx, phone)
	if err != nil {
		log.Error(
			"failed to generate otp",
			"phone", phone,
			"error", err,
		)
		return types.ErrInternalServer
	}

	if err := s.smsSvc.Send(phone, code); err != nil {
		log.Error(
			"failed to send otp",
			"phone", phone,
			"error", err,
		)
		return types.ErrInternalServer
	}

	return nil

}

func (s *auth) ValidateOTP(ctx context.Context, phone, code string) (*AuthResult, error) {
	log := s.logger.WithContext(ctx)
	// verify OTP
	if err := s.otpSvc.Validate(ctx, phone, code); err != nil {
		switch {
		case errors.Is(err, otp.ErrInvalid):
			return nil, types.ErrInvalidOTP
		case errors.Is(err, otp.ErrNotFound):
			return nil, types.ErrInvalidOTP
		default:
			log.Error(
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
			log.Error(
				"failed to create user",
				"phone",
				phone,
				"error",
				err,
			)

			return nil, types.ErrInternalServer
		}

		isNewUser = true

		log.Info(
			"user registered",
			"user_id", user.ID,
		)

	} else {
		log.Error(
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
		log.Error(
			"failed to generate jwt",
			"user_id", user.ID,
			"error", err,
		)
		return nil, types.ErrInternalServer
	}

	log.Info(
		"user authenticated",
		"user_id", user.ID,
		"is_new", isNewUser,
	)

	// consume or delete OTP
	if err := s.otpSvc.Consume(ctx, phone); err != nil {
		log.Error(
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
