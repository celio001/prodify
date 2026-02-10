package user_errors

import (
	"errors"

	"github.com/go-playground/validator/v10"
)

var (
	ErrUserNotFound       = errors.New("user not found")
	ErrUserCreationFailed = errors.New("failed to create user")
)

func CreateUserValidateError(err error) map[string]string {
	errors := make(map[string]string)

	if validationErrs, ok := err.(validator.ValidationErrors); ok {
		for _, fieldErr := range validationErrs {

			field := fieldErr.Field()
			tag := fieldErr.Tag()

			switch field {

			case "Name":
				switch tag {
				case "required":
					errors[field] = "name is required"
				case "min":
					errors[field] = "name must have at least 3 characters"
				case "max":
					errors[field] = "name must have at most 50 characters"
				}

			case "Email":
				switch tag {
				case "required":
					errors[field] = "email is required"
				case "email":
					errors[field] = "invalid email format"
				}

			case "Password":
				if tag == "required" {
					errors[field] = "password is required"
				}

			case "IsActive":
				if tag == "required" {
					errors[field] = "isActive is required"
				}
			}
		}
	}

	return errors
}


func UpdateUserValidateError(err error) map[string]string {
	errors := make(map[string]string)

	if validationErrs, ok := err.(validator.ValidationErrors); ok {
		for _, fieldErr := range validationErrs {

			field := fieldErr.Field()
			tag := fieldErr.Tag()

			switch field {

			case "Name":
				switch tag {
				case "min":
					errors[field] = "name must have at least 3 characters"
				case "max":
					errors[field] = "name must have at most 100 characters"
				}

			case "Email":
				if tag == "email" {
					errors[field] = "invalid email format"
				}
			}
		}
	}

	return errors
}
