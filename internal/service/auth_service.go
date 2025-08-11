package service

import (
	"github.com/gofiber/fiber/v2"
	"golang.org/x/crypto/bcrypt"
	"pplace_backend/internal/model"
	"pplace_backend/internal/model/dto/request"
	"time"
)

type AuthService struct {
	userService *UserService
}

func NewAuthService(userService *UserService) AuthService {
	return AuthService{userService: userService}
}

func (s *AuthService) Register(dto request.AuthDto) error {
	password, _ := bcrypt.GenerateFromPassword([]byte(dto.Password), bcrypt.DefaultCost)

	user := model.User{
		Username:   dto.Username,
		Password:   password,
		LastPlaced: time.Now(),
		Active:     true,
	}

	_, err := s.userService.Create(&user)
	if err != nil {
		return err
	}

	return nil
}

func (s *AuthService) Login(dto request.AuthDto) error {
	user, err := s.userService.GetByUsername(dto.Username)
	if err != nil {
		return err
	}

	if err := bcrypt.CompareHashAndPassword(user.Password, []byte(data["password"])); err != nil {
		c.Status(fiber.StatusBadRequest)
		return c.JSON(fiber.Map{
			"message": "incorrect password",
		})
	}

	return c.JSON(user)

}
