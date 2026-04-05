package auth_handler

import (
	auth_service "github.com/celio001/prodify/internal/auth/service"
	"github.com/celio001/prodify/internal/fiber/middleware"
	"github.com/gofiber/fiber/v2"
)

const (
	HandlerPath = "/auth"
)

func RegisterRouter(router fiber.Router, authService auth_service.AuthService) {

	handler := NewAuthHandler(authService)
	router.Post("/login", handler.AuthLoginHandler)
	router.Post("/register", handler.RegisterUserHandler)
	router.Patch("/reset-password", middleware.AuthMiddleware(), handler.AuthResetPasswordHandler, )
}
