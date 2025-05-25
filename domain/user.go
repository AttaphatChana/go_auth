package domain

import "errors"

type User struct {
	ID       string
	Username string
	Password string // stored hashed
}

var (
	ErrUserNotFound = errors.New("user not found")
)
