package service

import (
	"pplace_backend/internal/config"
	"pplace_backend/internal/model/dto/response"

	"github.com/rs/zerolog/log"
)

type InfoService struct {
	config *config.PPlaceConfig
}

func NewInfoService(config *config.PPlaceConfig) *InfoService {
	return &InfoService{config: config}
}

func (s *InfoService) GetPixelSheetInfo() response.SheetInfoDto {
	log.Info().Interface("version", s.config.Version).Interface("sheet", s.config.Sheet).Msg("Fetching service info")
	return response.SheetInfoDto{
		Size: map[string]int{
			"width":          int(s.config.Sheet.Width),
			"height":         int(s.config.Sheet.Height),
			"place_cooldown": int(s.config.Sheet.PlaceCooldown),
		},
		Version: s.config.Version,
	}
}
