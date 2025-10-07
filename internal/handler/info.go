package handler

import (
	"pplace_backend/internal/service"

	"github.com/gofiber/fiber/v2"
)

type InfoHandler struct {
	service *service.InfoService
}

func NewInfoHandler(service *service.InfoService) *InfoHandler {
	return &InfoHandler{service: service}
}

func (h *InfoHandler) GetPixelSheetInfo(ctx *fiber.Ctx) error {
	info := h.service.GetPixelSheetInfo()
	return ctx.JSON(info)
}
