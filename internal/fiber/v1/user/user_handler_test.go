package user_handler

import (
	"bytes"
	"encoding/json"
	"errors"
	"net/http"
	"net/http/httptest"
	"testing"

	user_errors "github.com/celio001/prodify/internal/user/errors"
	user_service_mock "github.com/celio001/prodify/internal/user/service/mock"
	user_types "github.com/celio001/prodify/internal/user/type"
	"github.com/celio001/prodify/pkg/logger"

	"github.com/gofiber/fiber/v2"
	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

func TestGetUserByPublicIDHandler_Success(t *testing.T) {

	logger.Init("dev")

	mockService := new(user_service_mock.MockUserService)
	handler := NewUserHandler(mockService)

	userID := uuid.New()

	mockResponse := &user_types.GetUserResponse{
		PublicID: userID.String(),
		Name:     "Celio",
		Email:    "celio@test.com",
		IsActive: true,
	}

	mockService.
		On("GetUserByPublicID", userID).
		Return(mockResponse, nil)

	app := fiber.New()

	app.Get("/v1/user", func(c *fiber.Ctx) error {
		c.Locals("user_id", userID.String())
		return handler.GetUserByPublicIDHandler(c)
	})

	req := httptest.NewRequest(http.MethodGet, "/v1/user", nil)
	resp, _ := app.Test(req)

	assert.Equal(t, fiber.StatusOK, resp.StatusCode)

	mockService.AssertExpectations(t)
}

func TestGetUserByPublicIDHandler_Unauthorized(t *testing.T) {

	logger.Init("dev")

	mockService := new(user_service_mock.MockUserService)
	handler := NewUserHandler(mockService)

	app := fiber.New()

	app.Get("/v1/user", handler.GetUserByPublicIDHandler)

	req := httptest.NewRequest(http.MethodGet, "/v1/user", nil)
	resp, _ := app.Test(req)

	assert.Equal(t, fiber.StatusUnauthorized, resp.StatusCode)
}

func TestGetUserByPublicIDHandler_InvalidUUID(t *testing.T) {

	logger.Init("dev")

	mockService := new(user_service_mock.MockUserService)
	handler := NewUserHandler(mockService)

	app := fiber.New()

	app.Get("/v1/user", func(c *fiber.Ctx) error {
		c.Locals("user_id", "invalid-uuid")
		return handler.GetUserByPublicIDHandler(c)
	})

	req := httptest.NewRequest(http.MethodGet, "/v1/user", nil)
	resp, _ := app.Test(req)

	assert.Equal(t, fiber.StatusBadRequest, resp.StatusCode)
}

func TestGetUserByPublicIDHandler_UserNotFound(t *testing.T) {

	logger.Init("dev")

	mockService := new(user_service_mock.MockUserService)
	handler := NewUserHandler(mockService)

	userID := uuid.New()

	mockService.
		On("GetUserByPublicID", userID).
		Return(nil, user_errors.ErrUserNotFound)

	app := fiber.New()

	app.Get("/v1/user", func(c *fiber.Ctx) error {
		c.Locals("user_id", userID.String())
		return handler.GetUserByPublicIDHandler(c)
	})

	req := httptest.NewRequest(http.MethodGet, "/v1/user", nil)
	resp, _ := app.Test(req)

	assert.Equal(t, fiber.StatusBadRequest, resp.StatusCode)

	mockService.AssertExpectations(t)
}

func TestGetUserByPublicIDHandler_InternalError(t *testing.T) {

	logger.Init("dev")

	mockService := new(user_service_mock.MockUserService)
	handler := NewUserHandler(mockService)

	userID := uuid.New()

	mockService.
		On("GetUserByPublicID", userID).
		Return(nil, errors.New("database error"))

	app := fiber.New()

	app.Get("/v1/user", func(c *fiber.Ctx) error {
		c.Locals("user_id", userID.String())
		return handler.GetUserByPublicIDHandler(c)
	})

	req := httptest.NewRequest(http.MethodGet, "/v1/user", nil)
	resp, _ := app.Test(req)

	assert.Equal(t, fiber.StatusInternalServerError, resp.StatusCode)

	mockService.AssertExpectations(t)
}

func TestCreateUserHandler_Success(t *testing.T) {

	logger.Init("dev")

	mockService := new(user_service_mock.MockUserService)
	handler := NewUserHandler(mockService)

	payload := user_types.CreateUserRequest{
		Name:     "Celio",
		Email:    "celio@test.com",
		Password: "123456",
		IsActive: true,
	}

	mockService.
		On("CreateUser", payload).
		Return(nil)

	body, _ := json.Marshal(payload)

	app := fiber.New()

	app.Post("/v1/user", handler.CreateUserHandler)

	req := httptest.NewRequest(http.MethodPost, "/v1/user", bytes.NewReader(body))
	req.Header.Set("Content-Type", "application/json")

	resp, _ := app.Test(req)

	assert.Equal(t, fiber.StatusCreated, resp.StatusCode)

	mockService.AssertExpectations(t)
}

func TestCreateUserHandler_InvalidJSON(t *testing.T) {

	logger.Init("dev")

	mockService := new(user_service_mock.MockUserService)
	handler := NewUserHandler(mockService)

	app := fiber.New()
	app.Post("/v1/user", handler.CreateUserHandler)

	req := httptest.NewRequest(
		http.MethodPost,
		"/v1/user",
		bytes.NewBufferString("{invalid json}"),
	)

	req.Header.Set("Content-Type", "application/json")

	resp, _ := app.Test(req)

	assert.Equal(t, fiber.StatusBadRequest, resp.StatusCode)
}

func TestCreateUserHandler_ValidationError(t *testing.T) {

	logger.Init("dev")

	mockService := new(user_service_mock.MockUserService)
	handler := NewUserHandler(mockService)

	payload := user_types.CreateUserRequest{
		Name:     "C", 
		Email:    "invalid-email",
		Password: "",
		IsActive: true,
	}

	body, _ := json.Marshal(payload)

	app := fiber.New()
	app.Post("/v1/user", handler.CreateUserHandler)

	req := httptest.NewRequest(http.MethodPost, "/v1/user", bytes.NewReader(body))
	req.Header.Set("Content-Type", "application/json")

	resp, _ := app.Test(req)

	assert.Equal(t, fiber.StatusBadRequest, resp.StatusCode)

	mockService.AssertNotCalled(t, "CreateUser", mock.Anything)
}

func TestCreateUserHandler_InternalError(t *testing.T) {

	logger.Init("dev")

	mockService := new(user_service_mock.MockUserService)
	handler := NewUserHandler(mockService)

	payload := user_types.CreateUserRequest{
		Name:     "Celio",
		Email:    "celio@test.com",
		Password: "123456",
		IsActive: true,
	}

	mockService.
		On("CreateUser", payload).
		Return(errors.New("db error"))

	body, _ := json.Marshal(payload)

	app := fiber.New()
	app.Post("/v1/user", handler.CreateUserHandler)

	req := httptest.NewRequest(http.MethodPost, "/v1/user", bytes.NewReader(body))
	req.Header.Set("Content-Type", "application/json")

	resp, _ := app.Test(req)

	assert.Equal(t, fiber.StatusInternalServerError, resp.StatusCode)

	mockService.AssertExpectations(t)
}

func TestUpdateUserHandler_Success(t *testing.T) {

	logger.Init("dev")

	mockService := new(user_service_mock.MockUserService)
	handler := NewUserHandler(mockService)

	userID := uuid.New()

	payload := user_types.UpdateUserRequest{
		Name:  "Novo Nome",
		Email: "novo@email.com",
	}

	mockService.
		On("UpdateUser", userID, payload).
		Return(nil)

	body, _ := json.Marshal(payload)

	app := fiber.New()

	app.Patch("/v1/user", func(c *fiber.Ctx) error {
		c.Locals("user_id", userID.String())
		return handler.UpdateUserHandler(c)
	})

	req := httptest.NewRequest(http.MethodPatch, "/v1/user", bytes.NewReader(body))
	req.Header.Set("Content-Type", "application/json")

	resp, _ := app.Test(req)

	assert.Equal(t, fiber.StatusOK, resp.StatusCode)

	mockService.AssertExpectations(t)
}

func TestUpdateUserHandler_Unauthorized(t *testing.T) {

	logger.Init("dev")

	mockService := new(user_service_mock.MockUserService)
	handler := NewUserHandler(mockService)

	app := fiber.New()
	app.Patch("/v1/user", handler.UpdateUserHandler)

	req := httptest.NewRequest(http.MethodPatch, "/v1/user", nil)

	resp, _ := app.Test(req)

	assert.Equal(t, fiber.StatusUnauthorized, resp.StatusCode)
}

func TestUpdateUserHandler_InvalidUUID(t *testing.T) {

	logger.Init("dev")

	mockService := new(user_service_mock.MockUserService)
	handler := NewUserHandler(mockService)

	app := fiber.New()

	app.Patch("/v1/user", func(c *fiber.Ctx) error {
		c.Locals("user_id", "invalid-uuid")
		return handler.UpdateUserHandler(c)
	})

	req := httptest.NewRequest(http.MethodPatch, "/v1/user", nil)

	resp, _ := app.Test(req)

	assert.Equal(t, fiber.StatusBadRequest, resp.StatusCode)
}

func TestUpdateUserHandler_InvalidJSON(t *testing.T) {

	logger.Init("dev")

	mockService := new(user_service_mock.MockUserService)
	handler := NewUserHandler(mockService)

	userID := uuid.New()

	app := fiber.New()

	app.Patch("/v1/user", func(c *fiber.Ctx) error {
		c.Locals("user_id", userID.String())
		return handler.UpdateUserHandler(c)
	})

	req := httptest.NewRequest(
		http.MethodPatch,
		"/v1/user",
		bytes.NewBufferString("{invalid json}"),
	)

	req.Header.Set("Content-Type", "application/json")

	resp, _ := app.Test(req)

	assert.Equal(t, fiber.StatusBadRequest, resp.StatusCode)
}

func TestUpdateUserHandler_ValidationError(t *testing.T) {

	logger.Init("dev")

	mockService := new(user_service_mock.MockUserService)
	handler := NewUserHandler(mockService)

	userID := uuid.New()

	payload := user_types.UpdateUserRequest{
		Name: "A", // inválido (min 3)
	}

	body, _ := json.Marshal(payload)

	app := fiber.New()

	app.Patch("/v1/user", func(c *fiber.Ctx) error {
		c.Locals("user_id", userID.String())
		return handler.UpdateUserHandler(c)
	})

	req := httptest.NewRequest(http.MethodPatch, "/v1/user", bytes.NewReader(body))
	req.Header.Set("Content-Type", "application/json")

	resp, _ := app.Test(req)

	assert.Equal(t, fiber.StatusBadRequest, resp.StatusCode)

	mockService.AssertNotCalled(t, "UpdateUser", mock.Anything, mock.Anything)
}

func TestUpdateUserHandler_UserNotFound(t *testing.T) {

	logger.Init("dev")

	mockService := new(user_service_mock.MockUserService)
	handler := NewUserHandler(mockService)

	userID := uuid.New()

	payload := user_types.UpdateUserRequest{
		Name: "Novo Nome",
	}

	mockService.
		On("UpdateUser", userID, payload).
		Return(user_errors.ErrUserNotFound)

	body, _ := json.Marshal(payload)

	app := fiber.New()

	app.Patch("/v1/user", func(c *fiber.Ctx) error {
		c.Locals("user_id", userID.String())
		return handler.UpdateUserHandler(c)
	})

	req := httptest.NewRequest(http.MethodPatch, "/v1/user", bytes.NewReader(body))
	req.Header.Set("Content-Type", "application/json")

	resp, _ := app.Test(req)

	assert.Equal(t, fiber.StatusBadRequest, resp.StatusCode)

	mockService.AssertExpectations(t)
}

func TestUpdateUserHandler_InternalError(t *testing.T) {

	logger.Init("dev")

	mockService := new(user_service_mock.MockUserService)
	handler := NewUserHandler(mockService)

	userID := uuid.New()

	payload := user_types.UpdateUserRequest{
		Name: "Novo Nome",
	}

	mockService.
		On("UpdateUser", userID, payload).
		Return(errors.New("db error"))

	body, _ := json.Marshal(payload)

	app := fiber.New()

	app.Patch("/v1/user", func(c *fiber.Ctx) error {
		c.Locals("user_id", userID.String())
		return handler.UpdateUserHandler(c)
	})

	req := httptest.NewRequest(http.MethodPatch, "/v1/user", bytes.NewReader(body))
	req.Header.Set("Content-Type", "application/json")

	resp, _ := app.Test(req)

	assert.Equal(t, fiber.StatusInternalServerError, resp.StatusCode)

	mockService.AssertExpectations(t)
}

func TestDeleteUserHandler_Success(t *testing.T) {

	logger.Init("dev")

	mockService := new(user_service_mock.MockUserService)
	handler := NewUserHandler(mockService)

	userID := uuid.New()

	mockService.
		On("SoftDeleteUser", userID).
		Return(nil)

	app := fiber.New()

	app.Delete("/v1/user", func(c *fiber.Ctx) error {
		c.Locals("user_id", userID.String())
		return handler.DeleteUserHandler(c)
	})

	req := httptest.NewRequest(http.MethodDelete, "/v1/user", nil)

	resp, _ := app.Test(req)

	assert.Equal(t, fiber.StatusOK, resp.StatusCode)

	mockService.AssertExpectations(t)
}

func TestDeleteUserHandler_Unauthorized(t *testing.T) {

	logger.Init("dev")

	mockService := new(user_service_mock.MockUserService)
	handler := NewUserHandler(mockService)

	app := fiber.New()
	app.Delete("/v1/user", handler.DeleteUserHandler)

	req := httptest.NewRequest(http.MethodDelete, "/v1/user", nil)

	resp, _ := app.Test(req)

	assert.Equal(t, fiber.StatusUnauthorized, resp.StatusCode)

	mockService.AssertNotCalled(t, "SoftDeleteUser", mock.Anything)
}

func TestDeleteUserHandler_InvalidUUID(t *testing.T) {

	logger.Init("dev")

	mockService := new(user_service_mock.MockUserService)
	handler := NewUserHandler(mockService)

	app := fiber.New()

	app.Delete("/v1/user", func(c *fiber.Ctx) error {
		c.Locals("user_id", "invalid-uuid")
		return handler.DeleteUserHandler(c)
	})

	req := httptest.NewRequest(http.MethodDelete, "/v1/user", nil)

	resp, _ := app.Test(req)

	assert.Equal(t, fiber.StatusBadRequest, resp.StatusCode)

	mockService.AssertNotCalled(t, "SoftDeleteUser", mock.Anything)
}

func TestDeleteUserHandler_UserNotFound(t *testing.T) {

	logger.Init("dev")

	mockService := new(user_service_mock.MockUserService)
	handler := NewUserHandler(mockService)

	userID := uuid.New()

	mockService.
		On("SoftDeleteUser", userID).
		Return(user_errors.ErrUserNotFound)

	app := fiber.New()

	app.Delete("/v1/user", func(c *fiber.Ctx) error {
		c.Locals("user_id", userID.String())
		return handler.DeleteUserHandler(c)
	})

	req := httptest.NewRequest(http.MethodDelete, "/v1/user", nil)

	resp, _ := app.Test(req)

	assert.Equal(t, fiber.StatusBadRequest, resp.StatusCode)

	mockService.AssertExpectations(t)
}

func TestDeleteUserHandler_InternalError(t *testing.T) {

	logger.Init("dev")

	mockService := new(user_service_mock.MockUserService)
	handler := NewUserHandler(mockService)

	userID := uuid.New()

	mockService.
		On("SoftDeleteUser", userID).
		Return(errors.New("database error"))

	app := fiber.New()

	app.Delete("/v1/user", func(c *fiber.Ctx) error {
		c.Locals("user_id", userID.String())
		return handler.DeleteUserHandler(c)
	})

	req := httptest.NewRequest(http.MethodDelete, "/v1/user", nil)

	resp, _ := app.Test(req)

	assert.Equal(t, fiber.StatusInternalServerError, resp.StatusCode)

	mockService.AssertExpectations(t)
}
