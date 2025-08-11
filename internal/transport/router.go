package transport

import (
	"github.com/gofiber/fiber/v2"
	"pplace_backend/internal/controller"
)

type Router struct {
	app            *fiber.App
	userController *controller.UserController
	authController *controller.AuthController
}

func NewRouter(app *fiber.App, userController *controller.UserController, authController *controller.AuthController) Router {
	router := Router{
		app:            app,
		userController: userController,
		authController: authController,
	}

	usersRoute := router.app.Group("/users")
	usersRoute.Get("/@me", router.userController.GetSelfInfo)

	authRoute := router.app.Group("/auth")
	authRoute.Post("/register", router.authController.Register)

	return router
}
