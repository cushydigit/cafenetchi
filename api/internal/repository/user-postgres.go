package repository

import (
	db "cafenetchi-api/internal/db/generated"
	"cafenetchi-api/internal/model"
	"context"
)

type UserPostgres struct {
	q *db.Queries
}

func NewUserPostgresRepo(q *db.Queries) *UserPostgres {
	return &UserPostgres{q: q}
}

func (r *UserPostgres) List(ctx context.Context) ([]model.User, error) {
	users, err := r.q.ListUsers(ctx)
	if err != nil {
		return nil, err
	}
	var us []model.User
	for _, u := range users {
		us = append(us, model.User{
			ID:    u.ID,
			Phone: u.Phone,
		})
	}
	return us, nil
}

func (r *UserPostgres) Get(ctx context.Context, id int64) (*model.User, error) {
	user, err := r.q.GetUserByID(ctx, id)
	if err != nil {
		return nil, err
	}
	return &model.User{
		ID:    user.ID,
		Phone: user.Phone,
	}, nil
}

func (r *UserPostgres) Create(ctx context.Context, user model.User) error {
	_, err := r.q.CreateUser(ctx, db.CreateUserParams{
		Phone: user.Phone,
	})
	return err
}

func (r *UserPostgres) Delete(ctx context.Context, id int64) error {
	return nil
}
