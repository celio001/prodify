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

func LimitBodyJSON(ctx *fiber.Ctx, limit int64, dest any) error {

	body := ctx.Body()

	if int64(len(body)) > limit {
		return ErrBodyTooLarge
	}

	if err := json.Unmarshal(body, dest); err != nil {
		return ErrJSONDecode
	}

	return nil
}