package errors

import "errors"

var (
	ErrUserNotFound       = errors.New("user not found")
	ErrInvalidUsername    = errors.New("invalid username")
	ErrInvalidCredentials = errors.New("invalid credentials")
	ErrNotFound           = errors.New("not found")
	ErrMultiplyRows       = errors.New("multiple rows returned")
)
