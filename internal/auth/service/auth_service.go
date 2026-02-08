package auth_service

import (
	auth_types "github.com/celio001/prodify/internal/auth/types"
	auth_errors "github.com/celio001/prodify/internal/auth/errors"
	user_repository "github.com/celio001/prodify/internal/user/repository"
	user_types "github.com/celio001/prodify/internal/user/type"
	"golang.org/x/crypto/bcrypt"
)

type authService struct {
	userRepo user_repository.UserRepository
}

type AuthService interface {
	Login(loginRequest auth_types.LoginRequest) (user_types.GetUserResponse, error)
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
