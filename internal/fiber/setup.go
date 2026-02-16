package fiber

import (
	auth_service "github.com/celio001/prodify/internal/auth/service"
	user_service "github.com/celio001/prodify/internal/user/service"
	product_repo "github.com/celio001/prodify/product"
	"github.com/gofiber/fiber/v2"
)

type HttpServer struct {
	app               *fiber.App
	productRepository product_repo.Repository
	auth_service      auth_service.AuthService
	userService       user_service.UserService
}

func CreateServer(productRepository product_repo.Repository, authRepository auth_service.AuthService, userService user_service.UserService) HttpServer {
	app := fiber.New(fiber.Config{
		DisableStartupMessage: true,
	})

	httpServer := HttpServer{
		app:               app,
		productRepository: productRepository,
		auth_service:      authRepository,
		userService:       userService,
	}

	return httpServer
}
