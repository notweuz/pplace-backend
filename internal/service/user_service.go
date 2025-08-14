package service

import (
	"fmt"
	"github.com/gofiber/fiber/v2"
	"github.com/golang-jwt/jwt/v5"
	"pplace_backend/internal/config"
	"pplace_backend/internal/database"
	error2 "pplace_backend/internal/error"
	"pplace_backend/internal/model"
)

type UserService struct {
	repository *database.UserRepository
	config     *config.PPlaceConfig
}

func NewUserService(repository *database.UserRepository, conf *config.PPlaceConfig) UserService {
	return UserService{repository: repository, config: conf}
}

func (us *UserService) Create(user *model.User) (*model.User, error) {
	result, err := us.repository.Create(user)

	if err != nil {
		return nil, err
	}

	return result, nil
}

func (us *UserService) Update(user *model.User) (*model.User, *error2.HttpError) {
	result, err := us.repository.Update(user)
	if err != nil {
		return nil, error2.NewHttpError(fiber.StatusInternalServerError, "Error while updating user", err.Error())
	}
	return result, nil
}

func (us *UserService) GetByID(id uint) (model.User, error) {
	result, err := us.repository.GetById(id)
	if err != nil {
		return model.User{}, err
	}

	return *result, nil
}

func (us *UserService) GetByUsername(username string) (*model.User, error) {
	result, err := us.repository.GetByUsername(username)
	if err != nil {
		return nil, err
	}

	return result, nil
}

func (us *UserService) GetSelfInfo(ctx *fiber.Ctx) (*model.User, *error2.HttpError) {
	user, err := us.GetCurrentUser(ctx)
	if err != nil {
		return nil, err
	}
	return user, nil
}

func (us *UserService) GetCurrentUser(ctx *fiber.Ctx) (*model.User, *error2.HttpError) {
	header := ctx.Get("Authorization")
	if len(header) < 7 {
		return nil, error2.NewHttpError(fiber.StatusUnauthorized, "Invalid authorization header")
	}

	headerToken := header[7:]

	token, err := jwt.Parse(headerToken, func(token *jwt.Token) (any, error) {
		return []byte(us.config.JWT.Secret), nil
	})

	if err != nil {
		return nil, error2.NewHttpError(fiber.StatusUnauthorized, "Invalid authorization header", err.Error())
	}

	claims, ok := token.Claims.(jwt.MapClaims)
	if !ok || !token.Valid {
		return nil, error2.NewHttpError(fiber.StatusUnauthorized, "Invalid authorization header")
	}

	idClaim, exists := claims["id"]
	if !exists || idClaim == nil {
		return nil, error2.NewHttpError(fiber.StatusUnauthorized, "Invalid authorization header")
	}

	userID := uint(idClaim.(float64))

	user, err := us.repository.GetById(userID)
	if err != nil {
		return nil, error2.NewHttpError(fiber.StatusUnauthorized, fmt.Sprintf("User with id %d not found", userID), err.Error())
	}

	return user, nil
}
