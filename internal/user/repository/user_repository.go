package user_repository

import (
	"context"
	"database/sql"
	"fmt"

	user_errors "github.com/celio001/prodify/internal/user/errors"
	user_types "github.com/celio001/prodify/internal/user/type"
	"github.com/celio001/prodify/pkg/logger"
	"github.com/google/uuid"
	"go.uber.org/zap"
	"golang.org/x/crypto/bcrypt"
)

const (
	getUserByPublicIDQuery = `SELECT id, publicId, name, email, password_hash, role, isActive, created_at, updated_at 
	FROM users 
	WHERE publicId = $1
	AND deleted_at IS NULL`

	getUserByEmailQuery = `SELECT publicId, name, email, password_hash, role, isActive, created_at, updated_at 
	FROM users 
	WHERE email = $1
	AND deleted_at IS NULL`

	createUserQuery = `INSERT INTO users (name, email, password_hash, is_active) 
	VALUES ($1, $2, $3, $4)`

	softDeleteUserQuery = `PDATE users
	SET deleted_at = now(), updated_at = now(), is_active = false
	WHERE user_id = $1`

	updateUserQuery = `UPDATE users
	SET
		name = COALESCE($2, name),
		email = COALESCE($4, email),
		updated_at = now()
	WHERE user_id = $1;`
)

type userRepository struct {
	Db *sql.DB
}

type UserRepository interface {
	GetUserByPublicID(publicId uuid.UUID) (*user_types.GetUserResponse, error)
	GetUserByEmail(email string) (*user_types.GetUserResponse, error)
	CreateUser(user user_types.CreateUserRequest) error
	SoftDeleteUser(user_id int64) error
	UpdateUser(user_id int64, user_params user_types.UpdateUserRequest) error
}

func NewUserRepository(Db *sql.DB) UserRepository {
	return &userRepository{
		Db: Db,
	}
}

func (r *userRepository) GetUserByPublicID(publicId uuid.UUID) (*user_types.GetUserResponse, error) {
	ctx := context.Background()

	row := r.Db.QueryRowContext(ctx, getUserByPublicIDQuery, publicId)

	var user user_types.GetUserResponse
	err := row.Scan(&user.PublicID,
		&user.ID,
		&user.Name,
		&user.Email,
		&user.PasswordHash,
		&user.Role,
		&user.IsActive,
		&user.CreatedAt,
		&user.UpdatedAt)
	if err != nil {
		if err == sql.ErrNoRows {
			logger.Log.Error("user not found", zap.String("error", err.Error()))
			return nil, user_errors.ErrUserNotFound
		}
		logger.Log.Error("user get user by public id", zap.String("error", err.Error()))
		return nil, err
	}
	return &user, nil
}

func (r *userRepository) GetUserByEmail(email string) (*user_types.GetUserResponse, error) {
	ctx := context.Background()

	row := r.Db.QueryRowContext(ctx, getUserByEmailQuery, email)

	var user user_types.GetUserResponse
	err := row.Scan(&user.PublicID,
		&user.Name,
		&user.Email,
		&user.PasswordHash,
		&user.Role,
		&user.IsActive,
		&user.CreatedAt,
		&user.UpdatedAt)
	if err != nil {
		if err == sql.ErrNoRows {
			logger.Log.Error("user not found", zap.String("error", err.Error()))
			return nil, user_errors.ErrUserNotFound
		}
		logger.Log.Error("user get user by email", zap.String("error", err.Error()))
		return nil, err
	}
	return &user, nil
}

func (r *userRepository) CreateUser(user user_types.CreateUserRequest) error {
	ctx := context.Background()

	passwordEncrypted, err := bcrypt.GenerateFromPassword([]byte(user.Password), bcrypt.DefaultCost)
	if err != nil {
		logger.Log.Error("error encrypted password", zap.String("error", err.Error()))
		return fmt.Errorf("%w: %v", user_errors.ErrUserCreationFailed, err)
	}

	_, err = r.Db.ExecContext(ctx, createUserQuery, user.Name, user.Email, passwordEncrypted, user.IsActive)
	if err != nil {
		return err
	}
	return nil
}

func (r *userRepository) SoftDeleteUser(user_id int64) error {
	ctx := context.Background()
	_, err := r.Db.ExecContext(ctx, softDeleteUserQuery, user_id)

	if err != nil {
		logger.Log.Error("error for soft delete user", zap.String("error", err.Error()))
		return err
	}

	return nil
}

func (r *userRepository) UpdateUser(user_id int64, user_params user_types.UpdateUserRequest) error {
	ctx := context.Background()
	_, err := r.Db.ExecContext(ctx, updateUserQuery, user_id, user_params.Name, user_params.Email)

	if err != nil {
		logger.Log.Error("error for update user", zap.String("error", err.Error()))
		return err
	}

	return nil
}