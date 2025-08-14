package validation

import (
	"regexp"

	"github.com/go-playground/validator/v10"
)

type ValidationError struct {
	Error string      `json:"error"`
	Field string      `json:"field"`
	Value interface{} `json:"value"`
}

func ValidateDTO(dto interface{}) []ValidationError {
	validate := validator.New()

	validate.RegisterValidation("username", func(fl validator.FieldLevel) bool {
		re := regexp.MustCompile(`^[a-zA-Z0-9_]+$`)
		return re.MatchString(fl.Field().String())
	})

	validate.RegisterValidation("password", func(fl validator.FieldLevel) bool {
		re := regexp.MustCompile(`^[a-zA-Z0-9!@#\$%\^&\*]+$`)
		return re.MatchString(fl.Field().String())
	})

	err := validate.Struct(dto)
	if err != nil {
		if _, ok := err.(*validator.InvalidValidationError); ok {
			return nil
		}
		var errors []ValidationError
		for _, e := range err.(validator.ValidationErrors) {
			errors = append(errors, ValidationError{
				Error: e.Error(),
				Field: e.Field(),
				Value: e.Value(),
			})
		}
		return errors
	}
	return nil
}
