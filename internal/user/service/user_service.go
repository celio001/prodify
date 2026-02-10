package user_service

// const minEntropyBits = 60

// if err := password_validator.Validate(user.Password, minEntropyBits); err != nil{
// 	return err
// }

import (
	user_repository "github.com/celio001/prodify/internal/user/repository"
	user_types "github.com/celio001/prodify/internal/user/type"
	"github.com/google/uuid"
)

type userService struct {
	userRepo user_repository.UserRepository
}

type UserService interface {
	GetUserByPublicID(publicID uuid.UUID) (*user_types.GetUserResponse, error)
	SoftDeleteUser(publicID uuid.UUID) error
	UpdateUser(publicID uuid.UUID, params user_types.UpdateUserRequest) error
	CreateUser(user user_types.CreateUserRequest) error
}

func NewUserService(userRepo user_repository.UserRepository) UserService {
	return &userService{
		userRepo: userRepo,
	}
}

func (s *userService) GetUserByPublicID(publicID uuid.UUID) (*user_types.GetUserResponse, error) {
	user, err := s.userRepo.GetUserByPublicID(publicID)
	if err != nil {
		return nil, err
	}
	return user, nil
}

func (s *userService) SoftDeleteUser(publicID uuid.UUID) error {
	user, err := s.userRepo.GetUserByPublicID(publicID)
	if err != nil {
		return err
	}
	err = s.userRepo.SoftDeleteUser(user.ID)
	if err != nil {
		return err
	}
	return nil
}

func (s *userService) UpdateUser(publicID uuid.UUID, params user_types.UpdateUserRequest) error {
	user, err := s.userRepo.GetUserByPublicID(publicID)
	if err != nil {
		return err
	}

	err = s.userRepo.UpdateUser(user.ID, params)
	if err != nil {
		return err
	}

	return nil
}

func (s *userService) CreateUser(user user_types.CreateUserRequest) error {
	err := s.userRepo.CreateUser(user)
	if err != nil {
		return err
	}

	return nil
}