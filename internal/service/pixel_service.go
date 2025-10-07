package service

import (
	"context"
	"pplace_backend/internal/config"
	"pplace_backend/internal/database"
	"pplace_backend/internal/model"
	"time"

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

	if (pixel.X > s.config.Sheet.Width) || (pixel.X < 1) || (pixel.Y > s.config.Sheet.Height) || (pixel.Y < 1) {
		log.Error().Uint("x", pixel.X).Uint("y", pixel.Y).Interface("current size", s.config.Sheet).Msg("Pixel coordinates out of range")
		return nil, fiber.NewError(fiber.StatusBadRequest, "Pixel coordinates out of range")
	}

	author, err := s.userService.GetSelfInfo(c)
	if err != nil {
		return nil, err
	}

	isReady, err := s.checkPlaceCooldown(author, ctx)
	if err != nil {
		return nil, err
	}

	if !isReady {
		return nil, fiber.NewError(fiber.StatusForbidden, "Cannot create pixel, user is on cooldown")
	}

	pixel.UserID = author.ID
	author.AmountPlaced++
	author.LastPlaced = time.Now()
	_, err = s.userService.Update(ctx, author)
	if err != nil {
		log.Error().Int("amount placed", author.AmountPlaced).Time("last placed", author.LastPlaced).Err(err).Msg("Failed to update user after placing pixel")
		return nil, err
	}

	log.Info().Uint("x", pixel.X).Uint("y", pixel.Y).Interface("color", pixel.Color).Msg("Creating pixel")
	return s.database.Create(ctx, pixel)
}

func (s *PixelService) Update(c *fiber.Ctx, ctx context.Context, pixel *model.Pixel) (*model.Pixel, error) {
	author, err := s.userService.GetSelfInfo(c)
	if err != nil {
		return nil, err
	}

	isReady, err := s.checkPlaceCooldown(author, ctx)
	if err != nil {
		return nil, err
	}

	if !isReady {
		return nil, fiber.NewError(fiber.StatusForbidden, "Cannot create pixel, user is on cooldown")
	}

	oldPixel, err := s.GetByID(ctx, pixel.ID)
	if err != nil {
		log.Error().Err(err).Msg("Failed to get pixel by ID")
		return nil, err
	}

	pixel.UserID = author.ID
	pixel.X = oldPixel.X
	pixel.Y = oldPixel.Y

	author.AmountPlaced++
	author.LastPlaced = time.Now()
	_, err = s.userService.Update(ctx, author)
	if err != nil {
		log.Error().Int("amount placed", author.AmountPlaced).Time("last placed", author.LastPlaced).Err(err).Msg("Failed to update user after placing pixel")
		return nil, err
	}

	log.Info().Uint("id", pixel.ID).Uint("x", pixel.X).Uint("y", pixel.Y).Interface("color", pixel.Color).Msg("Updating pixel")
	return s.database.Update(ctx, pixel)
}

func (s *PixelService) GetByID(ctx context.Context, id uint) (*model.Pixel, error) {
	log.Info().Uint("id", id).Msg("Getting pixel by ID")
	return s.database.GetByID(ctx, id)
}

func (s *PixelService) GetAll(ctx context.Context) ([]model.Pixel, error) {
	log.Info().Msg("Getting all pixels")
	return s.database.GetAll(ctx)
}

func (s *PixelService) GetByCoordinates(ctx context.Context, x, y uint) (*model.Pixel, error) {
	log.Info().Uint("x", x).Uint("y", y).Msg("Getting pixel by coordinates")
	return s.database.GetByCoordinates(ctx, x, y)
}

func (s *PixelService) GetAllByUser(ctx context.Context, userId uint) ([]model.Pixel, error) {
	log.Info().Uint("userId", userId).Msg("Getting all pixels by user ID")
	return s.database.GetAllByUserID(ctx, userId)
}

func (s *PixelService) GetAllByUserSelf(c *fiber.Ctx, ctx context.Context) ([]model.Pixel, error) {
	user, err := s.userService.GetSelfInfo(c)
	if err != nil {
		return nil, err
	}
	log.Info().Uint("userId", user.ID).Msg("Getting all pixels by self user ID")
	return s.database.GetAllByUserID(ctx, user.ID)
}

func (s *PixelService) Delete(c *fiber.Ctx, ctx context.Context, id uint) error {
	_, err := s.userService.GetSelfInfo(c)
	if err != nil {
		return err
	}
	err = s.database.Delete(ctx, id)
	if err != nil {
		log.Error().Err(err).Uint("id", id).Msg("Failed to delete pixel")
		return err
	}
	log.Info().Uint("id", id).Msg("Deleted pixel")
	return nil
}

func (s *PixelService) checkPlaceCooldown(user *model.User, ctx context.Context) (bool, error) {
	if user.LastPlaced.IsZero() {
		return true, nil
	}

	now := time.Now()
	elapsed := now.Sub(user.LastPlaced)
	cooldown := time.Duration(s.config.Sheet.PlaceCooldown) * time.Millisecond
	canPlace := elapsed >= cooldown || user.Admin

	log.Info().
		Uint("userId", user.ID).
		Dur("elapsed", elapsed).
		Dur("cooldown", cooldown).
		Bool("canPlace", canPlace).
		Bool("isAdmin", user.Admin).
		Msg("Cooldown check")

	return canPlace, nil
}
