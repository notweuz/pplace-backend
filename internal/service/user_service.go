package service

import (
	"github.com/gofiber/fiber/v2"
	"pplace_backend/internal/database"
	"pplace_backend/internal/model"
)

type UserService struct {
	Repository *database.UserRepository
}

func NewUserService(repository *database.UserRepository) UserService {
	return UserService{Repository: repository}
}

func (us *UserService) Create(user *model.User) (*model.User, error) {
	result, err := us.Repository.Create(user)

	if err != nil {
		return nil, err
	}

	return result, nil
}

func (us *UserService) GetSelfInfo(ctx *fiber.Ctx) (*model.User, error) {
	return nil, nil
}
