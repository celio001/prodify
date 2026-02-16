package auth_handler

import (
	auth_service "github.com/celio001/prodify/internal/auth/service"
	"github.com/gofiber/fiber/v2"
)

const (
	HandlerPath = "/auth"
)

func RegisterRouter(router fiber.Router, authService auth_service.AuthService) {

	handler := NewAuthHandler(authService)
	router.Post("/login", handler.AuthLoginHandler)
}
