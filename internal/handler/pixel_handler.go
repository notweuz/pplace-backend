package handler

import (
	"pplace_backend/internal/model"
	"pplace_backend/internal/model/dto/request"
	"pplace_backend/internal/model/dto/response"
	"pplace_backend/internal/service"

	"github.com/gofiber/fiber/v2"
)

type PixelHandler struct {
	service *service.PixelService
}

func NewPixelHandler(service *service.PixelService) *PixelHandler {
	return &PixelHandler{service: service}
}

func (h *PixelHandler) Create(c *fiber.Ctx) error {
	var pixelCreateDto request.PlacePixelDto
	if err := c.BodyParser(&pixelCreateDto); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Invalid request body",
		})
	}

	pixel := model.NewPixel(pixelCreateDto.X, pixelCreateDto.Y, pixelCreateDto.Color)
	createdPixel, err := h.service.Create(c, c.Context(), pixel)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "Failed to create pixel",
		})
	}

	authorDto := response.NewUserShortDto(createdPixel.UserID, createdPixel.User.Username)
	pixelDto := response.NewPixelDto(createdPixel.ID, createdPixel.X, createdPixel.Y, createdPixel.Color, *authorDto)
	return c.Status(fiber.StatusCreated).JSON(pixelDto)
}
