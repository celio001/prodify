package v1

import (
	"github.com/celio001/prodify/internal/fiber/v1/product"
	product_repo "github.com/celio001/prodify/product"
	"github.com/gofiber/fiber/v2"
)

const (
	HandlerPath = "/v1"
)

func RegisterRouter(router fiber.Router, productRepository product_repo.Repository) {
	productRouter := router.Group(product.HandlerPath)
	product.RegisterRouter(productRouter, productRepository)
}
