package transport

import (
	"pplace_backend/internal/handler"
	"pplace_backend/internal/service"

	"github.com/gofiber/fiber/v2"
	"github.com/rs/zerolog/log"
)

func SetupPixelRoutes(group fiber.Router, service *service.PixelService) {
	pixelHandler := handler.NewPixelHandler(service)

	log.Info().Msg("Setting up pixel routes")

	pixelsGroup := group.Group("/pixels")

	pixelsGroup.Post("/", pixelHandler.Create)
}
