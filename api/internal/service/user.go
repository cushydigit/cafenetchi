package service

import (
	"cafenetchi-api/internal/logger"
	"cafenetchi-api/internal/model"
	"cafenetchi-api/internal/repository"
	"context"
)

// TODO: change req type to input for service layer?
type User interface {
	GetByID(ctx context.Context, id int64) (*model.User, error)
	UpdateProfile(ctx context.Context, id int64, data model.UpdateUser) (*model.User, error)
}

type user struct {
	userRepo repository.User
	logger   *logger.Logger
}

func NewUser(r repository.User, l *logger.Logger) User {
	return &user{
		userRepo: r,
		logger:   l,
	}
}

func (s *user) GetByID(ctx context.Context, id int64) (*model.User, error) {
	u, err := s.userRepo.GetByID(ctx, id)
	if err != nil {
		return nil, err
	}
	return u, nil
}

func (s *user) UpdateProfile(ctx context.Context, id int64, data model.UpdateUser) (*model.User, error) {
	u, err := s.userRepo.Update(ctx, id, data)
	if err != nil {
		return nil, err
	}
	return u, nil
}
