package user_service_mock

import (
	user_types "github.com/celio001/prodify/internal/user/type"

	"github.com/google/uuid"
	"github.com/stretchr/testify/mock"
)

type MockUserService struct {
	mock.Mock
}

func (m *MockUserService) GetUserByPublicID(publicID uuid.UUID) (*user_types.GetUserResponse, error) {
	args := m.Called(publicID)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*user_types.GetUserResponse), args.Error(1)
}

func (m *MockUserService) SoftDeleteUser(publicID uuid.UUID) error {
	args := m.Called(publicID)
	return args.Error(0)
}

func (m *MockUserService) UpdateUser(publicID uuid.UUID, params user_types.UpdateUserRequest) error {
	args := m.Called(publicID, params)
	return args.Error(0)
}

func (m *MockUserService) CreateUser(user user_types.CreateUserRequest) error {
	args := m.Called(user)
	return args.Error(0)
}