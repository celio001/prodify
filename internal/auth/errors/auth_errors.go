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
			tag := fieldErr.Tag()

			switch field{
			case "Email":
				switch tag {
				case "required":
					errors[field] = "Email is required"
				case "email":
					errors[field] = "Invalid email format"
				}

			case "Password":
				switch tag {
				case "required":
					errors[field] = "Password is required"
				}
			}
		}
	}
	return errors
}

func RegisterValidateError(err error) map[string]string {
	errors := make(map[string]string)
	if validationErrs, ok := err.(validator.ValidationErrors); ok {
		for _, fieldErr := range validationErrs {
			field := fieldErr.Field()
			tag := fieldErr.Tag()
			switch field {
			case "Name":
				switch tag {
				case "required":
					errors[field] = "Name is required"
				case "min":
					errors[field] = "Name must be at least 3 characters"
				}
			case "Email":
				switch tag {
				case "required":
					errors[field] = "Email is required"
				case "email":
					errors[field] = "Invalid email format"
				}
			case "Password":
				switch tag {
				case "required":
					errors[field] = "Password is required"
				}
			}
		}
	}
	return errors
}

func ResetPasswordValidateError(err error) map[string]string {
	errors := make(map[string]string)

	if validationErrs, ok := err.(validator.ValidationErrors); ok {
		for _, fieldErr := range validationErrs {
			field := fieldErr.Field()
			switch field{
			case "OldPassword":
				errors[field] = "Old password is required"
			case "NewPassword":
				errors[field] = "New password is required"
			}
		}
	}
	return errors
}