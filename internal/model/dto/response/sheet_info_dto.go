package response

import "pplace_backend/internal/config"

type SheetInfoDto struct {
	Version string             `json:"version"`
	Size    config.SheetConfig `json:"size"`
}
