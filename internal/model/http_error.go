package model

import (
	"fmt"

	"github.com/rs/zerolog/log"
)

type HttpError struct {
	StatusCode int
	Message    string
	Errors     []string
}

func NewHttpError(statusCode int, message string, errors ...string) *HttpError {
	err := &HttpError{
		StatusCode: statusCode,
		Message:    message,
		Errors:     errors,
	}

	log.Error().Err(err).Msg(message)
	return err
}

func (e *HttpError) Error() string {
	return fmt.Sprintf("%d - %s", e.StatusCode, e.Message)
}
