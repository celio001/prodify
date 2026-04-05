package auth_service

import (
	auth_errors "github.com/celio001/prodify/internal/auth/errors"
	auth_types "github.com/celio001/prodify/internal/auth/types"
	user_repository "github.com/celio001/prodify/internal/user/repository"
	user_types "github.com/celio001/prodify/internal/user/type"
	"github.com/google/uuid"
	"golang.org/x/crypto/bcrypt"
)

type authService struct {
	userRepo user_repository.UserRepository
}

type AuthService interface {
	Login(loginRequest auth_types.LoginRequest) (user_types.GetUserResponse, error)
	RegisterUser(user auth_types.CreateUserRequest) (*auth_types.CreateUserResponse, error)
	ResetPassword(userPublicID uuid.UUID, resetPasswordRequest auth_types.ResetPasswordRequest) error
}

func NewAuthService(userRepo user_repository.UserRepository) AuthService {
	return &authService{
		userRepo: userRepo,
	}
}

func (s *authService) Login(loginRequest auth_types.LoginRequest) (user_types.GetUserResponse, error) {
	user, err := s.userRepo.GetUserByEmail(loginRequest.Email)
	if err != nil {
		return user_types.GetUserResponse{}, err
	}

	if err := bcrypt.CompareHashAndPassword([]byte(user.PasswordHash), []byte(loginRequest.Password)); err != nil {
		return user_types.GetUserResponse{}, auth_errors.ErrMatchDataUser
	}

	return *user, nil
}

func (s *authService) ResetPassword(userPublicID uuid.UUID, resetPasswordRequest auth_types.ResetPasswordRequest) error {
	user, err := s.userRepo.GetUserByPublicID(userPublicID)
	if err != nil {
		return err
	}

	err = s.userRepo.UpdateUserPassword(user.ID, resetPasswordRequest)
	if err != nil {
		return err
	}

	return nil
}

func (s *authService) RegisterUser(user auth_types.CreateUserRequest) (*auth_types.CreateUserResponse, error) {

	u := user_types.CreateUserRequest{
		Name:     user.Name,
		Email:    user.Email,
		Password: user.Password,
	}

	userExists, err := s.userRepo.GetUserByEmail(user.Email)
	if userExists != nil {
		return &auth_types.CreateUserResponse{}, auth_errors.ErrUserAlreadyExists
	}else if err != nil {
		return &auth_types.CreateUserResponse{}, err
	}

	userRepo, err := s.userRepo.CreateUser(u)
	if err != nil {
		return &auth_types.CreateUserResponse{}, err
	}
	return &auth_types.CreateUserResponse{
		Id:           userRepo.Id,
		PublicID:     userRepo.PublicID.String(),
		Name:         userRepo.Name,
		Email:        userRepo.Email,
		PasswordHash: userRepo.PasswordHash,
		IsActive:     userRepo.IsActive,
		CreatedAt:    userRepo.CreatedAt,
		UpdatedAt:    userRepo.UpdatedAt,
	}, nil
}
