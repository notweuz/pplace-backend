package validation

import "github.com/go-playground/validator/v10"

type ValidationError struct {
	Tag   string      `json:"tag"`
	Field string      `json:"field"`
	Value interface{} `json:"value"`
}

func ValidateDTO(dto interface{}) []ValidationError {
	validate := validator.New()
	err := validate.Struct(dto)
	if err != nil {
		if _, ok := err.(*validator.InvalidValidationError); ok {
			return nil
		}
		var errors []ValidationError
		for _, e := range err.(validator.ValidationErrors) {
			errors = append(errors, ValidationError{
				Tag:   e.Tag(),
				Field: e.Field(),
				Value: e.Value(),
			})
		}
		return errors
	}
	return nil
}
