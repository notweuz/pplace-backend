package transport

import (
	"pplace_backend/internal/handler"
	"pplace_backend/internal/service"

	"github.com/gofiber/fiber/v2"
	"github.com/rs/zerolog/log"
)

func SetupInfoRoutes(group fiber.Router, service *service.InfoService) {
	infoHandler := handler.NewInfoHandler(service)

	log.Info().Msg("Setting up info routes")
	infoGroup := group.Group("/info")

	infoGroup.Get("/", infoHandler.GetPixelSheetInfo)
}
