package user

import (
	"time"

	"github.com/google/uuid"
)

type User struct {
	Id           int64     `json:"id"`
	PublicId     uuid.UUID `json:"publicId"`
	Name         string    `json:"name"`
	Email        string    `json:"email"`
	PasswordHash string    `json:"password_hash"`
	Role         string    `json:"role"`
	IsActive     bool      `json:"isActive"`
	CreatedAt    time.Time `json:"created_at"`
	UpdatedAt    time.Time `json:"updated_at"`
	DeletedAt    time.Time `json:"deleted_at"`
}


