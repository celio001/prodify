package auth_handler

import (
	auth_errors "github.com/celio001/prodify/internal/auth/errors"
	auth_service "github.com/celio001/prodify/internal/auth/service"
	auth_types "github.com/celio001/prodify/internal/auth/types"
	user_errors "github.com/celio001/prodify/internal/user/errors"
	pkg_jwt "github.com/celio001/prodify/pkg/jwt"
	"github.com/celio001/prodify/pkg/logger"
	pkg_request "github.com/celio001/prodify/pkg/request"
	"github.com/go-playground/validator/v10"
	"github.com/gofiber/fiber/v2"
	"go.uber.org/zap"
)

type authHandler struct {
	authService auth_service.AuthService
}

type AuthHandler interface {
	AuthLoginHandler(ctx *fiber.Ctx) error
}

func NewAuthHandler(authService auth_service.AuthService) *authHandler {
	return &authHandler{
		authService: authService,
	}
}

const maxBodySize = 1 << 20 // 1MB
var validate = validator.New()

func (h *authHandler) AuthLoginHandler(ctx *fiber.Ctx) error {
	var loginRequest auth_types.LoginRequest

	if err := pkg_request.LimitBodyJSON(ctx, maxBodySize, &loginRequest); err != nil {
		return err
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

	access, err := pkg_jwt.CreateAccessToken(user.PublicID)
	if err != nil {
		logger.Log.Error("failed to create access token", zap.String("error", err.Error()))
		return ctx.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "failed to create access token"})
	}

	return ctx.Status(fiber.StatusOK).JSON(fiber.Map{"access_token": access})
}
