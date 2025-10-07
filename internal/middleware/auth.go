package middleware

import (
	"pplace_backend/internal/model/dto/response"
	"pplace_backend/internal/service"

	"github.com/gofiber/fiber/v2"
)

func AuthMiddleware(userService *service.UserService) fiber.Handler {
	return func(c *fiber.Ctx) error {
		tokenString, err := userService.ExtractToken(c)
		if err != nil {
			return response.NewHttpError(fiber.StatusUnauthorized, "Authorization token missing", []string{err.Error()})
		}

		user, err := userService.ParseAndValidateToken(tokenString)
		if err != nil {
			return response.NewHttpError(fiber.StatusUnauthorized, "Invalid or expired token", []string{err.Error()})
		}

		c.Locals("user", user)
		return c.Next()
	}
}
