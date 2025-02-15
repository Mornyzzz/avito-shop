package errors

import "errors"

var (
	ErrUserNotFound       = errors.New("user not found")
	ErrInvalidUsername    = errors.New("invalid username")
	ErrInvalidPassword    = errors.New("invalid password")
	ErrInvalidCredentials = errors.New("invalid credentials")
	ErrNotFound           = errors.New("not found")
)
