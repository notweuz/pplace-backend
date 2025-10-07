package middleware

import (
	"time"

	"github.com/gofiber/fiber/v2"
	"github.com/rs/zerolog/log"
)

func LoggingMiddleware() fiber.Handler {
	return func(c *fiber.Ctx) error {
		start := time.Now()

		log.Debug().Msgf("HTTP %s %s from %s", c.Method(), c.Path(), c.IP())

		err := c.Next()

		duration := time.Since(start)
		status := c.Response().StatusCode()

		log.Info().Msgf("HTTP %s %s completed with status %d in %v", c.Method(), c.Path(), status, duration)

		return err
	}
}
