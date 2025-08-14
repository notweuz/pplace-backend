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

func (c *PixelController) GetPixelByCoordinates(ctx *fiber.Ctx) error {
	x, err := ctx.ParamsInt("x")
	if err != nil {
		errorDto := response.HttpErrorDto{
			StatusCode: fiber.StatusBadRequest,
			Message:    "Invalid x parameter",
			Errors:     []string{"x parameter must be an integer"},
		}
		return ctx.Status(fiber.StatusBadRequest).JSON(errorDto)
	}

	y, err := ctx.ParamsInt("y")
	if err != nil {
		errorDto := response.HttpErrorDto{
			StatusCode: fiber.StatusBadRequest,
			Message:    "Invalid y parameter",
			Errors:     []string{"y parameter must be an integer"},
		}
		return ctx.Status(fiber.StatusBadRequest).JSON(errorDto)
	}

	pixel, _ := c.service.GetPixelByCoordinates(uint(x), uint(y))
	if pixel == nil {
		errorDto := response.HttpErrorDto{
			StatusCode: fiber.StatusNotFound,
			Message:    "Pixel not found",
			Errors:     []string{},
		}
		return ctx.Status(fiber.StatusNotFound).JSON(errorDto)
	}

	pixelDto := &response.PixelDto{
		ID:    pixel.ID,
		X:     pixel.X,
		Y:     pixel.Y,
		Color: pixel.Color,
		Author: response.UserDto{
			ID:       pixel.User.ID,
			Username: pixel.User.Username,
		},
	}

	return ctx.JSON(pixelDto)
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
