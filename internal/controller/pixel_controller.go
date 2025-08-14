package controller

import (
	"github.com/gofiber/fiber/v2"
	"pplace_backend/internal/model/dto/request"
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

func (c *PixelController) PlacePixel(ctx *fiber.Ctx) error {
	var data request.PixelPlaceDto

	if err := ctx.BodyParser(&data); err != nil {
		return err
	}

	pixel, err := c.service.PlacePixel(data, ctx)
	if err != nil {
		errorDto := response.HttpErrorDto{
			StatusCode: err.StatusCode,
			Message:    err.Message,
			Errors:     err.Errors,
		}
		return ctx.Status(err.StatusCode).JSON(errorDto)
	}

	return ctx.JSON(pixel)
}
