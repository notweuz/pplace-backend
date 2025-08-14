package layer

import (
	"pplace_backend/internal/config"
	"pplace_backend/internal/controller"
	"pplace_backend/internal/service"
)

type AuthLayer struct {
	Service    *service.AuthService
	Controller *controller.AuthController
}

func NewAuthLayer(userService *service.UserService, config *config.PPlaceConfig) AuthLayer {
	s := service.NewAuthService(userService, config)
	c := controller.NewAuthController(&s)

	return AuthLayer{
		Service:    &s,
		Controller: &c,
	}
}
