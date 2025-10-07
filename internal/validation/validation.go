package validation

import (
	"regexp"

	"github.com/go-playground/validator/v10"
)

type Error struct {
	Error string      `json:"error"`
	Field string      `json:"field"`
	Value interface{} `json:"value"`
}

func ValidateDTO(dto interface{}) []Error {
	validate := validator.New()

	err := validate.RegisterValidation("username", func(fl validator.FieldLevel) bool {
		re := regexp.MustCompile(`^[a-zA-Z0-9_]+$`)
		return re.MatchString(fl.Field().String())
	})
	if err != nil {
		return nil
	}

	err = validate.RegisterValidation("password", func(fl validator.FieldLevel) bool {
		re := regexp.MustCompile(`^[a-zA-Z0-9!@#\$%\^&\*]+$`)
		return re.MatchString(fl.Field().String())
	})
	if err != nil {
		return nil
	}

	err = validate.RegisterValidation("color", func(fl validator.FieldLevel) bool {
		re := regexp.MustCompile(`^#([A-Fa-f0-9]{6})$`)
		return re.MatchString(fl.Field().String())
	})
	if err != nil {
		return nil
	}

	err = validate.Struct(dto)
	if err != nil {
		if _, ok := err.(*validator.InvalidValidationError); ok {
			return nil
		}
		var errors []Error
		for _, e := range err.(validator.ValidationErrors) {
			errors = append(errors, Error{
				Error: e.Error(),
				Field: e.Field(),
				Value: e.Value(),
			})
		}
		return errors
	}
	return nil
}
