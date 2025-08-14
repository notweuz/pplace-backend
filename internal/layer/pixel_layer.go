package layer

import (
	"gorm.io/gorm"
	"pplace_backend/internal/config"
	"pplace_backend/internal/controller"
	"pplace_backend/internal/database"
	"pplace_backend/internal/service"
)

type PixelLayer struct {
	Repository *database.PixelRepository
	Service    *service.PixelService
	Controller *controller.PixelController
}

func NewPixelLayer(db *gorm.DB, config *config.PPlaceConfig, userService *service.UserService) PixelLayer {
	r := database.NewPixelRepository(db)
	s := service.NewPixelService(&r, config, userService)
	c := controller.NewPixelController(&s)

	return PixelLayer{
		Repository: &r,
		Service:    &s,
		Controller: &c,
	}
}
