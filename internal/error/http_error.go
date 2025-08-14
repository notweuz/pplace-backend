package error

import "fmt"

type HttpError struct {
	StatusCode int
	Message    string
	Errors     []string
}

func NewHttpError(statusCode int, message string, errors ...string) *HttpError {
	return &HttpError{
		StatusCode: statusCode,
		Message:    message,
		Errors:     errors,
	}
}

func (e *HttpError) Error() string {
	return fmt.Sprintf("%d - %s", e.StatusCode, e.Message)
}
