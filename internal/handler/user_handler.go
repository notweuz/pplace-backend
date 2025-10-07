package handler

import (
	"pplace_backend/internal/model"
	"pplace_backend/internal/model/dto/request"
	"pplace_backend/internal/model/dto/response"
	"pplace_backend/internal/service"
	"pplace_backend/internal/validation"
	"strings"

	"github.com/gofiber/fiber/v2"
)

type UserHandler struct {
	service *service.UserService
}

func NewUserHandler(userService *service.UserService) *UserHandler {
	return &UserHandler{service: userService}
}

func (h *UserHandler) GetSelfInfo(c *fiber.Ctx) error {
	user, err := h.service.GetSelfInfo(c)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": err.Error(),
		})
	}

	userDto := response.NewUserDto(user.ID, user.Username, user.LastPlaced)
	return c.JSON(userDto)
}

func (h *UserHandler) CreateUser(c *fiber.Ctx) error {
	var user model.User
	if err := c.BodyParser(&user); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Invalid request body",
		})
	}

	createdUser, err := h.service.Create(c.Context(), &user)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": err.Error(),
		})
	}

	userDto := response.NewUserDto(createdUser.ID, createdUser.Username, createdUser.LastPlaced)
	return c.Status(fiber.StatusCreated).JSON(userDto)
}

func (h *UserHandler) UpdateUser(c *fiber.Ctx) error {
	var updateData request.UpdateUserDto
	if err := c.BodyParser(&updateData); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Invalid request body",
		})
	}

	if errors := validation.ValidateDTO(&updateData); errors != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"errors": errors,
		})
	}

	if updateData.Username == "" && updateData.Password == "" {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "At least one field (username or password) must be provided",
		})
	}

	currentUser, err := h.service.GetSelfInfo(c)
	if err != nil {
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
			"error": "Unauthorized",
		})
	}

	updatedUser, err := h.service.UpdateProfile(c.Context(), currentUser.ID, updateData.Username, updateData.Password)
	if err != nil {
		return h.handleUpdateError(c, err)
	}

	userDto := response.NewUserDto(updatedUser.ID, updatedUser.Username, updatedUser.LastPlaced)
	return c.JSON(userDto)
}

func (h *UserHandler) GetUserByID(c *fiber.Ctx) error {
	id, err := c.ParamsInt("id")
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Invalid user ID",
		})
	}

	user, err := h.service.GetByID(c.Context(), uint(id))
	if err != nil {
		return c.Status(fiber.StatusNotFound).JSON(fiber.Map{
			"error": "User not found",
		})
	}

	userDto := response.NewUserDto(user.ID, user.Username, user.LastPlaced)
	return c.JSON(userDto)
}

func (h *UserHandler) GetUserByUsername(c *fiber.Ctx) error {
	username := c.Params("username")
	if username == "" {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Username is required",
		})
	}

	user, err := h.service.GetByUsername(c.Context(), username)
	if err != nil {
		return c.Status(fiber.StatusNotFound).JSON(fiber.Map{
			"error": "User not found",
		})
	}

	userDto := response.NewUserDto(user.ID, user.Username, user.LastPlaced)
	return c.JSON(userDto)
}

func (h *UserHandler) handleUpdateError(c *fiber.Ctx, err error) error {
	errMsg := err.Error()
	switch {
	case strings.Contains(errMsg, "already taken"):
		return c.Status(fiber.StatusConflict).JSON(fiber.Map{
			"error": errMsg,
		})
	case strings.Contains(errMsg, "not found"):
		return c.Status(fiber.StatusNotFound).JSON(fiber.Map{
			"error": errMsg,
		})
	default:
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": errMsg,
		})
	}
}
