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

	token, err := c.service.Register(data)
	if err != nil {
		if err.Error() == "username exists" {
			return ctx.Status(fiber.StatusConflict).JSON(fiber.Map{
				"error": "Username already exists",
			})
		}

		return ctx.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": err.Error(),
		})
	}

	return ctx.Status(fiber.StatusCreated).JSON(token)
}
