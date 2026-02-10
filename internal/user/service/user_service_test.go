package user_service

import (
	"errors"
	"testing"

	user_errors "github.com/celio001/prodify/internal/user/errors"
	user_mock "github.com/celio001/prodify/internal/user/repository/mock"
	user_types "github.com/celio001/prodify/internal/user/type"

	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
)

func TestSoftDeleteUser_Success(t *testing.T) {

	mockRepo := new(user_mock.MockUserRepository)
	service := NewUserService(mockRepo)

	publicID := uuid.New()

	user := &user_types.GetUserResponse{
		ID: 1,
	}

	mockRepo.
		On("GetUserByPublicID", publicID).
		Return(user, nil)

	mockRepo.
		On("SoftDeleteUser", int64(1)).
		Return(nil)

	err := service.SoftDeleteUser(publicID)

	assert.NoError(t, err)
	mockRepo.AssertExpectations(t)
}

func TestSoftDeleteUser_GetUserError(t *testing.T) {

	mockRepo := new(user_mock.MockUserRepository)
	service := NewUserService(mockRepo)

	publicID := uuid.New()
	mockRepo.
		On("GetUserByPublicID", publicID).
		Return(nil, user_errors.ErrUserNotFound)

	err := service.SoftDeleteUser(publicID)

	assert.Error(t, err)
	assert.Equal(t, user_errors.ErrUserNotFound, err)

	mockRepo.AssertExpectations(t)
}

func TestSoftDeleteUser_DeleteError(t *testing.T) {

	mockRepo := new(user_mock.MockUserRepository)
	service := NewUserService(mockRepo)

	publicID := uuid.New()
	deleteError := errors.New("delete error")

	user := &user_types.GetUserResponse{
		ID: 1,
	}

	mockRepo.
		On("GetUserByPublicID", publicID).
		Return(user, nil)

	mockRepo.
		On("SoftDeleteUser", int64(1)).
		Return(deleteError)

	err := service.SoftDeleteUser(publicID)

	assert.Error(t, err)
	assert.Equal(t, deleteError, err)

	mockRepo.AssertExpectations(t)
}

func TestUpdateUser_Success(t *testing.T) {

	mockRepo := new(user_mock.MockUserRepository)
	service := NewUserService(mockRepo)

	publicID := uuid.New()

	params := user_types.UpdateUserRequest{
		Name:  "Novo Nome",
		Email: "novo@email.com",
	}

	user := &user_types.GetUserResponse{
		ID: 1,
	}

	mockRepo.
		On("GetUserByPublicID", publicID).
		Return(user, nil)

	mockRepo.
		On("UpdateUser", int64(1), params).
		Return(nil)

	err := service.UpdateUser(publicID, params)

	assert.NoError(t, err)
	mockRepo.AssertExpectations(t)
}

func TestUpdateUser_GetUserError(t *testing.T) {

	mockRepo := new(user_mock.MockUserRepository)
	service := NewUserService(mockRepo)

	publicID := uuid.New()
	params := user_types.UpdateUserRequest{}

	expectedError := errors.New("get user error")

	mockRepo.
		On("GetUserByPublicID", publicID).
		Return(nil, expectedError)

	err := service.UpdateUser(publicID, params)

	assert.Error(t, err)
	assert.Equal(t, expectedError, err)

	mockRepo.AssertExpectations(t)
}

func TestUpdateUser_UpdateError(t *testing.T) {

	mockRepo := new(user_mock.MockUserRepository)
	service := NewUserService(mockRepo)

	publicID := uuid.New()

	params := user_types.UpdateUserRequest{
		Name: "Novo Nome",
	}

	updateError := errors.New("update error")

	user := &user_types.GetUserResponse{
		ID: 1,
	}

	mockRepo.
		On("GetUserByPublicID", publicID).
		Return(user, nil)

	mockRepo.
		On("UpdateUser", int64(1), params).
		Return(updateError)

	err := service.UpdateUser(publicID, params)

	assert.Error(t, err)
	assert.Equal(t, updateError, err)

	mockRepo.AssertExpectations(t)
}

func TestCreateUser_Success(t *testing.T) {

	mockRepo := new(user_mock.MockUserRepository)
	service := NewUserService(mockRepo)

	userRequest := user_types.CreateUserRequest{
		Name:     "Célio",
		Email:    "celio@email.com",
		Password: "123456",
		IsActive: true,
	}

	mockRepo.
		On("CreateUser", userRequest).
		Return(nil)

	err := service.CreateUser(userRequest)

	assert.NoError(t, err)
	mockRepo.AssertExpectations(t)
}

func TestCreateUser_RepositoryError(t *testing.T) {

	mockRepo := new(user_mock.MockUserRepository)
	service := NewUserService(mockRepo)

	userRequest := user_types.CreateUserRequest{
		Name:     "Célio",
		Email:    "celio@email.com",
		Password: "123456",
		IsActive: true,
	}

	expectedError := errors.New("database error")

	mockRepo.
		On("CreateUser", userRequest).
		Return(expectedError)

	err := service.CreateUser(userRequest)

	assert.Error(t, err)
	assert.Equal(t, expectedError, err)

	mockRepo.AssertExpectations(t)
}
