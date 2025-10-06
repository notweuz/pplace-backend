package service

import (
	"context"
	"fmt"
	"time"

	"pplace_backend/internal/config"
	"pplace_backend/internal/model"
	"pplace_backend/internal/model/dto/request"
	"pplace_backend/internal/model/dto/response"

	"github.com/golang-jwt/jwt/v5"
	"golang.org/x/crypto/bcrypt"
)

type AuthService struct {
	userService *UserService
	config      *config.PPlaceConfig
}

func NewAuthService(userService *UserService, config *config.PPlaceConfig) *AuthService {
	return &AuthService{userService: userService, config: config}
}

func (s *AuthService) Register(ctx context.Context, dto request.AuthDto) (*response.AuthTokenDto, error) {
	user, err := s.userService.GetByUsername(ctx, dto.Username)
	if err != nil {
		return nil, fmt.Errorf("error while getting user by username: %w", err)
	}
	if user != nil {
		return nil, fmt.Errorf("user with this username already exists")
	}

	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(dto.Password), bcrypt.DefaultCost)
	if err != nil {
		return nil, fmt.Errorf("error while hashing password: %w", err)
	}

	newUser := model.NewUser(dto.Username, string(hashedPassword))

	createdUser, err := s.userService.Create(ctx, newUser)
	if err != nil {
		return nil, fmt.Errorf("error while creating user: %w", err)
	}

	tokenString, err := s.generateToken(createdUser)
	if err != nil {
		return nil, fmt.Errorf("error while generating token: %w", err)
	}

	return response.NewAuthTokenDto(tokenString), nil
}

func (s *AuthService) Login(ctx context.Context, dto request.AuthDto) (*response.AuthTokenDto, error) {
	user, err := s.userService.GetByUsername(ctx, dto.Username)
	if err != nil {
		return nil, fmt.Errorf("error while getting user: %w", err)
	}
	if user == nil {
		return nil, fmt.Errorf("user with that username not found")
	}

	if err := bcrypt.CompareHashAndPassword(user.Password, []byte(dto.Password)); err != nil {
		return nil, fmt.Errorf("invalid password")
	}

	tokenString, err := s.generateToken(user)
	if err != nil {
		return nil, fmt.Errorf("error while generating token: %w", err)
	}

	return response.NewAuthTokenDto(tokenString), nil
}

func (s *AuthService) generateToken(user *model.User) (string, error) {
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"id":       user.ID,
		"username": user.Username,
		"exp":      time.Now().Add(time.Hour * time.Duration(s.config.JWT.Expiration)).Unix(),
	})

	tokenString, err := token.SignedString([]byte(s.config.JWT.Secret))
	if err != nil {
		return "", err
	}

	return tokenString, nil
}
