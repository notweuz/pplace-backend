package layer

import (
	"gorm.io/gorm"
	"pplace_backend/internal/config"
	"pplace_backend/internal/controller"
	"pplace_backend/internal/database"
	"pplace_backend/internal/service"
)

type UserLayer struct {
	Repository *database.UserRepository
	Service    *service.UserService
	Controller *controller.UserController
}

func NewUserLayer(db *gorm.DB, config *config.PPlaceConfig) UserLayer {
	r := database.NewUserRepository(db)
	s := service.NewUserService(&r, config)
	c := controller.NewUserController(&s)

	return UserLayer{
		Repository: &r,
		Service:    &s,
		Controller: &c,
	}
}
