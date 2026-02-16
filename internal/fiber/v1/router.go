package v1

import (
	auth_service "github.com/celio001/prodify/internal/auth/service"
	auth_handler "github.com/celio001/prodify/internal/fiber/v1/auth"
	product_handler "github.com/celio001/prodify/internal/fiber/v1/product"
	user_handler "github.com/celio001/prodify/internal/fiber/v1/user"
	user_service "github.com/celio001/prodify/internal/user/service"
	product_repo "github.com/celio001/prodify/product"
	"github.com/gofiber/fiber/v2"
)

const (
	HandlerPath = "/v1"
)

func RegisterRouter(router fiber.Router, productRepository product_repo.Repository, authSvc auth_service.AuthService, userSvc user_service.UserService) {
	productRouter := router.Group(product_handler.HandlerPath)
	authRouter := router.Group(auth_handler.HandlerPath)
	userRouter := router.Group(user_handler.HandlerPath)

	auth_handler.RegisterRouter(authRouter, authSvc)
	user_handler.RegisterRouter(userRouter, userSvc)
	
	product_handler.RegisterRouter(productRouter, productRepository)
	
}
