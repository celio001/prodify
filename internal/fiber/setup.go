package fiber

import (
	product_repo "github.com/celio001/prodify/product"
	"github.com/gofiber/fiber/v2"
)

type HttpServer struct {
	app               *fiber.App
	productRepository product_repo.Repository
}

func CreateServer(productRepository product_repo.Repository) HttpServer {
	app := fiber.New(fiber.Config{
		DisableStartupMessage: true,
	})

	httpServer := HttpServer{
		app: app,
		productRepository: productRepository,
	}

	return httpServer
}
