package service

import (
	"context"
	"fmt"
	"pplace_backend/internal/config"
	"pplace_backend/internal/database"
	"pplace_backend/internal/model"
	"strings"

	"github.com/gofiber/fiber/v2"
	"github.com/golang-jwt/jwt/v5"
	"golang.org/x/crypto/bcrypt"
)

type UserService struct {
	database *database.UserDatabase
	config   *config.PPlaceConfig
}

func NewUserService(db *database.UserDatabase, c *config.PPlaceConfig) *UserService {
	return &UserService{database: db, config: c}
}

func (s *UserService) Create(ctx context.Context, user *model.User) (*model.User, error) {
	return s.database.Create(ctx, user)
}

func (s *UserService) Update(ctx context.Context, user *model.User) (*model.User, error) {
	return s.database.Update(ctx, user)
}

func (s *UserService) GetByID(ctx context.Context, id uint) (*model.User, error) {
	return s.database.GetById(ctx, id)
}

func (s *UserService) GetByUsername(ctx context.Context, username string) (*model.User, error) {
	return s.database.GetByUsername(ctx, username)
}

func (s *UserService) GetSelfInfo(c *fiber.Ctx) (*model.User, error) {
	user, ok := c.Locals("user").(*model.User)
	if !ok || user == nil {
		return nil, fmt.Errorf("user not found in context")
	}
	return user, nil
}

func (s *UserService) UpdateProfile(ctx context.Context, userID uint, username, password string) (*model.User, error) {
	currentUser, err := s.database.GetById(ctx, userID)
	if err != nil || currentUser == nil {
		return nil, fmt.Errorf("user not found: %w", err)
	}

	if username != "" {
		existingUser, err := s.database.GetByUsername(ctx, username)
		if err == nil && existingUser != nil && existingUser.ID != userID {
			return nil, fmt.Errorf("username already taken")
		}
		currentUser.Username = username
	}

	if password != "" {
		hashedPassword, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
		if err != nil {
			return nil, fmt.Errorf("error while hashing password: %w", err)
		}
		currentUser.Password = hashedPassword
	}

	return s.database.Update(ctx, currentUser)
}

func (s *UserService) ParseAndValidateToken(tokenString string) (*model.User, error) {
	token, err := jwt.Parse(tokenString, func(token *jwt.Token) (any, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
		}
		return []byte(s.config.JWT.Secret), nil
	})

	if err != nil {
		return nil, fmt.Errorf("invalid token: %w", err)
	}

	claims, ok := token.Claims.(jwt.MapClaims)
	if !ok {
		return nil, fmt.Errorf("invalid token claims")
	}

	idClaim, exists := claims["id"]
	if !exists {
		return nil, fmt.Errorf("missing user ID in token")
	}

	var userID uint
	switch v := idClaim.(type) {
	case float64:
		userID = uint(v)
	case int:
		userID = uint(v)
	default:
		return nil, fmt.Errorf("invalid user ID format")
	}

	user, err := s.database.GetById(context.Background(), userID)
	if err != nil {
		return nil, fmt.Errorf("user not found: %w", err)
	}

	return user, nil
}

func (s *UserService) ExtractToken(c *fiber.Ctx) (string, error) {
	header := c.Get("Authorization")
	if !strings.HasPrefix(header, "Bearer ") {
		return "", fmt.Errorf("invalid authorization header format")
	}
	return strings.TrimPrefix(header, "Bearer "), nil
}
