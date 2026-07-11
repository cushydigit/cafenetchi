package service

import (
	db "cafenetchi-api/internal/db/generated"
	"cafenetchi-api/internal/logger"
	"cafenetchi-api/internal/mapper"
	"cafenetchi-api/internal/model"
	"cafenetchi-api/internal/utils"
	"context"
	"database/sql"
	"errors"
	"time"
)

type Auth interface {
	SendOTP(ctx context.Context, phone string) error
	ValidateOTP(ctx context.Context, phone, code string) (*model.User, string, bool, error)
}

// Auth or AuthService
type auth struct {
	userQueries *db.Queries
	otpSvc      OTP
	smsSvc      SMS
	jwtSecret   string
	logger      *logger.Logger
}

func NewAuth(uq *db.Queries, o OTP, s SMS, jwtSecret string, l *logger.Logger) Auth {
	return &auth{
		userQueries: uq,
		otpSvc:      o,
		smsSvc:      s,
		jwtSecret:   jwtSecret,
		logger:      l,
	}

}

func (s *auth) SendOTP(ctx context.Context, phone string) error {
	// TODO: Business rule: Maybe check rate limiting here later
	otpCode, err := s.otpSvc.GenerateOTP(ctx, phone)
	if err != nil {
		return err
	}

	return s.smsSvc.SendOTP(phone, otpCode)
}

func (s *auth) ValidateOTP(ctx context.Context, phone, code string) (*model.User, string, bool, error) {
	// verify OTP
	validated, err := s.otpSvc.ValidateOTP(ctx, phone, code)
	if err != nil {
		return nil, "", false, err
	}
	if !validated {
		return nil, "", false, errors.New("invalid otp code")
	}
	// TODO: Find or Create User
	var isNewUser bool
	var user db.User
	user, err = s.userQueries.GetUserByPhone(ctx, phone)
	if err != nil {
		if !errors.Is(err, sql.ErrNoRows) {
			return nil, "", false, err
		}
		// user doesn't exist
		user, err = s.userQueries.CreateUser(ctx, db.CreateUserParams{
			Phone: phone,
		})

		if err != nil {
			return nil, "", false, err
		}
		isNewUser = true
	}

	token, err := utils.GenerateJWT(
		user.ID,
		user.Phone,
		"user",
		s.jwtSecret,
		time.Second*3600*24,
	)
	if err != nil {
		return nil, "", false, err
	}

	return mapper.DBUserToModel(&user), token, isNewUser, nil

}
