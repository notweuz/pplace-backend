package layer

import (
	"pplace_backend/internal/controller"
	"pplace_backend/internal/service"
)

type AuthLayer struct {
	Service    *service.AuthService
	Controller *controller.AuthController
}

func NewAuthLayer(userService *service.UserService) AuthLayer {
	s := service.NewAuthService(userService)
	c := controller.NewAuthController(&s)

	return AuthLayer{
		Service:    &s,
		Controller: &c,
	}
}
