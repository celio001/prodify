package auth_handler

import (
	"errors"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	auth_errors "github.com/celio001/prodify/internal/auth/errors"
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
