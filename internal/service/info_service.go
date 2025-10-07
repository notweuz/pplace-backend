package service

import (
	"pplace_backend/internal/config"
	"pplace_backend/internal/model/dto/response"
)

type InfoService struct {
	config *config.PPlaceConfig
}

func NewInfoService(config *config.PPlaceConfig) *InfoService {
	return &InfoService{config: config}
}

func (s *InfoService) GetPixelSheetInfo() response.SheetInfoDto {
	return response.SheetInfoDto{
		Size: map[string]int{
			"width":  int(s.config.Sheet.Width),
			"height": int(s.config.Sheet.Height),
		},
		Version: s.config.Version,
	}
}
