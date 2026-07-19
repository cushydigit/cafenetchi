package repository

import "errors"

var (
	ErrUserNotFound    = errors.New("user not found")
	ErrPhoneExists     = errors.New("phone already exists")
	ErrDuplicateRecord = errors.New("duplicate record")
)
