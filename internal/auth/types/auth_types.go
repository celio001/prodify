package auth_types

import "time"

type LoginRequest struct {
	Email    string `json:"email" validate:"required,email"`
	Password string `json:"password" validate:"required"`
}

type ResetPasswordRequest struct {
	OldPasswordHash string `json:"old_password_hash" validate:"required"`
	NewPassword     string `json:"new_password" validate:"required"`
}

type CreateUserRequest struct {
	Name     string `json:"name" validate:"required,min=3,max=50"`
	Email    string `json:"email" validate:"required,email"`
	Password string `json:"password_hash" validate:"required"`
}

type CreateUserResponse struct {
	Id           int64     `json:"id"`
	PublicID     string    `json:"publicId"`
	Name         string    `json:"name"`
	Email        string    `json:"email"`
	PasswordHash string    `json:"passwordHash"`
	IsActive     bool      `json:"isActive"`
	CreatedAt    time.Time `json:"createdAt"`
	UpdatedAt    time.Time `json:"updatedAt"`
}
