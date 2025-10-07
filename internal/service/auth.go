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
	"github.com/rs/zerolog/log"
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
		log.Error().Err(err).Msgf("UserService GetByUsername failed: %v", err)
		return nil, fmt.Errorf("error while getting user by username: %w", err)
	}
	if user != nil {
		log.Warn().Msgf("User %s already exists", user.Username)
		return nil, fmt.Errorf("user with this username already exists")
	}

	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(dto.Password), bcrypt.DefaultCost)
	if err != nil {
		log.Error().Err(err).Msg("Failed to hash password")
		return nil, fmt.Errorf("error while hashing password: %w", err)
	}

	newUser := model.NewUser(dto.Username, string(hashedPassword))

	createdUser, err := s.userService.Create(ctx, newUser)
	if err != nil {
		log.Error().Err(err).Msgf("Failed to create user %s", dto.Username)
		return nil, fmt.Errorf("error while creating user: %w", err)
	}

	tokenString, err := s.generateToken(createdUser)
	if err != nil {
		log.Error().Err(err).Msg("Failed to generate token")
		return nil, fmt.Errorf("error while generating token: %w", err)
	}

	log.Info().Msgf("User %s registered successfully", createdUser.Username)
	return response.NewAuthTokenDto(tokenString), nil
}

func (s *AuthService) Login(ctx context.Context, dto request.AuthDto) (*response.AuthTokenDto, error) {
	user, err := s.userService.GetByUsername(ctx, dto.Username)
	if err != nil {
		log.Error().Err(err).Msgf("UserService GetByUsername failed: %v", err)
		return nil, fmt.Errorf("error while getting user: %w", err)
	}
	if user == nil {
		log.Warn().Msgf("User with username %s not found", dto.Username)
		return nil, fmt.Errorf("user with that username not found")
	}

	if err := bcrypt.CompareHashAndPassword(user.Password, []byte(dto.Password)); err != nil {
		log.Warn().Msgf("Invalid password for user %s", dto.Username)
		return nil, fmt.Errorf("invalid password")
	}

	tokenString, err := s.generateToken(user)
	if err != nil {
		log.Error().Err(err).Msg("Failed to generate token")
		return nil, fmt.Errorf("error while generating token: %w", err)
	}

	log.Info().Msgf("User %s logged in successfully", user.Username)
	return response.NewAuthTokenDto(tokenString), nil
}

func (s *AuthService) generateToken(user *model.User) (string, error) {
	claims := model.UserClaims{
		ID: user.ID,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(time.Hour * time.Duration(s.config.JWT.Expiration))),
			IssuedAt:  jwt.NewNumericDate(time.Now()),
			Issuer:    "pplace_backend",
			Subject:   fmt.Sprintf("%d", user.ID),
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString([]byte(s.config.JWT.Secret))
}
