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
		return c.Status(fiber.StatusBadRequest).JSON(
			model.NewHttpError(
				fiber.StatusBadRequest,
				"Invalid request body",
				[]string{err.Error()},
			),
		)
	}

	pixel := model.NewPixel(0, pixelCreateDto.X, pixelCreateDto.Y, pixelCreateDto.Color)
	createdPixel, err := h.service.Create(c, c.Context(), pixel)
	if err != nil {
		return h.handlePixelError(c, err)
	}

	authorDto := response.NewUserShortDto(createdPixel.UserID, createdPixel.User.Username)
	pixelDto := response.NewPixelDto(createdPixel.ID, createdPixel.X, createdPixel.Y, createdPixel.Color, *authorDto)
	return c.Status(fiber.StatusCreated).JSON(pixelDto)
}

func (h *PixelHandler) Update(c *fiber.Ctx) error {
	id, err := c.ParamsInt("id")
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(
			model.NewHttpError(fiber.StatusBadRequest, "Invalid pixel ID", []string{err.Error()}),
		)
	}

	var pixelUpdateDto request.UpdatePixelDto
	if err := c.BodyParser(&pixelUpdateDto); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(
			model.NewHttpError(fiber.StatusBadRequest, "Invalid request body", []string{err.Error()}),
		)
	}

	pixel := model.NewPixel(uint(id), 0, 0, pixelUpdateDto.Color)
	updatedPixel, err := h.service.Update(c, c.Context(), pixel)
	if err != nil {
		return h.handlePixelError(c, err)
	}

	authorDto := response.NewUserShortDto(updatedPixel.UserID, updatedPixel.User.Username)
	pixelDto := response.NewPixelDto(updatedPixel.ID, updatedPixel.X, updatedPixel.Y, updatedPixel.Color, *authorDto)
	return c.Status(fiber.StatusOK).JSON(pixelDto)
}

func (h *PixelHandler) handlePixelError(c *fiber.Ctx, err error) error {
	errMsg := err.Error()

	switch {
	case errMsg == "pixel not found":
		return c.Status(fiber.StatusNotFound).JSON(model.NewHttpError(fiber.StatusNotFound, errMsg, []string{errMsg}))
	case errMsg == "pixel already exists":
		return c.Status(fiber.StatusConflict).JSON(model.NewHttpError(fiber.StatusConflict, errMsg, []string{errMsg}))
	default:
		return c.Status(fiber.StatusInternalServerError).JSON(model.NewHttpError(fiber.StatusInternalServerError, "Failed to process pixel", []string{errMsg}))
	}
}
