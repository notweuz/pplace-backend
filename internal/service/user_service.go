package service

import (
	"github.com/gofiber/fiber/v2"
	"pplace_backend/internal/database"
	"pplace_backend/internal/model"
)

type UserService struct {
	repository *database.UserRepository
}

func NewUserService(repository *database.UserRepository) UserService {
	return UserService{repository: repository}
}

func (us *UserService) Create(user *model.User) (*model.User, error) {
	result, err := us.repository.Create(user)

	if err != nil {
		return nil, err
	}

	return result, nil
}

func (us *UserService) GetSelfInfo(ctx *fiber.Ctx) (*model.User, error) {
	return nil, nil
}
