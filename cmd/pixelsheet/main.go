package main

import (
	"github.com/gofiber/fiber/v2"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"log"
	"pplace_backend/internal/controller"
	"pplace_backend/internal/database"
	"pplace_backend/internal/model"
	"pplace_backend/internal/service"
	"pplace_backend/internal/transport"
)

func main() {
	dsn := "host=localhost user=pplace password=pplace123 dbname=pplace-dev port=5432 sslmode=disable"
	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{})
	if err != nil {
		log.Fatal("Error connecting to database: ", err)
	}
	log.Println("Connected to database")

	err = db.AutoMigrate(&model.User{})
	if err != nil {
		log.Fatal("Migration error: ", err)
	}
	log.Println("Performed migrations")

	app := fiber.New()

	userRepository := database.NewUserRepository(db)
	userService := service.NewUserService(&userRepository)
	userController := controller.NewUserController(&userService)
	log.Println("Created layers")

	_ = transport.NewRouter(app, userController)
	log.Println("Added routes")

	log.Fatal(app.Listen(":8000"))
}
