package response

import (
	"fmt"
)

type HttpError struct {
	StatusCode int      `json:"status_code"`
	Message    string   `json:"message"`
	Errors     []string `json:"errors"`
}

func NewHttpError(statusCode int, message string, errors []string) *HttpError {
	err := &HttpError{
		StatusCode: statusCode,
		Message:    message,
		Errors:     errors,
	}

	return err
}

func (e *HttpError) Error() string {
	return fmt.Sprintf("%d - %s", e.StatusCode, e.Message)
}
