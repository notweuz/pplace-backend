package middleware

import (
	"errors"

	"pplace_backend/internal/model/dto/response"

	"github.com/gofiber/fiber/v2"
	"github.com/rs/zerolog/log"
)

func CustomErrorHandler() func(c *fiber.Ctx, err error) error {
	return func(c *fiber.Ctx, err error) error {
		log.Error().Err(err).Msg("request error")

		var he *response.HttpError
		if errors.As(err, &he) && he != nil {
			return c.Status(he.StatusCode).JSON(he)
		}

		var fe *fiber.Error
		if errors.As(err, &fe) && fe != nil {
			return c.Status(fe.Code).JSON(response.NewHttpError(fe.Code, fe.Message, nil))
		}

		generic := response.NewHttpError(
			fiber.StatusInternalServerError,
			"Internal Server Error",
			[]string{err.Error()},
		)
		return c.Status(fiber.StatusInternalServerError).JSON(generic)
	}
}
