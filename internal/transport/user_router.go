package transport

import (
	"pplace_backend/internal/layer/handler"
	"pplace_backend/internal/layer/middleware"
	"pplace_backend/internal/layer/service"

	"github.com/gofiber/fiber/v2"
	"github.com/rs/zerolog/log"
)

func SetupUserRoutes(app *fiber.App, service *service.UserService) {
	userHandler := handler.NewUserHandler(service)
	authMiddleware := middleware.AuthMiddleware(service)

	log.Info().Msg("Setting up user routes")

	app.Post("/api/users", userHandler.CreateUser)
	app.Get("/api/users/:id", userHandler.GetUserByID)
	app.Get("/api/users/username/:username", userHandler.GetUserByUsername)

	app.Get("/api/users/me", authMiddleware, userHandler.GetSelfInfo)
	app.Put("/api/users/me", authMiddleware, userHandler.UpdateUser)
}
