package controller

import (
	"github.com/gofiber/fiber/v2"
	"pplace_backend/internal/model/dto/response"
	"pplace_backend/internal/service"
)

type PixelController struct {
	service *service.PixelService
}

func NewPixelController(service *service.PixelService) PixelController {
	return PixelController{
		service: service,
	}
}

func (c *PixelController) GetAllPixels(ctx *fiber.Ctx) error {
	pixels, err := c.service.GetAllPixels()
	if err != nil {
		errorDto := response.HttpErrorDto{
			StatusCode: err.StatusCode,
			Message:    err.Message,
			Errors:     err.Errors,
		}
		return ctx.Status(err.StatusCode).JSON(errorDto)
	}

	return ctx.JSON(pixels)
}
