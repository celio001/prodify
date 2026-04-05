package user_types

import (
	"time"

	"github.com/google/uuid"
)

type CreateUserRequest struct {
	Name     string `json:"name" validate:"required,min=3,max=50"`
	Email    string `json:"email" validate:"required,email"`
	Password string `json:"password" validate:"required"`
}

type CreateUserResponse struct {
	Id       int64     `json:"id"`
	PublicID uuid.UUID `json:"publicId"`
	Name     string    `json:"name"`
	Email    string    `json:"email"`
	PasswordHash string    `json:"passwordHash"`
	IsActive bool      `json:"isActive"`
	CreatedAt time.Time `json:"createdAt"`
	UpdatedAt time.Time `json:"updatedAt"`
}

type UpdateUserRequest struct {
	Name  string `json:"name,omitempty" validate:"omitempty,min=3,max=100"`
	Email string `json:"email,omitempty" validate:"omitempty,email"`
}
type GetUserResponse struct {
	ID           int64     `json:"id"`
	PublicID     string    `json:"publicId"`
	Name         string    `json:"name"`
	Email        string    `json:"email"`
	PasswordHash string    `json:"passwordHash"`
	IsActive     bool      `json:"isActive"`
	CreatedAt    time.Time `json:"createdAt"`
	UpdatedAt    time.Time `json:"updatedAt"`
}
