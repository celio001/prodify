package user_handler

import (
	"github.com/celio001/prodify/pkg/logger"
	pkg_request "github.com/celio001/prodify/pkg/request"
	uuidvalidator "github.com/celio001/prodify/pkg/uuid-validator"

	user_errors "github.com/celio001/prodify/internal/user/errors"
	user_service "github.com/celio001/prodify/internal/user/service"
	user_types "github.com/celio001/prodify/internal/user/type"

	"github.com/go-playground/validator/v10"
	"github.com/gofiber/fiber/v2"
	"go.uber.org/zap"
)

type UserHandler interface {
	GetUserByPublicIDHandler(c *fiber.Ctx) error
	UpdateUserHandler(c *fiber.Ctx) error
	DeleteUserHandler(c *fiber.Ctx) error
	CreateUserHandler(c *fiber.Ctx) error
}

type userHandler struct {
	userService user_service.UserService
}

func NewUserHandler(userService user_service.UserService) UserHandler {
	return &userHandler{userService: userService}
}

const maxBodySize = 1 << 20
var validate = validator.New()

// @Summary Get authenticated user profile
// @Description Returns the authenticated user profile using the user ID extracted from the access token
// @Tags user
// @Accept json
// @Produce json
// @Security BearerAuth
// @Success 200 {object} map[string]interface{} "User loaded successfully"
// @Failure 400 {object} map[string]string "Invalid user ID or user not found"
// @Failure 401 {object} map[string]string "User not authenticated"
// @Failure 500 {object} map[string]string "Internal server error"
// @Router /v1/user [get]
func (h *userHandler) GetUserByPublicIDHandler(c *fiber.Ctx) error {

	userID := c.Locals("user_id")
	if userID == nil {
		return c.Status(fiber.StatusUnauthorized).
			JSON(fiber.Map{"error": "user not authenticated"})
	}

	id, err := uuidvalidator.ValidateUuid(userID.(string))
	if err != nil {
		logger.Log.Error("invalid uuid", zap.Error(err))
		return c.Status(fiber.StatusBadRequest).
			JSON(fiber.Map{"error": "INVALID_USER_ID"})
	}

	user, err := h.userService.GetUserByPublicID(id)
	if err != nil {
		switch err {
		case user_errors.ErrUserNotFound:
			return c.Status(fiber.StatusBadRequest).
				JSON(fiber.Map{"error": "USER_NOT_FOUND"})
		default:
			return c.Status(fiber.StatusInternalServerError).
				JSON(fiber.Map{"error": "INTERNAL_ERROR"})
		}
	}

	return c.Status(fiber.StatusOK).
		JSON(fiber.Map{
			"message": "user loaded successfully",
			"data":    user,
		})
}

// @Summary Create user
// @Description Creates a new user account
// @Tags user
// @Accept json
// @Produce json
// @Param request body user_types.CreateUserRequest true "User creation payload"
// @Success 201 {object} map[string]string "User created successfully"
// @Failure 400 {object} map[string]interface{} "Invalid request body or validation error"
// @Failure 500 {object} map[string]string "Internal server error"
// @Router /v1/user [post]
func (h *userHandler) CreateUserHandler(c *fiber.Ctx) error {

	var req user_types.CreateUserRequest

	if err := pkg_request.LimitBodyJSON(c, maxBodySize, &req); err != nil {
		return c.Status(fiber.StatusBadRequest).
			JSON(fiber.Map{"error": err.Error()})
	}

	if err := validate.Struct(req); err != nil {
		logger.Log.Error("invalid create user payload", zap.Error(err))
		return c.Status(fiber.StatusBadRequest).
			JSON(fiber.Map{"error": user_errors.CreateUserValidateError(err)})
	}

	if err := h.userService.CreateUser(req); err != nil {
		return c.Status(fiber.StatusInternalServerError).
			JSON(fiber.Map{"error": "INTERNAL_ERROR"})
	}

	return c.Status(fiber.StatusCreated).
		JSON(fiber.Map{"message": "user created successfully"})
}

// @Summary Update authenticated user
// @Description Updates authenticated user profile using the user ID extracted from the access token
// @Tags user
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param request body user_types.UpdateUserRequest true "User update payload"
// @Success 200 {object} map[string]string "User updated successfully"
// @Failure 400 {object} map[string]interface{} "Invalid request body, invalid user ID or user not found"
// @Failure 401 {object} map[string]string "User not authenticated"
// @Failure 500 {object} map[string]string "Internal server error"
// @Router /v1/user [patch]
func (h *userHandler) UpdateUserHandler(c *fiber.Ctx) error {

	var req user_types.UpdateUserRequest

	userID := c.Locals("user_id")
	if userID == nil {
		return c.Status(fiber.StatusUnauthorized).
			JSON(fiber.Map{"error": "user not authenticated"})
	}

	id, err := uuidvalidator.ValidateUuid(userID.(string))
	if err != nil {
		logger.Log.Error("invalid uuid", zap.Error(err))
		return c.Status(fiber.StatusBadRequest).
			JSON(fiber.Map{"error": "INVALID_USER_ID"})
	}

	if err := pkg_request.LimitBodyJSON(c, maxBodySize, &req); err != nil {
		return c.Status(fiber.StatusBadRequest).
			JSON(fiber.Map{"error": err.Error()})
	}

	if err := validate.Struct(req); err != nil {
		logger.Log.Error("invalid update payload", zap.Error(err))
		return c.Status(fiber.StatusBadRequest).
			JSON(fiber.Map{"error": user_errors.UpdateUserValidateError(err)})
	}

	err = h.userService.UpdateUser(id, req)
	if err != nil {
		switch err {
		case user_errors.ErrUserNotFound:
			return c.Status(fiber.StatusBadRequest).
				JSON(fiber.Map{"error": "USER_NOT_FOUND"})
		default:
			return c.Status(fiber.StatusInternalServerError).
				JSON(fiber.Map{"error": "INTERNAL_ERROR"})
		}
	}

	return c.Status(fiber.StatusOK).
		JSON(fiber.Map{"message": "user updated successfully"})
}

// @Summary Delete authenticated user
// @Description Soft deletes the authenticated user using the user ID extracted from the access token
// @Tags user
// @Accept json
// @Produce json
// @Security BearerAuth
// @Success 200 {object} map[string]string "User successfully deleted"
// @Failure 400 {object} map[string]string "Invalid user ID or user not found"
// @Failure 401 {object} map[string]string "User not authenticated"
// @Failure 500 {object} map[string]string "Internal server error"
// @Router /v1/user [delete]
func (h *userHandler) DeleteUserHandler(c *fiber.Ctx) error {

	userID := c.Locals("user_id")
	if userID == nil {
		return c.Status(fiber.StatusUnauthorized).
			JSON(fiber.Map{"error": "user not authenticated"})
	}

	id, err := uuidvalidator.ValidateUuid(userID.(string))
	if err != nil {
		logger.Log.Error("invalid uuid", zap.Error(err))
		return c.Status(fiber.StatusBadRequest).
			JSON(fiber.Map{"error": "INVALID_USER_ID"})
	}

	err = h.userService.SoftDeleteUser(id)
	if err != nil {
		switch err {
		case user_errors.ErrUserNotFound:
			return c.Status(fiber.StatusBadRequest).
				JSON(fiber.Map{"error": "USER_NOT_FOUND"})
		default:
			return c.Status(fiber.StatusInternalServerError).
				JSON(fiber.Map{"error": "INTERNAL_ERROR"})
		}
	}

	return c.Status(fiber.StatusOK).
		JSON(fiber.Map{"message": "user successfully deleted"})
}
