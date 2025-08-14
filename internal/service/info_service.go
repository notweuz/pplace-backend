package service

import (
	config2 "pplace_backend/internal/config"
	"pplace_backend/internal/model/dto/response"
)

type InfoService struct {
	config *config2.PPlaceConfig
}

func NewInfoService(cfg *config2.PPlaceConfig) InfoService {
	return InfoService{
		config: cfg,
	}
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
