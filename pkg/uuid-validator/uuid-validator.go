package uuidvalidator

import (
	"errors"

	"github.com/google/uuid"
)

var ErrInvalidUuid = errors.New("invalid user id")

func ValidateUuid(userID string) (id uuid.UUID, err error) {

	if id, err = uuid.Parse(userID); err != nil {
		return uuid.UUID{}, ErrInvalidUuid
	}

	return id, nil

}
