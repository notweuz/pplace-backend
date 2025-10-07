package handler

import (
	"pplace_backend/internal/model"
	"pplace_backend/internal/service"

	"pplace_backend/internal/model/dto/request"
	"pplace_backend/internal/validation"

	"github.com/gofiber/fiber/v2"
)

type AuthHandler struct {
	authService *service.AuthService
}

func NewAuthHandler(authService *service.AuthService) *AuthHandler {
	return &AuthHandler{authService: authService}
}

func (h *AuthHandler) Register(c *fiber.Ctx) error {
	var data request.AuthDto
	if err := c.BodyParser(&data); err != nil {
		return model.NewHttpError(fiber.StatusBadRequest, "Wrong request body provided", []string{err.Error()})
	}

	if errors := validation.ValidateDTO(&data); errors != nil {
		stringErrors := make([]string, len(errors))
		for i, err := range errors {
			stringErrors[i] = err.Error
		}
		return model.NewHttpError(fiber.StatusBadRequest, "Request body validation failed", stringErrors)
	}

	token, err := h.authService.Register(c.Context(), data)
	if err != nil {
		return err
	}

	return c.Status(fiber.StatusCreated).JSON(token)
}

func (h *AuthHandler) Login(c *fiber.Ctx) error {
	var data request.AuthDto
	if err := c.BodyParser(&data); err != nil {
		return model.NewHttpError(fiber.StatusBadRequest, "Wrong request body provided", []string{err.Error()})
	}

	if errors := validation.ValidateDTO(&data); errors != nil {
		stringErrors := make([]string, len(errors))
		for i, err := range errors {
			stringErrors[i] = err.Error
		}
		return model.NewHttpError(fiber.StatusBadRequest, "Request body validation failed", stringErrors)
	}

	token, err := h.authService.Login(c.Context(), data)
	if err != nil {
		return err
	}

	return c.Status(fiber.StatusOK).JSON(token)
}
