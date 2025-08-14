package service

import (
	"errors"
	"github.com/gofiber/fiber/v2"
	"gorm.io/gorm"
	"pplace_backend/internal/config"
	"pplace_backend/internal/database"
	error2 "pplace_backend/internal/error"
	"pplace_backend/internal/model"
	"pplace_backend/internal/model/dto/request"
	"pplace_backend/internal/model/dto/response"
)

type PixelService struct {
	repository  *database.PixelRepository
	config      *config.PPlaceConfig
	userService *UserService
}

func NewPixelService(repository *database.PixelRepository, config *config.PPlaceConfig, userService *UserService) PixelService {
	return PixelService{
		repository:  repository,
		userService: userService,
		config:      config,
	}
}

func (ps *PixelService) GetAllPixels() (*response.PixelsListDto, *error2.HttpError) {
	pixels, err := ps.repository.GetAll()
	if err != nil {
		return nil, error2.NewHttpError(fiber.StatusInternalServerError, "Error while getting all pixels")
	}

	pixelsDtos := make([]response.PixelShortDto, len(pixels))
	for i, pixel := range pixels {
		pixelsDtos[i] = response.PixelShortDto{
			ID:    pixel.ID,
			X:     pixel.X,
			Y:     pixel.Y,
			Color: pixel.Color,
		}
	}

	pixelsList := response.PixelsListDto{
		Pixels: pixelsDtos,
	}

	return &pixelsList, nil
}

func (ps *PixelService) GetPixelByCoordinates(x, y uint) (*model.Pixel, *error2.HttpError) {
	pixel, err := ps.repository.GetByCoordinates(x, y)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, nil
		}
		return nil, error2.NewHttpError(fiber.StatusNotFound, "Pixel not found")
	}

	return pixel, nil
}

func (ps *PixelService) Create(pixel *model.Pixel) (*model.Pixel, *error2.HttpError) {
	create, err := ps.repository.Create(pixel)
	if err != nil {
		return nil, error2.NewHttpError(fiber.StatusInternalServerError, "Error while creating pixel")
	}
	return create, nil
}

func (ps *PixelService) Update(pixel *model.Pixel) (*model.Pixel, *error2.HttpError) {
	update, err := ps.repository.Update(pixel)
	if err != nil {
		return nil, error2.NewHttpError(fiber.StatusInternalServerError, "Error while updating pixel")
	}
	return update, nil
}

func (ps *PixelService) PlacePixel(dto request.PixelPlaceDto, ctx *fiber.Ctx) (*response.PixelDto, *error2.HttpError) {
	if dto.X > ps.config.Sheet.Width || dto.Y > ps.config.Sheet.Height {
		return nil, error2.NewHttpError(fiber.StatusBadRequest, "Invalid coordinates")
	}
	if dto.Color == "" {
		return nil, error2.NewHttpError(fiber.StatusBadRequest, "Invalid color")
	}

	user, err := ps.userService.GetCurrentUser(ctx)
	if err != nil {
		return nil, error2.NewHttpError(fiber.StatusUnauthorized, "Invalid authorization header", err.Error())
	} else if user.PixelsStock <= 0 {
		return nil, error2.NewHttpError(fiber.StatusForbidden, "You have no pixels left")
	}

	pixel, err := ps.GetPixelByCoordinates(dto.X, dto.Y)
	if err != nil {
		return nil, error2.NewHttpError(fiber.StatusNotFound, "Pixel not found", err.Error())
	}

	if pixel == nil {
		pixelToCreate := &model.Pixel{
			ID:     0,
			X:      dto.X,
			Y:      dto.Y,
			Color:  dto.Color,
			UserID: user.ID,
			User:   *user,
		}
		pixel, err = ps.Create(pixelToCreate)
		if err != nil {
			return nil, error2.NewHttpError(fiber.StatusInternalServerError, "Error while creating pixel", err.Error())
		}
	} else {
		pixel.Color = dto.Color
		pixel, err = ps.Update(pixel)
		if err != nil {
			return nil, error2.NewHttpError(fiber.StatusInternalServerError, "Error while updating pixel", err.Error())
		}
	}
	user.PixelsStock--
	_, err = ps.userService.Update(user)
	if err != nil {
		return nil, err
	}

	return &response.PixelDto{
		ID:    pixel.ID,
		X:     pixel.X,
		Y:     pixel.Y,
		Color: pixel.Color,
		Author: response.UserDto{
			ID:       pixel.User.ID,
			Username: pixel.User.Username,
		},
	}, nil
}
