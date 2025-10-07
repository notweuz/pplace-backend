package middleware

import (
	"errors"
	"strings"

	"pplace_backend/internal/model"

	"github.com/gofiber/fiber/v2"
	"github.com/rs/zerolog/log"
)

func CustomErrorHandler() func(c *fiber.Ctx, err error) error {
	return func(c *fiber.Ctx, err error) error {
		log.Error().Err(err).Msg("request error")

		var he *model.HttpError
		if errors.As(err, &he) && he != nil {
			return c.Status(he.StatusCode).JSON(he)
		}

		var fe *fiber.Error
		if errors.As(err, &fe) && fe != nil {
			return c.Status(fe.Code).JSON(model.NewHttpError(fe.Code, fe.Message, nil))
		}

		msg := strings.ToLower(err.Error())
		switch {
		case strings.Contains(msg, "not found") || strings.Contains(msg, "not exists"):
			return c.Status(fiber.StatusNotFound).JSON(model.NewHttpError(fiber.StatusNotFound, err.Error(), []string{err.Error()}))
		case strings.Contains(msg, "already exists") || strings.Contains(msg, "already taken"):
			return c.Status(fiber.StatusConflict).JSON(model.NewHttpError(fiber.StatusConflict, err.Error(), []string{err.Error()}))
		case strings.Contains(msg, "cooldown"):
			return c.Status(fiber.StatusForbidden).JSON(model.NewHttpError(fiber.StatusForbidden, err.Error(), []string{err.Error()}))
		case strings.Contains(msg, "invalid password") || strings.Contains(msg, "unauthorized") || strings.Contains(msg, "unauthenticated"):
			return c.Status(fiber.StatusUnauthorized).JSON(model.NewHttpError(fiber.StatusUnauthorized, err.Error(), []string{err.Error()}))
		}

		generic := model.NewHttpError(fiber.StatusInternalServerError, "Internal Server Error", []string{err.Error()})
		return c.Status(fiber.StatusInternalServerError).JSON(generic)
	}
}
