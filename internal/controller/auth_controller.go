package controller

import (
	"github.com/gofiber/fiber/v2"
	"pplace_backend/internal/model/dto/request"
	"pplace_backend/internal/service"
	"pplace_backend/internal/validation"
)

type AuthController struct {
	service *service.AuthService
}

func NewAuthController(service *service.AuthService) AuthController {
	return AuthController{service: service}
}

func (c *AuthController) Register(ctx *fiber.Ctx) error {
	var data request.AuthDto

	if err := ctx.BodyParser(&data); err != nil {
		return err
	}

	errors := validation.ValidateDTO(&data)
	if errors != nil {
		return ctx.Status(fiber.StatusBadRequest).JSON(errors)
	}

	return ctx.JSON(data)
}
