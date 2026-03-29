package auth_handler

import (
	auth_errors "github.com/celio001/prodify/internal/auth/errors"
	auth_service "github.com/celio001/prodify/internal/auth/service"
	auth_types "github.com/celio001/prodify/internal/auth/types"
	user_errors "github.com/celio001/prodify/internal/user/errors"
	pkg_jwt "github.com/celio001/prodify/pkg/jwt"
	"github.com/celio001/prodify/pkg/logger"
	pkg_request "github.com/celio001/prodify/pkg/request"
	uuidvalidator "github.com/celio001/prodify/pkg/uuid-validator"
	"github.com/go-playground/validator/v10"
	"github.com/gofiber/fiber/v2"
	"go.uber.org/zap"
)

type authHandler struct {
	authService auth_service.AuthService
}

type AuthHandler interface {
	AuthLoginHandler(ctx *fiber.Ctx) error
	RegisterUserHandler(ctx *fiber.Ctx) error
}

func NewAuthHandler(authService auth_service.AuthService) *authHandler {
	return &authHandler{
		authService: authService,
	}
}

const maxBodySize = 1 << 20 // 1MB
var validate = validator.New()

// @Summary User login
// @Description Authenticates a user using email and password and returns a JWT access token
// @Tags auth
// @Accept json
// @Produce json
// @Param request body auth_types.LoginRequest true "Login payload"
// @Success 200 {object} map[string]string "Access token generated successfully"
// @Failure 400 {object} map[string]interface{} "Invalid request body or validation error"
// @Failure 401 {object} map[string]string "Invalid credentials or user not found"
// @Failure 500 {object} map[string]string "Internal server error"
// @Router /v1/auth/login [post]
func (h *authHandler) AuthLoginHandler(ctx *fiber.Ctx) error {
	var loginRequest auth_types.LoginRequest

	if err := pkg_request.LimitBodyJSON(ctx, maxBodySize, &loginRequest); err != nil {
		return ctx.Status(fiber.StatusBadRequest).
			JSON(fiber.Map{"error": err.Error()})
	}

	if err := validate.Struct(loginRequest); err != nil {
		logger.Log.Error("invalid login payload", zap.String("error", err.Error()))
		return ctx.Status(fiber.StatusBadRequest).
			JSON(fiber.Map{"error": auth_errors.LoginValidateError(err)})
	}

	user, err := h.authService.Login(loginRequest)
	if err != nil {
		switch err {
		case auth_errors.ErrMatchDataUser:
			return ctx.Status(fiber.StatusUnauthorized).JSON(fiber.Map{"error": err.Error()})
		case user_errors.ErrUserNotFound:
			return ctx.Status(fiber.StatusUnauthorized).JSON(fiber.Map{"error": err.Error()})
		default:
			logger.Log.Error("failed to login user", zap.String("error", err.Error()))
			return ctx.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "failed to login user"})
		}

	}

	token, err := pkg_jwt.CreateAccessToken(user.PublicID)
	if err != nil {
		logger.Log.Error("failed to create access token", zap.String("error", err.Error()))
		return ctx.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "failed to create access token"})
	}

	refresh_token, err := pkg_jwt.CreateRefreshToken(user.PublicID)
	if err != nil {
		logger.Log.Error("failed to create refresh token", zap.String("error", err.Error()))
		return ctx.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "failed to create refresh token"})
	}

	return ctx.Status(fiber.StatusOK).JSON(fiber.Map{
		"access_token":  token,
		"refresh_token": refresh_token,
		"token_type":    "Bearer",
	})
}

func (h *authHandler) RegisterUserHandler(ctx *fiber.Ctx) error {
	var createUserRequest auth_types.CreateUserRequest

	if err := pkg_request.LimitBodyJSON(ctx, maxBodySize, &createUserRequest); err != nil {
		return ctx.Status(fiber.StatusBadRequest).
			JSON(fiber.Map{"error": err.Error()})
	}

	if err := validate.Struct(createUserRequest); err != nil {
		logger.Log.Error("invalid create user payload", zap.String("error", err.Error()))
		return ctx.Status(fiber.StatusBadRequest).
			JSON(fiber.Map{"error": auth_errors.RegisterValidateError(err)})
	}

	user, err := h.authService.RegisterUser(createUserRequest)
	if err != nil {
		switch err {
		case user_errors.ErrUserCreationFailed:
			return ctx.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": err.Error()})
		default:
			logger.Log.Error("failed to register user", zap.String("error", err.Error()))
			return ctx.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "failed to register user"})
		}
	}

	token, err := pkg_jwt.CreateAccessToken(user.PublicID)
	if err != nil {
		logger.Log.Error("failed to create access token", zap.String("error", err.Error()))
		return ctx.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "failed to create access token"})
	}

	refresh_token, err := pkg_jwt.CreateRefreshToken(user.PublicID)
	if err != nil {
		logger.Log.Error("failed to create refresh token", zap.String("error", err.Error()))
		return ctx.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "failed to create refresh token"})
	}

	return ctx.Status(fiber.StatusOK).JSON(fiber.Map{
		"access_token":  token,
		"refresh_token": refresh_token,
		"token_type":    "Bearer",
	})
}

// @Summary Reset user password
// @Description Allows an authenticated user to reset their password
// @Tags auth
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param request body auth_types.ResetPasswordRequest true "Reset password payload"
// @Success 200 {object} map[string]string "Password reset successfully"
// @Failure 400 {object} map[string]interface{} "Invalid request body, validation error or invalid user ID"
// @Failure 401 {object} map[string]string "User not authenticated or business rule violation"
// @Failure 500 {object} map[string]string "Internal server error"
// @Router /v1/auth/reset-password [patch]
func (h *authHandler) AuthResetPasswordHandler(ctx *fiber.Ctx) error {
	var resetPasswordRequest auth_types.ResetPasswordRequest

	userID := ctx.Locals("user_id")
	if userID == nil {
		return ctx.Status(fiber.StatusUnauthorized).
			JSON(fiber.Map{"error": "user not authenticated"})
	}

	userIDParsed, err := uuidvalidator.ValidateUuid(userID.(string))
	if err != nil {
		logger.Log.Error("invalid uuid", zap.Error(err))
		return ctx.Status(fiber.StatusBadRequest).
			JSON(fiber.Map{"error": "INVALID_USER_ID"})
	}

	if err := pkg_request.LimitBodyJSON(ctx, maxBodySize, &resetPasswordRequest); err != nil {
		return ctx.Status(fiber.StatusBadRequest).
			JSON(fiber.Map{"error": err.Error()})
	}

	if err := validate.Struct(resetPasswordRequest); err != nil {
		logger.Log.Error("invalid reset password payload", zap.String("error", err.Error()))
		return ctx.Status(fiber.StatusBadRequest).
			JSON(fiber.Map{"error": auth_errors.ResetPasswordValidateError(err)})
	}

	if err := h.authService.ResetPassword(userIDParsed, resetPasswordRequest); err != nil {
		switch err {
		case user_errors.ErrSamePassword:
			return ctx.Status(fiber.StatusUnauthorized).JSON(fiber.Map{"error": err.Error()})
		case user_errors.ErrUserCreationFailed:
			return ctx.Status(fiber.StatusUnauthorized).JSON(fiber.Map{"error": err.Error()})
		case user_errors.ErrUserNotFound:
			return ctx.Status(fiber.StatusUnauthorized).JSON(fiber.Map{"error": err.Error()})
		case auth_errors.ErrMatchDataUser:
			return ctx.Status(fiber.StatusUnauthorized).JSON(fiber.Map{"error": err.Error()})
		default:
			logger.Log.Error("failed to reset password", zap.String("error", err.Error()))
			return ctx.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "failed to reset password"})
		}
	}

	return ctx.Status(fiber.StatusOK).JSON(fiber.Map{"message": "password reset successfully"})
}
