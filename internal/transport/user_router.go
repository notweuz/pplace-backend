package transport

import (
	"pplace_backend/internal/layer/handler"
	"pplace_backend/internal/layer/middleware"
	"pplace_backend/internal/layer/service"

	"github.com/gofiber/fiber/v2"
	"github.com/rs/zerolog/log"
)

func SetupUserRoutes(group fiber.Router, service *service.UserService) {
	userHandler := handler.NewUserHandler(service)
	authMiddleware := middleware.AuthMiddleware(service)

	log.Info().Msg("Setting up user routes")

	userGroup := group.Group("/users")
	userGroup.Post("/", userHandler.CreateUser)
	userGroup.Get("/:id", userHandler.GetUserByID)
	userGroup.Get("/username/:username", userHandler.GetUserByUsername)

	userGroup.Get("/me", authMiddleware, userHandler.GetSelfInfo)
	userGroup.Put("/me", authMiddleware, userHandler.UpdateUser)
}
