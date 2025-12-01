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

	id := c.Query("id")

	prod, err := h.productRepository.FindByID(c.Context(), id)
	if err != nil {
		if err == product.ErrProductNotFound {
			return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": product.ErrProductNotFound.Error()})
		} else {
			return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": err.Error()})
		}
	}

	return c.Status(fiber.StatusOK).JSON(prod)
}
