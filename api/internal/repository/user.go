package repository

import (
	db "cafenetchi-api/internal/db/generated"
	"cafenetchi-api/internal/mapper"
	"cafenetchi-api/internal/model"
	"context"
	"database/sql"
	"errors"

	"github.com/jackc/pgx/v5/pgconn"
	"github.com/jackc/pgx/v5/pgtype"
)

type User interface {
	// GetByID returns the user with given ID.
	// Returns ErrUserNotFound if no matching user exits.
	GetByID(ctx context.Context, id int64) (*model.User, error)
	// GetBeyPhone return the user with given ID.
	// Returns ErrUserNotFound if not matching user exits.
	GetByPhone(ctx context.Context, phone string) (*model.User, error)

	// Create creates a new user using given phone number.
	// Returns ErrPhoneAlreadyExists if matching phone already exists.
	Create(ctx context.Context, phone string) (*model.User, error)

	// Update update the user with given id and data
	// Returns updated user
	Update(ctx context.Context, id int64, data model.UpdateUser) (*model.User, error)
}

type user struct {
	queries *db.Queries
}

func NewUser(q *db.Queries) User {
	return &user{queries: q}
}

func (r *user) GetByID(ctx context.Context, id int64) (*model.User, error) {
	u, err := r.queries.GetUserByID(ctx, id)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, ErrUserNotFound
		}
	}
	return mapper.DBUserToModel(&u), nil
}

func (r *user) GetByPhone(ctx context.Context, phone string) (*model.User, error) {
	u, err := r.queries.GetUserByPhone(ctx, phone)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, ErrUserNotFound
		}
		return nil, err
	}
	return mapper.DBUserToModel(&u), nil
}

func (r *user) Create(ctx context.Context, phone string) (*model.User, error) {
	u, err := r.queries.CreateUser(ctx, db.CreateUserParams{
		Phone: phone,
	})
	if err != nil {
		var pgErr *pgconn.PgError
		if errors.As(err, &pgErr) {
			switch pgErr.Code {
			case "23505":
				return nil, ErrPhoneExists
			}
		}
		return nil, err
	}
	return mapper.DBUserToModel(&u), nil
}

func (r *user) Update(ctx context.Context, id int64, data model.UpdateUser) (*model.User, error) {
	u, err := r.queries.UpdateUser(ctx, db.UpdateUserParams{
		ID: id,
		FullName: pgtype.Text{
			String: data.FullName,
			Valid:  data.FullName != "",
		},
		AvatarUrl: pgtype.Text{
			String: data.AvatarURL,
			Valid:  data.AvatarURL != "",
		},
	})
	if err != nil {
		return nil, err
	}
	return mapper.DBUserToModel(&u), nil
}
