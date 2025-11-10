package product

import (
	"github.com/celio001/prodify/product"
	"github.com/gofiber/fiber/v2"
)

const (
	HandlerPath = "/product"
)

func RegisterRouter(router fiber.Router, productRepository product.Repository) {

	handler := NewProductHandler(productRepository)
	router.Get("", handler.GetProduct)
}
