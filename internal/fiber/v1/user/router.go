package user_handler

import (
	"github.com/celio001/prodify/internal/fiber/middleware"
	user_service "github.com/celio001/prodify/internal/user/service"
	"github.com/gofiber/fiber/v2"
)

const (
	HandlerPath = "/user"
)

func RegisterRouter(router fiber.Router, userService user_service.UserService) {

	userHandler := NewUserHandler(userService)
	router.Get("/", middleware.AuthMiddleware(), userHandler.GetUserByPublicIDHandler)
	router.Patch("/", middleware.AuthMiddleware(), userHandler.UpdateUserHandler)
	router.Delete("/", middleware.AuthMiddleware(), userHandler.DeleteUserHandler)
	router.Post("/", middleware.AuthMiddleware(), userHandler.CreateUserHandler)
}
