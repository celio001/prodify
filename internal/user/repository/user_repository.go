package user_repository

import (
	"context"
	"database/sql"
	"fmt"

	auth_types "github.com/celio001/prodify/internal/auth/types"
	user_errors "github.com/celio001/prodify/internal/user/errors"
	user_types "github.com/celio001/prodify/internal/user/type"
	"github.com/celio001/prodify/pkg/logger"
	"github.com/google/uuid"
	"go.uber.org/zap"
	"golang.org/x/crypto/bcrypt"
)

const (
	getUserByPublicIDQuery = `SELECT id, public_id, name, email, password_hash, is_active, created_at, updated_at 
	FROM users 
	WHERE public_id = $1
	AND deleted_at IS NULL`

	getUserByEmailQuery = `SELECT public_id, name, email, password_hash, is_active, created_at, updated_at 
	FROM users 
	WHERE email = $1
	AND deleted_at IS NULL`

	createUserQuery = `INSERT INTO users (name, email, password_hash, is_active) 
	VALUES ($1, $2, $3, $4)`

	softDeleteUserQuery = `UPDATE users
	SET deleted_at = now(), updated_at = now(), is_active = false
	WHERE user_id = $1`

	updateUserQuery = `UPDATE users
	SET
		name = COALESCE($2, name),
		email = COALESCE($4, email),
		updated_at = now()
	WHERE user_id = $1;`

	updateUserPasswordQuery = `UPDATE users
	SET
	password_hash = $2,
	updated_at = now()
	WHERE user_id = $1;`
)

type userRepository struct {
	Db *sql.DB
}

type UserRepository interface {
	GetUserByPublicID(publicId uuid.UUID) (*user_types.GetUserResponse, error)
	GetUserByEmail(email string) (*user_types.GetUserResponse, error)
	CreateUser(user user_types.CreateUserRequest) (*user_types.CreateUserResponse, error)
	SoftDeleteUser(user_id int64) error
	UpdateUser(user_id int64, user_params user_types.UpdateUserRequest) error
	UpdateUserPassword(user_id int64, resetPasswordRequest auth_types.ResetPasswordRequest) error
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
	err := row.Scan(
		&user.ID,
		&user.PublicID,
		&user.Name,
		&user.Email,
		&user.PasswordHash,
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

func (r *userRepository) CreateUser(user user_types.CreateUserRequest) (*user_types.CreateUserResponse, error) {
	ctx := context.Background()

	passwordEncrypted, err := bcrypt.GenerateFromPassword([]byte(user.Password), bcrypt.DefaultCost)
	if err != nil {
		logger.Log.Error("error encrypted password", zap.String("error", err.Error()))
		return nil, fmt.Errorf("%w: %v", user_errors.ErrUserCreationFailed, err)
	}

	row, err := r.Db.QueryContext(ctx, createUserQuery, user.Name, user.Email, passwordEncrypted, true)
	if err != nil {
		return nil, err
	}
	defer row.Close()

	var u user_types.CreateUserResponse

	err = row.Scan(
		&u.Id,
		&u.PublicID,
		&u.Name,
		&u.Email,
		&u.PasswordHash,
		&u.IsActive,
		&u.CreatedAt,
		&u.UpdatedAt,
	)

	return &u, nil
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

func (r *userRepository) UpdateUserPassword(user_id int64, resetPasswordRequest auth_types.ResetPasswordRequest) error {
	ctx := context.Background()

	err := bcrypt.CompareHashAndPassword([]byte(resetPasswordRequest.OldPasswordHash), []byte(resetPasswordRequest.NewPassword))
	if err == nil {
		logger.Log.Error("new password cannot be the same as the old password")
		return user_errors.ErrSamePassword
	}

	passwordEncrypted, err := bcrypt.GenerateFromPassword([]byte(resetPasswordRequest.NewPassword), bcrypt.DefaultCost)
	if err != nil {
		logger.Log.Error("error encrypting password", zap.String("error", err.Error()))
		return fmt.Errorf("%w: %v", user_errors.ErrUserCreationFailed, err)
	}

	_, err = r.Db.ExecContext(ctx, updateUserPasswordQuery, user_id, passwordEncrypted)
	if err != nil {
		logger.Log.Error("error updating user password", zap.String("error", err.Error()))
		return err
	}
	return nil
}
