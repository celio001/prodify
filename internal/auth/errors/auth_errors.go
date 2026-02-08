package auth_errors

import (
	"errors"
	"github.com/go-playground/validator/v10"
)

var (
	ErrMatchDataUser = errors.New("email or password incorrect")
)

func LoginValidateError(err error) map[string]string {
	errors := make(map[string]string)

	if validationErrs, ok := err.(validator.ValidationErrors); ok {
		for _, fieldErr := range validationErrs {
			field := fieldErr.Field()
			switch field{
			case "Email":
				errors[field] = "Email is required"
			case "Password":
				errors[field] = "Password is required"
			}
		}
	}
	return errors
}