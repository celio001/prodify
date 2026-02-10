package user_repository_mock

import (
	user_types "github.com/celio001/prodify/internal/user/type"
	"github.com/google/uuid"
	"github.com/stretchr/testify/mock"
)

type MockUserRepository struct {
	mock.Mock
}

func (m *MockUserRepository) GetUserByPublicID(publicID uuid.UUID) (*user_types.GetUserResponse, error) {
	args := m.Called(publicID)

	if args.Get(0) == nil {
		return nil, args.Error(1)
	}

	return args.Get(0).(*user_types.GetUserResponse), args.Error(1)
}

func (m *MockUserRepository) GetUserByEmail(email string) (*user_types.GetUserResponse, error) {
	args := m.Called(email)

	if args.Get(0) == nil {
		return nil, args.Error(1)
	}

	return args.Get(0).(*user_types.GetUserResponse), args.Error(1)
}

func (m *MockUserRepository) CreateUser(user user_types.CreateUserRequest) error {
	args := m.Called(user)
	return args.Error(0)
}

func (m *MockUserRepository) SoftDeleteUser(userID int64) error {
	args := m.Called(userID)
	return args.Error(0)
}

func (m *MockUserRepository) UpdateUser(userID int64, params user_types.UpdateUserRequest) error {
	args := m.Called(userID, params)
	return args.Error(0)
}
