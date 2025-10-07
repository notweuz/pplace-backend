package transport

import (
	"pplace_backend/internal/handler"
	"pplace_backend/internal/middleware" // Добавить импорт
	"pplace_backend/internal/service"

	"github.com/gofiber/fiber/v2"
	"github.com/rs/zerolog/log"
)

func SetupPixelRoutes(group fiber.Router, service *service.PixelService, userService *service.UserService) {
	pixelHandler := handler.NewPixelHandler(service)
	authMiddleware := middleware.AuthMiddleware(userService)

	log.Info().Msg("Setting up pixel routes")
	pixelsGroup := group.Group("/pixels")

	pixelsGroup.Post("/", authMiddleware, pixelHandler.Create)
	pixelsGroup.Patch("/:id", authMiddleware, pixelHandler.Update)
}
