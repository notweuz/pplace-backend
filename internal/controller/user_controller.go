package controller

import (
	"github.com/gofiber/fiber/v2"
	"log"
	"pplace_backend/internal/model/dto/response"
	"pplace_backend/internal/service"
)

type UserController struct {
	Service *service.UserService
}

func NewUserController(service *service.UserService) UserController {
	return UserController{Service: service}
}

func (uc *UserController) GetSelfInfo(ctx *fiber.Ctx) error {
	info, err := uc.Service.GetSelfInfo(ctx)
	if err != nil {
		log.Println("GetSelfInfo error:", err)
		return ctx.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "Failed to get user info",
		})
	}

	userDto := response.UserDto{
		Id:       info.ID,
		Username: info.Username,
	}

	return ctx.JSON(userDto)
}
