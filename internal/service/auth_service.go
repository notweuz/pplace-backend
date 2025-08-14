package service

import (
	"github.com/gofiber/fiber/v2"
	"github.com/golang-jwt/jwt/v5"
	"golang.org/x/crypto/bcrypt"
	"pplace_backend/internal/config"
	error2 "pplace_backend/internal/error"
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

func (s *AuthService) Register(dto request.AuthDto) (*response.AuthTokenDto, *error2.HttpError) {
	user, err := s.userService.GetByUsername(dto.Username)
	if err != nil {
		return nil, error2.NewHttpError(fiber.StatusInternalServerError, "Error while getting user by username", err.Error())
	}

	if user != nil {
		return nil, error2.NewHttpError(fiber.StatusConflict, "User with this username already exists")
	}

	password, _ := bcrypt.GenerateFromPassword([]byte(dto.Password), bcrypt.DefaultCost)

	user = &model.User{
		Username:   dto.Username,
		Password:   password,
		LastPlaced: time.Now(),
		Active:     true,
	}

	createdUser, err := s.userService.Create(user)
	if err != nil {
		return nil, error2.NewHttpError(fiber.StatusInternalServerError, "Error while creating user", err.Error())
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"id":       createdUser.ID,
		"username": createdUser.Username,
		"exp":      time.Now().Add(time.Hour * time.Duration(s.config.JWT.Expiration)).Unix(),
	})

	tokenString, err := token.SignedString([]byte(s.config.JWT.Secret))
	if err != nil {
		return nil, error2.NewHttpError(fiber.StatusInternalServerError, "Error while signing token", err.Error())
	}

	return &response.AuthTokenDto{Token: tokenString}, nil
}

func (s *AuthService) Login(dto request.AuthDto) (*response.AuthTokenDto, *error2.HttpError) {
	user, err := s.userService.GetByUsername(dto.Username)
	if err != nil || user == nil {
		return nil, error2.NewHttpError(fiber.StatusNotFound, "User with that username not found")
	}

	if err := bcrypt.CompareHashAndPassword(user.Password, []byte(dto.Password)); err != nil {
		return nil, error2.NewHttpError(fiber.StatusUnauthorized, "Invalid password")
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"id":       user.ID,
		"username": user.Username,
		"exp":      time.Now().Add(time.Hour * time.Duration(s.config.JWT.Expiration)).Unix(),
	})

	tokenString, err := token.SignedString([]byte(s.config.JWT.Secret))
	if err != nil {
		return nil, error2.NewHttpError(fiber.StatusInternalServerError, "Error while signing token", err.Error())
	}

	return &response.AuthTokenDto{Token: tokenString}, nil
}
