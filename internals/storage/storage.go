package storage

import "errors"

var (
	ErrUserNotFound = errors.New("user not found")
	ErrAppNotFound  = errors.New("application not found")
	ErrUserExists   = errors.New("user already exists")
)
