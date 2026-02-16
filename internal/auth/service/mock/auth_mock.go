package auth_mock

import (
	auth_types "github.com/celio001/prodify/internal/auth/types"
	user_types "github.com/celio001/prodify/internal/user/type"
	"github.com/stretchr/testify/mock"
	"github.com/google/uuid"
)

type MockAuthService struct {
	mock.Mock
}

func (m *MockAuthService) Login(req auth_types.LoginRequest) (user_types.GetUserResponse, error) {
	args := m.Called(req)

	return args.Get(0).(user_types.GetUserResponse), args.Error(1)
}

func (m *MockAuthService) ResetPassword(userPublicID uuid.UUID, resetPasswordRequest auth_types.ResetPasswordRequest) error {
	args := m.Called(userPublicID, resetPasswordRequest)
	return args.Error(0)
}