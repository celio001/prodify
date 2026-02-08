package auth_service

import (
	"errors"
	"testing"

	auth_errors "github.com/celio001/prodify/internal/auth/errors"
	auth_types "github.com/celio001/prodify/internal/auth/types"
	user_errors "github.com/celio001/prodify/internal/user/errors"
	user_mock "github.com/celio001/prodify/internal/user/repository/mock"
	user_types "github.com/celio001/prodify/internal/user/type"

	"github.com/stretchr/testify/assert"
	"golang.org/x/crypto/bcrypt"
)

func TestLogin(t *testing.T) {

	validPassword := "123456"

	hash, _ := bcrypt.GenerateFromPassword([]byte(validPassword), bcrypt.DefaultCost)

	tests := []struct {
		name        string
		mockReturn  *user_types.GetUserResponse
		mockError   error
		request     auth_types.LoginRequest
		expectError error
	}{
		{
			name: "success",
			mockReturn: &user_types.GetUserResponse{
				Email:        "test@mail.com",
				PasswordHash: string(hash),
			},
			mockError: nil,
			request: auth_types.LoginRequest{
				Email:    "test@mail.com",
				Password: validPassword,
			},
			expectError: nil,
		},
		{
			name:       "user not found",
			mockReturn: nil,
			mockError:  user_errors.ErrUserNotFound,
			request: auth_types.LoginRequest{
				Email:    "test@mail.com",
				Password: validPassword,
			},
			expectError: user_errors.ErrUserNotFound,
		},
		{
			name:       "repository error",
			mockReturn: nil,
			mockError:  errors.New("db error"),
			request: auth_types.LoginRequest{
				Email:    "test@mail.com",
				Password: validPassword,
			},
			expectError: errors.New("db error"),
		},
		{
			name: "invalid password",
			mockReturn: &user_types.GetUserResponse{
				Email:        "test@mail.com",
				PasswordHash: string(hash),
			},
			mockError: nil,
			request: auth_types.LoginRequest{
				Email:    "test@mail.com",
				Password: "wrong-password",
			},
			expectError: auth_errors.ErrMatchDataUser,
		},
	}

	for _, tt := range tests {

		t.Run(tt.name, func(t *testing.T) {

			mockRepo := new(user_mock.MockUserRepository)

			mockRepo.
				On("GetUserByEmail", tt.request.Email).
				Return(tt.mockReturn, tt.mockError)

			service := NewAuthService(mockRepo)

			result, err := service.Login(tt.request)

			if tt.expectError == nil {
				assert.NoError(t, err)
				assert.Equal(t, tt.request.Email, result.Email)
			} else {
				assert.Error(t, err)

				if errors.Is(tt.expectError, auth_errors.ErrMatchDataUser) ||
					errors.Is(tt.expectError, user_errors.ErrUserNotFound) {
					assert.ErrorIs(t, err, tt.expectError)
				}
			}

			mockRepo.AssertExpectations(t)
		})
	}
}
