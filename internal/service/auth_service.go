package service

import (
	"fmt"
	"github.com/golang-jwt/jwt/v5"
	"golang.org/x/crypto/bcrypt"
	"pplace_backend/internal/config"
	"pplace_backend/internal/model"
	"pplace_backend/internal/model/dto/request"
	"pplace_backend/internal/model/dto/response"
	"time"
)

type AuthService struct {
	userService *UserService
	config      *config.PPlaceConfig
}

func NewAuthService(userService *UserService, config *config.PPlaceConfig) AuthService {
	return AuthService{userService: userService, config: config}
}

func (s *AuthService) Register(dto request.AuthDto) (response.AuthTokenDto, error) {
	if _, err := s.userService.GetByUsername(dto.Username); err == nil {
		return response.AuthTokenDto{}, fmt.Errorf("username exists")
	}

	password, _ := bcrypt.GenerateFromPassword([]byte(dto.Password), bcrypt.DefaultCost)

	user := model.User{
		Username:   dto.Username,
		Password:   password,
		LastPlaced: time.Now(),
		Active:     true,
	}

	createdUser, err := s.userService.Create(&user)
	if err != nil {
		return response.AuthTokenDto{}, err
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"id":       createdUser.ID,
		"username": createdUser.Username,
		"exp":      time.Now().Add(time.Hour * time.Duration(s.config.JWT.Expiration)).Unix(),
	})

	tokenString, err := token.SignedString([]byte(s.config.JWT.Secret))
	if err != nil {
		return response.AuthTokenDto{}, fmt.Errorf("failed to generate token: %w", err)
	}

	return response.AuthTokenDto{Token: tokenString}, nil
}

func (s *AuthService) Login(dto request.AuthDto) error {
	//user, err := s.userService.GetByUsername(dto.Username)
	//if err != nil {
	//	return err
	//}

	//return c.JSON(user)
	return nil
}
