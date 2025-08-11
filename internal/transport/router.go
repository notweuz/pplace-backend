package transport

import (
	"github.com/gofiber/fiber/v2"
	"pplace_backend/internal/controller"
)

type Router struct {
	app            *fiber.App
	userController controller.UserController
}

func NewRouter(app *fiber.App, userController controller.UserController) Router {
	router := Router{
		app:            app,
		userController: userController,
	}

	usersRoute := router.app.Group("/users")
	usersRoute.Get("/@me", router.userController.GetSelfInfo)

	return router
}
