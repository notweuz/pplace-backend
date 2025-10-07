package service

import (
	"context"
	"pplace_backend/internal/config"
	"pplace_backend/internal/database"
	"pplace_backend/internal/model"

	"github.com/gofiber/fiber/v2"
	"github.com/rs/zerolog/log"
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
	_, err := s.GetByCoordinates(ctx, pixel.X, pixel.Y)
	if err == nil {
		log.Error().Uint("x", pixel.X).Uint("y", pixel.Y).Msg("Cannot create pixel, pixel on that place already exists")
		return nil, fiber.NewError(fiber.StatusConflict, "Already exists")
	}

	author, err := s.userService.GetSelfInfo(c)
	if err != nil {
		log.Error().Err(err).Msg("")
		return nil, err
	}
	pixel.UserID = author.ID
	return s.database.Create(ctx, pixel)
}

func (s *PixelService) Update(c *fiber.Ctx, ctx context.Context, pixel *model.Pixel) (*model.Pixel, error) {
	author, err := s.userService.GetSelfInfo(c)
	if err != nil {
		log.Error().Err(err).Msg("")
		return nil, err
	}

	oldPixel, err := s.GetByID(ctx, pixel.ID)
	if err != nil {
		log.Error().Err(err).Msg("")
		return nil, err
	}

	pixel.UserID = author.ID
	pixel.X = oldPixel.X
	pixel.Y = oldPixel.Y
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
