package auth_types

type LoginRequest struct {
	Email    string `json:"email" validate:"required,email"`
	Password string `json:"password" validate:"required"`
}

type ResetPasswordRequest struct {
	OldPasswordHash string `json:"old_password_hash" validate:"required"`
	NewPassword string `json:"new_password" validate:"required"`
}