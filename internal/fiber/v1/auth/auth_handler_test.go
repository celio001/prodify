package auth_handler

import (
	"errors"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	auth_errors "github.com/celio001/prodify/internal/auth/errors"
	auth_service "github.com/celio001/prodify/internal/auth/service"
	auth_mock "github.com/celio001/prodify/internal/auth/service/mock"
	user_errors "github.com/celio001/prodify/internal/user/errors"
	user_types "github.com/celio001/prodify/internal/user/type"
	"github.com/celio001/prodify/pkg/logger"

	"github.com/gofiber/fiber/v2"
	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

func setupTestApp(service *auth_mock.MockAuthService) *fiber.App {
	app := fiber.New()

	handler := &authHandler{
		authService: service,
	}

	app.Post("/login", handler.AuthLoginHandler)

	return app
}

func setupTestAppWithUser(service auth_service.AuthService, userID string) *fiber.App {
	app := fiber.New()

	handler := &authHandler{
		authService: service,
	}

	app.Patch("/reset-password", func(c *fiber.Ctx) error {
		c.Locals("user_id", userID)
		return handler.AuthResetPasswordHandler(c)
	})

	return app
}

func TestAuthLoginHandler_Success(t *testing.T) {

	logger.Init("dev")

	mockService := new(auth_mock.MockAuthService)

	userID := uuid.New().String()

	mockService.
		On("Login", mock.Anything).
		Return(user_types.GetUserResponse{
			PublicID: userID,
		}, nil)

	app := setupTestApp(mockService)

	body := `{
		"email":"test@mail.com",
		"password":"123456"
	}`

	req := httptest.NewRequest(http.MethodPost, "/login", strings.NewReader(body))
	req.Header.Set("Content-Type", "application/json")

	resp, err := app.Test(req)

	assert.NoError(t, err)
	assert.Equal(t, fiber.StatusOK, resp.StatusCode)

	mockService.AssertExpectations(t)
}

func TestAuthLoginHandler_InvalidPayload(t *testing.T) {

	logger.Init("dev")

	mockService := new(auth_mock.MockAuthService)
	app := setupTestApp(mockService)

	body := `{
		"email":"invalid-email"
	}`

	req := httptest.NewRequest(http.MethodPost, "/login", strings.NewReader(body))
	req.Header.Set("Content-Type", "application/json")

	resp, err := app.Test(req)

	assert.NoError(t, err)
	assert.Equal(t, fiber.StatusBadRequest, resp.StatusCode)

	mockService.AssertNotCalled(t, "Login")
}

func TestAuthLoginHandler_UserNotFound(t *testing.T) {

	logger.Init("dev")

	mockService := new(auth_mock.MockAuthService)

	mockService.
		On("Login", mock.Anything).
		Return(user_types.GetUserResponse{}, user_errors.ErrUserNotFound)

	app := setupTestApp(mockService)

	body := `{
		"email":"test@mail.com",
		"password":"123456"
	}`

	req := httptest.NewRequest(http.MethodPost, "/login", strings.NewReader(body))
	req.Header.Set("Content-Type", "application/json")

	resp, err := app.Test(req)

	assert.NoError(t, err)
	assert.Equal(t, fiber.StatusUnauthorized, resp.StatusCode)

	mockService.AssertExpectations(t)
}

func TestAuthLoginHandler_InvalidPassword(t *testing.T) {

	logger.Init("dev")

	mockService := new(auth_mock.MockAuthService)

	mockService.
		On("Login", mock.Anything).
		Return(user_types.GetUserResponse{}, auth_errors.ErrMatchDataUser)

	app := setupTestApp(mockService)

	body := `{
		"email":"test@mail.com",
		"password":"123456"
	}`

	req := httptest.NewRequest(http.MethodPost, "/login", strings.NewReader(body))
	req.Header.Set("Content-Type", "application/json")

	resp, err := app.Test(req)

	assert.NoError(t, err)
	assert.Equal(t, fiber.StatusUnauthorized, resp.StatusCode)

	mockService.AssertExpectations(t)
}

func TestAuthLoginHandler_InternalError(t *testing.T) {

	logger.Init("dev")

	mockService := new(auth_mock.MockAuthService)

	mockService.
		On("Login", mock.Anything).
		Return(user_types.GetUserResponse{}, errors.New("db error"))

	app := setupTestApp(mockService)

	body := `{
		"email":"test@mail.com",
		"password":"123456"
	}`

	req := httptest.NewRequest(http.MethodPost, "/login", strings.NewReader(body))
	req.Header.Set("Content-Type", "application/json")

	resp, err := app.Test(req)

	assert.NoError(t, err)
	assert.Equal(t, fiber.StatusInternalServerError, resp.StatusCode)

	mockService.AssertExpectations(t)
}

func TestAuthResetPasswordHandler_Success(t *testing.T) {
	logger.Init("dev")

	mockService := new(auth_mock.MockAuthService)
	userID := uuid.New().String()

	mockService.
		On("ResetPassword", mock.Anything, mock.Anything).
		Return(nil)

	app := setupTestAppWithUser(mockService, userID)

	body := `{
		"old_password_hash":"123456",
		"new_password":"654321"
	}`

	req := httptest.NewRequest(http.MethodPatch, "/reset-password", strings.NewReader(body))
	req.Header.Set("Content-Type", "application/json")

	resp, err := app.Test(req)

	assert.NoError(t, err)
	assert.Equal(t, fiber.StatusOK, resp.StatusCode)

	mockService.AssertExpectations(t)
}

func TestAuthResetPasswordHandler_UserNotAuthenticated(t *testing.T) {
	logger.Init("dev")

	mockService := new(auth_mock.MockAuthService)

	app := fiber.New()

	handler := &authHandler{
		authService: mockService,
	}

	app.Patch("/reset-password", func(c *fiber.Ctx) error {
		return handler.AuthResetPasswordHandler(c)
	})

	body := `{
		"old_password_hash":"123456",
		"new_password":"654321"
	}`

	req := httptest.NewRequest(http.MethodPatch, "/reset-password", strings.NewReader(body))
	req.Header.Set("Content-Type", "application/json")

	resp, err := app.Test(req)

	assert.NoError(t, err)
	assert.Equal(t, fiber.StatusUnauthorized, resp.StatusCode)

	mockService.AssertNotCalled(t, "ResetPassword")
}

func TestAuthResetPasswordHandler_InvalidUUID(t *testing.T) {
	logger.Init("dev")

	mockService := new(auth_mock.MockAuthService)
	app := setupTestAppWithUser(mockService, "invalid-uuid")

	body := `{
		"old_password":"123456",
		"new_password":"654321"
	}`

	req := httptest.NewRequest(http.MethodPatch, "/reset-password", strings.NewReader(body))
	req.Header.Set("Content-Type", "application/json")

	resp, err := app.Test(req)

	assert.NoError(t, err)
	assert.Equal(t, fiber.StatusBadRequest, resp.StatusCode)

	mockService.AssertNotCalled(t, "ResetPassword")
}

func TestAuthResetPasswordHandler_InvalidJSON(t *testing.T) {
	logger.Init("dev")

	mockService := new(auth_mock.MockAuthService)
	userID := uuid.New().String()

	app := setupTestAppWithUser(mockService, userID)

	body := `invalid-json`

	req := httptest.NewRequest(http.MethodPatch, "/reset-password", strings.NewReader(body))
	req.Header.Set("Content-Type", "application/json")

	resp, err := app.Test(req)

	assert.NoError(t, err)
	assert.Equal(t, fiber.StatusBadRequest, resp.StatusCode)

	mockService.AssertNotCalled(t, "ResetPassword")
}

func TestAuthResetPasswordHandler_ValidationError(t *testing.T) {
	logger.Init("dev")

	mockService := new(auth_mock.MockAuthService)
	userID := uuid.New().String()

	app := setupTestAppWithUser(mockService, userID)

	body := `{
		"old_password_hash":"123456"
	}`

	req := httptest.NewRequest(http.MethodPatch, "/reset-password", strings.NewReader(body))
	req.Header.Set("Content-Type", "application/json")

	resp, err := app.Test(req)

	assert.NoError(t, err)
	assert.Equal(t, fiber.StatusBadRequest, resp.StatusCode)

	mockService.AssertNotCalled(t, "ResetPassword")
}

func TestAuthResetPasswordHandler_SamePassword(t *testing.T) {
	logger.Init("dev")

	mockService := new(auth_mock.MockAuthService)
	userID := uuid.New().String()

	mockService.
		On("ResetPassword", mock.Anything, mock.Anything).
		Return(user_errors.ErrSamePassword)

	app := setupTestAppWithUser(mockService, userID)

	body := `{
		"old_password_hash":"123456",
		"new_password":"123456"
	}`

	req := httptest.NewRequest(http.MethodPatch, "/reset-password", strings.NewReader(body))
	req.Header.Set("Content-Type", "application/json")

	resp, err := app.Test(req)

	assert.NoError(t, err)
	assert.Equal(t, fiber.StatusUnauthorized, resp.StatusCode)

	mockService.AssertExpectations(t)
}

func TestAuthResetPasswordHandler_UserNotFound(t *testing.T) {
	logger.Init("dev")

	mockService := new(auth_mock.MockAuthService)
	userID := uuid.New().String()

	mockService.
		On("ResetPassword", mock.Anything, mock.Anything).
		Return(user_errors.ErrUserNotFound)

	app := setupTestAppWithUser(mockService, userID)

	body := `{
		"old_password_hash":"123456",
		"new_password":"654321"
	}`

	req := httptest.NewRequest(http.MethodPatch, "/reset-password", strings.NewReader(body))
	req.Header.Set("Content-Type", "application/json")

	resp, err := app.Test(req)

	assert.NoError(t, err)
	assert.Equal(t, fiber.StatusUnauthorized, resp.StatusCode)

	mockService.AssertExpectations(t)
}

func TestAuthResetPasswordHandler_InvalidOldPassword(t *testing.T) {
	logger.Init("dev")

	mockService := new(auth_mock.MockAuthService)
	userID := uuid.New().String()

	mockService.
		On("ResetPassword", mock.Anything, mock.Anything).
		Return(auth_errors.ErrMatchDataUser)

	app := setupTestAppWithUser(mockService, userID)

	body := `{
		"old_password_hash":"wrong-password",
		"new_password":"654321"
	}`

	req := httptest.NewRequest(http.MethodPatch, "/reset-password", strings.NewReader(body))
	req.Header.Set("Content-Type", "application/json")

	resp, err := app.Test(req)

	assert.NoError(t, err)
	assert.Equal(t, fiber.StatusUnauthorized, resp.StatusCode)

	mockService.AssertExpectations(t)
}

func TestAuthResetPasswordHandler_InternalError(t *testing.T) {
	logger.Init("dev")

	mockService := new(auth_mock.MockAuthService)
	userID := uuid.New().String()

	mockService.
		On("ResetPassword", mock.Anything, mock.Anything).
		Return(errors.New("database error"))

	app := setupTestAppWithUser(mockService, userID)

	body := `{
		"old_password_hash":"123456",
		"new_password":"654321"
	}`

	req := httptest.NewRequest(http.MethodPatch, "/reset-password", strings.NewReader(body))
	req.Header.Set("Content-Type", "application/json")

	resp, err := app.Test(req)

	assert.NoError(t, err)
	assert.Equal(t, fiber.StatusInternalServerError, resp.StatusCode)

	mockService.AssertExpectations(t)
}
