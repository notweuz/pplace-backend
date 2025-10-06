package handler

import (
	"strings"

	"pplace_backend/internal/layer/service"
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
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Invalid request body",
		})
	}

	if errors := validation.ValidateDTO(&data); errors != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"errors": errors,
		})
	}

	token, err := h.authService.Register(c.Context(), data)
	if err != nil {
		return h.handleAuthError(c, err)
	}

	return c.Status(fiber.StatusCreated).JSON(token)
}

func (h *AuthHandler) Login(c *fiber.Ctx) error {
	var data request.AuthDto
	if err := c.BodyParser(&data); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Invalid request body",
		})
	}

	if errors := validation.ValidateDTO(&data); errors != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"errors": errors,
		})
	}

	token, err := h.authService.Login(c.Context(), data)
	if err != nil {
		return h.handleAuthError(c, err)
	}

	return c.Status(fiber.StatusOK).JSON(token)
}

func (h *AuthHandler) handleAuthError(c *fiber.Ctx, err error) error {
	errMsg := err.Error()

	switch {
	case strings.Contains(errMsg, "already exists"):
		return c.Status(fiber.StatusConflict).JSON(fiber.Map{
			"error": errMsg,
		})
	case strings.Contains(errMsg, "not found"):
		return c.Status(fiber.StatusNotFound).JSON(fiber.Map{
			"error": errMsg,
		})
	case strings.Contains(errMsg, "invalid password"):
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
			"error": errMsg,
		})
	default:
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": errMsg,
		})
	}
}
