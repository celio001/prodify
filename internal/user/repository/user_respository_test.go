package user_repository

import (
	"database/sql"
	"fmt"
	"regexp"
	"testing"
	"time"

	"github.com/DATA-DOG/go-sqlmock"
	user_errors "github.com/celio001/prodify/internal/user/errors"
	user_types "github.com/celio001/prodify/internal/user/type"
	"github.com/celio001/prodify/pkg/logger"
	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
)

func TestGetUserByPublicID(t *testing.T) {
	logger.Init("dev")

	publicID := uuid.New()
	now := time.Now()

	tests := []struct {
		name        string
		mockRows    *sqlmock.Rows
		mockError   error
		expectError error
	}{
		{
			name: "success",
			mockRows: sqlmock.NewRows([]string{
				"publicId",
				"name",
				"email",
				"password_hash",
				"role",
				"isActive",
				"created_at",
				"updated_at",
			}).AddRow(
				publicID,
				"Célio",
				"celio@email.com",
				"hash",
				"user",
				true,
				now,
				now,
			),
			mockError:   nil,
			expectError: nil,
		},
		{
			name:        "user not found",
			mockRows:    nil,
			mockError:   sql.ErrNoRows,
			expectError: user_errors.ErrUserNotFound,
		},
		{
			name:        "database scan error",
			mockRows:    nil,
			mockError:   fmt.Errorf("scan error"),
			expectError: fmt.Errorf("scan error"),
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {

			db, mock, err := sqlmock.New()
			assert.NoError(t, err)
			defer db.Close()

			repo := NewUserRepository(db)

			expect := mock.ExpectQuery(regexp.QuoteMeta(getUserByPublicIDQuery)).
				WithArgs(publicID)

			if tt.mockError != nil {
				expect.WillReturnError(tt.mockError)
			} else {
				expect.WillReturnRows(tt.mockRows)
			}

			user, err := repo.GetUserByPublicID(publicID)

			if tt.expectError == nil {
				assert.NoError(t, err)
				assert.NotNil(t, user)
				assert.Equal(t, "Célio", user.Name)
			} else if tt.expectError == user_errors.ErrUserNotFound {
				assert.Nil(t, user)
				assert.ErrorIs(t, err, user_errors.ErrUserNotFound)
			} else {
				assert.Nil(t, user)
				assert.Error(t, err)
			}

			assert.NoError(t, mock.ExpectationsWereMet())
		})
	}
}

func TestGetUserByEmail(t *testing.T) {
	logger.Init("dev")

	email := "celio@email.com"
	now := time.Now()

	tests := []struct {
		name        string
		mockRows    *sqlmock.Rows
		mockError   error
		expectError error
	}{
		{
			name: "success",
			mockRows: sqlmock.NewRows([]string{
				"publicId",
				"name",
				"email",
				"password_hash",
				"role",
				"isActive",
				"created_at",
				"updated_at",
			}).AddRow(
				uuid.New(),
				"Célio",
				email,
				"hash",
				"user",
				true,
				now,
				now,
			),
			mockError:   nil,
			expectError: nil,
		},
		{
			name:        "user not found",
			mockRows:    nil,
			mockError:   sql.ErrNoRows,
			expectError: user_errors.ErrUserNotFound,
		},
		{
			name:        "database scan error",
			mockRows:    nil,
			mockError:   fmt.Errorf("scan error"),
			expectError: fmt.Errorf("scan error"),
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {

			db, mock, err := sqlmock.New()
			assert.NoError(t, err)
			defer db.Close()

			repo := NewUserRepository(db)

			expect := mock.ExpectQuery(regexp.QuoteMeta(getUserByEmailQuery)).
				WithArgs(email)

			if tt.mockError != nil {
				expect.WillReturnError(tt.mockError)
			} else {
				expect.WillReturnRows(tt.mockRows)
			}

			user, err := repo.GetUserByEmail(email)

			if tt.expectError == nil {
				assert.NoError(t, err)
				assert.NotNil(t, user)
				assert.Equal(t, email, user.Email)
				assert.Equal(t, "Célio", user.Name)
			} else if tt.expectError == user_errors.ErrUserNotFound {
				assert.Nil(t, user)
				assert.ErrorIs(t, err, user_errors.ErrUserNotFound)
			} else {
				assert.Nil(t, user)
				assert.Error(t, err)
			}

			assert.NoError(t, mock.ExpectationsWereMet())
		})
	}
}

func TestCreateUser(t *testing.T) {
	logger.Init("dev")

	tests := []struct {
		name          string
		user          user_types.CreateUserRequest
		mockExecError error
		expectError   bool
		expectWrapErr bool
	}{
		{
			name: "success",
			user: user_types.CreateUserRequest{
				Name:     "Célio",
				Email:    "celio@email.com",
				Password: "123456",
				IsActive: true,
			},
			mockExecError: nil,
			expectError:   false,
		},
		{
			name: "database error",
			user: user_types.CreateUserRequest{
				Name:     "Célio",
				Email:    "celio@email.com",
				Password: "123456",
				IsActive: true,
			},
			mockExecError: fmt.Errorf("db error"),
			expectError:   true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {

			db, mock, err := sqlmock.New()
			assert.NoError(t, err)
			defer db.Close()

			repo := NewUserRepository(db)

			if tt.mockExecError != nil {
				mock.ExpectExec(regexp.QuoteMeta(createUserQuery)).
					WithArgs(
						tt.user.Name,
						tt.user.Email,
						sqlmock.AnyArg(),
						tt.user.IsActive,
					).
					WillReturnError(tt.mockExecError)
			} else {
				mock.ExpectExec(regexp.QuoteMeta(createUserQuery)).
					WithArgs(
						tt.user.Name,
						tt.user.Email,
						sqlmock.AnyArg(),
						tt.user.IsActive,
					).
					WillReturnResult(sqlmock.NewResult(1, 1))
			}

			err = repo.CreateUser(tt.user)

			if tt.expectError {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
			}

			assert.NoError(t, mock.ExpectationsWereMet())
		})
	}
}

func TestSoftDeleteUser(t *testing.T) {
	logger.Init("dev")

	tests := []struct {
		name        string
		mockError   error
		expectError error
	}{
		{
			name:        "success",
			mockError:   nil,
			expectError: nil,
		},
		{
			name:        "database error",
			mockError:   fmt.Errorf("db error"),
			expectError: fmt.Errorf("db error"),
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {

			db, mock, err := sqlmock.New()
			assert.NoError(t, err)
			defer db.Close()

			repo := NewUserRepository(db)

			userID := int64(1)

			expect := mock.ExpectExec(regexp.QuoteMeta(softDeleteUserQuery)).
				WithArgs(userID)

			if tt.mockError != nil {
				expect.WillReturnError(tt.mockError)
			} else {
				expect.WillReturnResult(sqlmock.NewResult(0, 1))
			}

			err = repo.SoftDeleteUser(userID)

			if tt.expectError == nil {
				assert.NoError(t, err)
			} else if tt.expectError == user_errors.ErrUserNotFound {
				assert.ErrorIs(t, err, user_errors.ErrUserNotFound)
			} else {
				assert.Error(t, err)
			}

			assert.NoError(t, mock.ExpectationsWereMet())
		})
	}
}

func TestUpdateUser(t *testing.T) {
	logger.Init("dev")

	tests := []struct {
		name        string
		mockError   error
		expectError bool
	}{
		{
			name:        "success",
			mockError:   nil,
			expectError: false,
		},
		{
			name:        "database error",
			mockError:   fmt.Errorf("db error"),
			expectError: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {

			db, mock, err := sqlmock.New()
			assert.NoError(t, err)
			defer db.Close()

			repo := NewUserRepository(db)

			userID := int64(1)
			params := user_types.UpdateUserRequest{
				Name:  "John Doe",
				Email: "john@email.com",
			}

			expect := mock.ExpectExec(regexp.QuoteMeta(updateUserQuery)).
				WithArgs(userID, params.Name, params.Email)

			if tt.mockError != nil {
				expect.WillReturnError(tt.mockError)
			} else {
				expect.WillReturnResult(sqlmock.NewResult(0, 1))
			}

			err = repo.UpdateUser(userID, params)

			if tt.expectError {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
			}

			assert.NoError(t, mock.ExpectationsWereMet())
		})
	}
}
