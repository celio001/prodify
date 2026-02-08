package user_errors

import "errors"

var (
	ErrUserNotFound       = errors.New("user not found")
	ErrUserCreationFailed = errors.New("failed to create user")
)
