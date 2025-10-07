package service

import (
	"context"
	"pplace_backend/internal/config"
	"pplace_backend/internal/database"
	"pplace_backend/internal/model"

	"github.com/gofiber/fiber/v2"
)

type PixelService struct {
	database    *database.PixelDatabase
	config      *config.PPlaceConfig
	userService *UserService
}

func NewPixelService(db *database.PixelDatabase, config *config.PPlaceConfig, userService *UserService) *PixelService {
	return &PixelService{
		database:    db,
		config:      config,
		userService: userService,
	}
}

func (s *PixelService) Create(c *fiber.Ctx, ctx context.Context, pixel *model.Pixel) (*model.Pixel, error) {
	author, err := s.userService.GetSelfInfo(c)
	if err != nil {
		return nil, err
	}
	pixel.UserID = author.ID
	return s.database.Create(ctx, pixel)
}

func (s *PixelService) Update(ctx context.Context, pixel *model.Pixel) (*model.Pixel, error) {
	return s.database.Update(ctx, pixel)
}

func (s *PixelService) GetByID(ctx context.Context, id uint) (*model.Pixel, error) {
	return s.database.GetByID(ctx, id)
}

func (s *PixelService) GetAll(ctx context.Context) ([]model.Pixel, error) {
	return s.database.GetAll(ctx)
}

func (s *PixelService) GetByCoordinates(ctx context.Context, x, y uint) (*model.Pixel, error) {
	return s.database.GetByCoordinates(ctx, x, y)
}

func (s *PixelService) GetAllByUser(ctx context.Context, userId uint) ([]model.Pixel, error) {
	return s.database.GetAllByUserID(ctx, userId)
}

func (s *PixelService) GetAllByUserSelf(c *fiber.Ctx, ctx context.Context) ([]model.Pixel, error) {
	user, err := s.userService.GetSelfInfo(c)
	if err != nil {
		return nil, err
	}
	return s.database.GetAllByUserID(ctx, user.ID)
}
