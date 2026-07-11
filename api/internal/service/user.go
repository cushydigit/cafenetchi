package service

import (
	db "cafenetchi-api/internal/db/generated"
	"cafenetchi-api/internal/logger"
	"cafenetchi-api/internal/mapper"
	"cafenetchi-api/internal/model"
	"cafenetchi-api/internal/types"
	"context"
)

// TODO: change req type to input for service layer?
type User interface {
	GetByID(ctx context.Context, id int64) (*model.User, error)
	UpdateProfile(ctx context.Context, id int64, req types.UpdateProfileRequest) error
}

type user struct {
	logger  *logger.Logger
	queries *db.Queries
}

func NewUser(l *logger.Logger, q *db.Queries) User {
	u := &user{
		logger:  l,
		queries: q,
	}
	return u
}

func (s *user) GetByID(ctx context.Context, id int64) (*model.User, error) {
	u, err := s.queries.GetUserByID(ctx, id)
	if err != nil {
		return nil, err
	}

	return mapper.DBUserToModel(&u), nil

}

func (s *user) UpdateProfile(ctx context.Context, id int64, req types.UpdateProfileRequest) error {

	s.queries.UpdateUser(ctx, db.UpdateUserParams{
		ID:        id,
		FullName:  mapper.StringToPgText(req.FullName),
		AvatarUrl: mapper.StringToPgText(req.Avatar),
	})
	return nil

}
