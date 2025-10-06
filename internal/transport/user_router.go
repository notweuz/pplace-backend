package transport

import (
	"pplace_backend/internal/layer/handler"
	"pplace_backend/internal/layer/middleware"
	"pplace_backend/internal/layer/service"

	"github.com/gofiber/fiber/v2"
)

func SetupUserRoutes(app *fiber.App, service *service.UserService) {
	userHandler := handler.NewUserHandler(service)
	authMiddleware := middleware.AuthMiddleware(service)

	app.Post("/api/users", userHandler.CreateUser)
	app.Get("/api/users/:id", userHandler.GetUserByID)
	app.Get("/api/users/username/:username", userHandler.GetUserByUsername)

	app.Get("/api/users/me", authMiddleware, userHandler.GetSelfInfo)
	app.Put("/api/users/me", authMiddleware, userHandler.UpdateUser)
}
