package pkg_request

import (
	"encoding/json"
	"errors"

	"github.com/gofiber/fiber/v2"
)

var (
	ErrBodyTooLarge = errors.New("request body too large")
	ErrJSONDecode   = errors.New("failed to decode json body")
)

func LimitBodyJSON(c *fiber.Ctx, limit int, dest any) error {
	body := c.Body()
	if body == nil || len(body) == 0 {
		return ErrJSONDecode
	}

	if len(body) > limit {
		return ErrBodyTooLarge
	}

	if err := json.Unmarshal(body, dest); err != nil {
		return ErrJSONDecode
	}

	return nil
}