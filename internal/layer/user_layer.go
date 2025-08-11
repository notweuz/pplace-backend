package layer

import (
	"pplace_backend/internal/controller"
	"pplace_backend/internal/database"
	"pplace_backend/internal/service"
)

type UserLayer struct {
	Repository *database.UserRepository
	Service    *service.UserService
	Controller *controller.UserController
}

func NewUserLayer(r *database.UserRepository, s *service.UserService, c *controller.UserController) UserLayer {
	return UserLayer{
		Repository: r,
		Service:    s,
		Controller: c,
	}
}
