package product

import (
	"github.com/celio001/prodify/product"
	"github.com/gofiber/fiber/v2"
)

type ProductHandler struct {
	productRepository product.Repository
}

func NewProductHandler(productRepository product.Repository) *ProductHandler {
	return &ProductHandler{
		productRepository: productRepository,
	}
}

func (h *ProductHandler) GetProduct(c *fiber.Ctx) error {
	return nil
}
