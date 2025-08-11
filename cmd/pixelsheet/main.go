package main

import (
	"fmt"
	"github.com/gofiber/fiber/v2"
	"gopkg.in/yaml.v3"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"log"
	"os"
	config2 "pplace_backend/internal/config"
	"pplace_backend/internal/controller"
	"pplace_backend/internal/database"
	"pplace_backend/internal/model"
	"pplace_backend/internal/service"
	"pplace_backend/internal/transport"
	"strconv"
)

func main() {
	data, err := os.ReadFile("configs/application.yml")
	if err != nil {
		log.Fatal("Error while getting YAML config file: ", err)
	}

	var config config2.Config
	err = yaml.Unmarshal(data, &config)
	if err != nil {
		log.Fatal("Error while parsing config file: ", err)
	}
	log.Println("Loaded config")

	dbConfig := config.PPlace.Database

	dsn := fmt.Sprintf("host=%v user=%v password=%v dbname=%v port=%v sslmode=disable", dbConfig.Host, dbConfig.User, dbConfig.Password, dbConfig.DBName, dbConfig.Port) // "host=localhost user= password=pplace123 dbname=pplace-dev port=5432 sslmode=disable"
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

	log.Fatal(app.Listen(":" + strconv.Itoa(int(config.PPlace.Port))))
}
