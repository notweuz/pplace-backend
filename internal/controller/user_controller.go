package controller

import (
	"github.com/gofiber/fiber/v2"
	"pplace_backend/internal/model/dto/response"
	"pplace_backend/internal/service"
)

type UserController struct {
	service *service.UserService
}

func NewUserController(service *service.UserService) UserController {
	return UserController{service: service}
}

func (uc *UserController) GetSelfInfo(ctx *fiber.Ctx) error {
	info, err := uc.service.GetSelfInfo(ctx)
	if err != nil {
		errorDto := response.HttpErrorDto{
			StatusCode: err.StatusCode,
			Message:    err.Message,
		}
		return ctx.Status(err.StatusCode).JSON(errorDto)
	}

	userDto := response.UserDto{
		Id:       info.ID,
		Username: info.Username,
	}

	return ctx.JSON(userDto)
}
