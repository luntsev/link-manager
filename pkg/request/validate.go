package request

import "github.com/go-playground/validator"

func IsValid[T any](payload T) error {
	validate := validator.New()
	return validate.Struct(payload)
}
