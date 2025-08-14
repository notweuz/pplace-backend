package error

import "fmt"

type HttpError struct {
	StatusCode int
	Message    string
}

func NewHttpError(statusCode int, message string) *HttpError {
	return &HttpError{
		StatusCode: statusCode,
		Message:    message,
	}
}

func (e *HttpError) Error() string {
	return fmt.Sprintf("%d - %s", e.StatusCode, e.Message)
}
