package transport

import (
	"pplace_backend/internal/layer/handler"
	"pplace_backend/internal/layer/service"

	"github.com/gofiber/fiber/v2"
)

func SetupAuthRoutes(app *fiber.App, service *service.AuthService) {
	authHandler := handler.NewAuthHandler(service)

	app.Post("/api/auth/register", authHandler.Register)
	app.Post("/api/auth/login", authHandler.Login)
}
