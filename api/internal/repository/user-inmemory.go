package repository

import (
	"cafenetchi-api/internal/model"
	"context"
	"errors"
	"sync"
)

type UserInMemory struct {
	users  map[int64]*model.User
	mu     sync.RWMutex
	nextID int64
}

func NewInMemoryUserRepo() *UserInMemory {
	return &UserInMemory{
		users:  make(map[int64]*model.User),
		nextID: 1,
	}

}

func (r *UserInMemory) List(ctx context.Context) ([]model.User, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()
	var users []model.User
	for _, u := range r.users {
		users = append(users, *u)
	}
	return users, nil

}

func (r *UserInMemory) Get(ctx context.Context, id int64) (*model.User, error) {
	r.mu.Lock()
	defer r.mu.Unlock()

	user, ok := r.users[id]
	if !ok {
		return nil, errors.New("User not found in inmemory repository")
	}

	// defensive copy for in-memory
	userCopy := *user
	return &userCopy, nil
}

func (r *UserInMemory) Create(ctx context.Context, user model.User) error {
	r.mu.Lock()
	defer r.mu.Unlock()

	if user.Phone == "" {
		return errors.New("invalid phone number")
	}

	for _, u := range r.users {
		if u.Phone == user.Phone {
			return errors.New("invalid phone number")
		}
	}

	user.ID = uint(r.nextID)
	r.users[int64(user.ID)] = &user
	r.nextID++

	return nil
}

func (r *UserInMemory) Delete(ctx context.Context, id int64) error {
	r.mu.Lock()
	defer r.mu.Unlock()

	if _, exists := r.users[id]; !exists {
		return errors.New("product not found")
	}

	delete(r.users, id)
	return nil
}
