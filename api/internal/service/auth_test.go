package service

import (
	"cafenetchi-api/internal/logger"
	"cafenetchi-api/internal/model"
	"cafenetchi-api/internal/otp"
	"cafenetchi-api/internal/repository"
	"cafenetchi-api/internal/sms"
	"cafenetchi-api/internal/types"
	"context"
	"errors"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

var (
	PHONE = "09123456789"
	CODE  = "12345"
)

type fakeUserRepo struct {
	getUser *model.User
	getErr  error

	createUser *model.User
	createErr  error
}

func (f *fakeUserRepo) GetByPhone(ctx context.Context, phone string) (*model.User, error) {
	return f.getUser, f.getErr
}

func (f *fakeUserRepo) Create(ctx context.Context, phone string) (*model.User, error) {
	return f.createUser, f.createErr
}

func (f *fakeUserRepo) Get(ctx context.Context, id int64) (*model.User, error) {
	panic("not implemented")
}

func (f *fakeUserRepo) Update(ctx context.Context, id int64, data model.UpdateUser) (*model.User, error) {
	panic("not implemented")
}

type fakeOTP struct {
	validateErr error
	generateOTP string
	generateErr error
}

func (f *fakeOTP) Generate(ctx context.Context, phone string) (string, error) {
	return f.generateOTP, f.generateErr
}

func (f *fakeOTP) Validate(ctx context.Context, phone, code string) error {
	return f.validateErr
}

func (f *fakeOTP) Consume(ctx context.Context, phone string) error {
	return nil
}

type fakeSMS struct {
	sendErr error
}

func (f *fakeSMS) Send(phone, code string) error {
	return f.sendErr
}

func (f *fakeSMS) SendCustom(phone, msg string) error {
	return nil
}

func TestSendOTP_Success(t *testing.T) {
	s := NewAuth(
		&fakeUserRepo{},
		&fakeOTP{},
		&fakeSMS{},
		"secret",
		logger.NewDefault(),
	)

	err := s.SendOTP(context.Background(), PHONE)
	assert.NoError(t, err)
}

func TestSendOTP_PhoneRequired(t *testing.T) {
	s := NewAuth(
		&fakeUserRepo{},
		&fakeOTP{},
		&fakeSMS{},
		"secret",
		logger.NewDefault(),
	)

	err := s.SendOTP(context.Background(), "")
	assert.ErrorIs(t, err, types.ErrPhoneRequired)
}

func TestSendOTP_GeneratedErr(t *testing.T) {
	s := NewAuth(
		&fakeUserRepo{},
		&fakeOTP{
			generateErr: errors.New("boom"),
		},
		&fakeSMS{},
		"secret",
		logger.NewDefault(),
	)

	err := s.SendOTP(context.Background(), PHONE)
	assert.Error(t, err)
}

func TestSendOTP_SMSFailure(t *testing.T) {
	s := NewAuth(
		&fakeUserRepo{},
		&fakeOTP{
			generateOTP: CODE,
		},
		&fakeSMS{
			sendErr: sms.ErrSend,
		},
		"secret",
		logger.NewDefault(),
	)

	err := s.SendOTP(context.Background(), PHONE)

	assert.ErrorIs(t, err, types.ErrInternalServer)
}

func TestValidateOTP_ExistingUser(t *testing.T) {
	user := &model.User{
		ID:    1,
		Phone: PHONE,
	}

	s := NewAuth(
		&fakeUserRepo{
			getUser: user,
		},
		&fakeOTP{},
		&fakeSMS{},
		"secret",
		logger.NewDefault(),
	)

	result, err := s.ValidateOTP(context.Background(), user.Phone, "123456")
	require.NoError(t, err)

	assert.False(t, result.IsNewUser)
	assert.Equal(t, user.ID, result.User.ID)
	assert.NotEmpty(t, result.Token)
}

func TestValidateOTP_NewUser(t *testing.T) {
	repo := &fakeUserRepo{
		getErr: repository.ErrUserNotFound,
		createUser: &model.User{
			ID:    1,
			Phone: PHONE,
		},
	}

	s := NewAuth(
		repo,
		&fakeOTP{},
		&fakeSMS{},
		"secret",
		logger.NewDefault(),
	)

	result, err := s.ValidateOTP(context.Background(), PHONE, CODE)
	require.NoError(t, err)

	assert.True(t, result.IsNewUser)

}

func TestValidateOTP_InvalidOTP(t *testing.T) {
	s := NewAuth(
		&fakeUserRepo{},
		&fakeOTP{
			validateErr: otp.ErrInvalid,
		},
		&fakeSMS{},
		"secret",
		logger.NewDefault(),
	)

	_, err := s.ValidateOTP(context.Background(), PHONE, CODE)
	assert.ErrorIs(t, err, types.ErrInvalidOTP)

}
