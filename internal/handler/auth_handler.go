package handler

import (
	"pplace_backend/internal/model"
	"pplace_backend/internal/service"
	"strings"

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
		return c.Status(fiber.StatusBadRequest).JSON(
			model.NewHttpError(
				fiber.StatusBadRequest,
				"Wrong request body provided",
				[]string{err.Error()},
			),
		)
	}

	if errors := validation.ValidateDTO(&data); errors != nil {
		stringErrors := make([]string, len(errors))
		for i, err := range errors {
			stringErrors[i] = err.Error
		}
		return c.Status(fiber.StatusBadRequest).JSON(
			model.NewHttpError(
				fiber.StatusBadRequest,
				"Request body validation failed",
				stringErrors,
			),
		)
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
		return c.Status(fiber.StatusBadRequest).JSON(
			model.NewHttpError(
				fiber.StatusBadRequest,
				"Wrong request body provided",
				[]string{err.Error()},
			),
		)
	}

	if errors := validation.ValidateDTO(&data); errors != nil {
		stringErrors := make([]string, len(errors))
		for i, err := range errors {
			stringErrors[i] = err.Error
		}
		return c.Status(fiber.StatusBadRequest).JSON(
			model.NewHttpError(
				fiber.StatusBadRequest,
				"Request body validation failed",
				stringErrors,
			),
		)
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
		return c.Status(fiber.StatusConflict).JSON(model.NewHttpError(fiber.StatusConflict, errMsg, []string{err.Error()}))
	case strings.Contains(errMsg, "not found"):
		return c.Status(fiber.StatusNotFound).JSON(model.NewHttpError(fiber.StatusNotFound, errMsg, []string{err.Error()}))
	case strings.Contains(errMsg, "invalid password"):
		return c.Status(fiber.StatusUnauthorized).JSON(model.NewHttpError(fiber.StatusUnauthorized, errMsg, []string{err.Error()}))
	default:
		return c.Status(fiber.StatusInternalServerError).JSON(model.NewHttpError(fiber.StatusInternalServerError, errMsg, []string{err.Error()}))
	}
}
