package service

import (
	"github.com/gofiber/fiber/v2"
	"pplace_backend/internal/config"
	"pplace_backend/internal/database"
	error2 "pplace_backend/internal/error"
	"pplace_backend/internal/model/dto/response"
)

type PixelService struct {
	repository *database.PixelRepository
	config     *config.PPlaceConfig
}

func NewPixelService(repository *database.PixelRepository, config *config.PPlaceConfig) PixelService {
	return PixelService{
		repository: repository,
		config:     config,
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
