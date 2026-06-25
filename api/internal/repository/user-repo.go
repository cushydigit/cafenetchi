package repository

import (
	"cafenetchi-api/internal/model"
	"context"
)

type UserRepository interface {
	List(ctx context.Context) ([]model.User, error)
	Get(ctx context.Context, id int64) (*model.User, error)
	Create(ctx context.Context, user model.User) error
	Delete(ctx context.Context, id int64) error
}
