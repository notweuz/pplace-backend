package controller

import (
	"github.com/gofiber/fiber/v2"
	"pplace_backend/internal/service"
)

type InfoController struct {
	service *service.InfoService
}

func NewInfoController(service *service.InfoService) InfoController {
	return InfoController{service: service}
}

func (ic *InfoController) GetPixelSheetInfo(ctx *fiber.Ctx) error {
	info := ic.service.GetPixelSheetInfo()
	return ctx.JSON(info)
}
