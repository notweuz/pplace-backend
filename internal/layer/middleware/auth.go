package middleware

import (
	"pplace_backend/internal/layer/service"

	"github.com/gofiber/fiber/v2"
)

func AuthMiddleware(userService *service.UserService) fiber.Handler {
	return func(c *fiber.Ctx) error {
		tokenString, err := userService.ExtractToken(c)
		if err != nil {
			return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
				"error": "Invalid authorization header",
			})
		}

		user, err := userService.ParseAndValidateToken(tokenString)
		if err != nil {
			return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
				"error":   "Invalid or expired token",
				"details": err.Error(),
			})
		}

		c.Locals("user", user)
		return c.Next()
	}
}
