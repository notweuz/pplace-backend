package transport

import (
	"github.com/gofiber/fiber/v2"
	"pplace_backend/internal/controller"
)

type Router struct {
	app             *fiber.App
	userController  *controller.UserController
	authController  *controller.AuthController
	infoController  *controller.InfoController
	pixelController *controller.PixelController
}

func NewRouter(app *fiber.App, userController *controller.UserController, authController *controller.AuthController, infoController *controller.InfoController, pixelController *controller.PixelController) Router {
	router := Router{
		app:             app,
		userController:  userController,
		authController:  authController,
		infoController:  infoController,
		pixelController: pixelController,
	}

	usersRoute := router.app.Group("/users")
	usersRoute.Get("/@me", router.userController.GetSelfInfo)

	authRoute := router.app.Group("/auth")
	authRoute.Post("/register", router.authController.Register)
	authRoute.Post("/login", router.authController.Login)

	infoRoute := router.app.Group("/info")
	infoRoute.Get("/", router.infoController.GetPixelSheetInfo)

	pixelRoute := router.app.Group("/pixels")
	pixelRoute.Get("/", router.pixelController.GetAllPixels)

	return router
}
