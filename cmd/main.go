package main

import (
	"fmt"
	"os"

	config2 "pplace_backend/internal/config"
	"pplace_backend/internal/database"
	"pplace_backend/internal/middleware"
	"pplace_backend/internal/model"
	service2 "pplace_backend/internal/service"
	"pplace_backend/internal/transport"
	"pplace_backend/internal/ws"

	"github.com/gofiber/fiber/v2"
	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
	"github.com/rs/zerolog/pkgerrors"
	"gopkg.in/yaml.v3"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

func main() {
	setupLogger()
	log.Info().Msg("Starting pplace server")

	data, err := os.ReadFile("configs/application.yml")
	if err != nil {
		log.Fatal().Err(err).Msg("Failed to read application.yml")
	}

	log.Info().Msg("Loaded application.yml")
	var config config2.Config
	err = yaml.Unmarshal(data, &config)
	if err != nil {
		log.Fatal().Err(err).Msg("Failed to parse application.yml")
	}

	log.Info().Msg("Parsed application configuration")
	dbConfig := config.PPlace.Database
	dsn := fmt.Sprintf("host=%v user=%v password=%v dbname=%v port=%v sslmode=disable", dbConfig.Host, dbConfig.User, dbConfig.Password, dbConfig.DBName, dbConfig.Port)
	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{
		Logger: logger.Default.LogMode(logger.Silent),
	})
	if err != nil {
		log.Fatal().Err(err).Msg("Failed to connect to database")
	}

	log.Info().Msg("Connected to database")
	err = db.AutoMigrate(&model.User{}, &model.Pixel{})
	if err != nil {
		log.Fatal().Err(err).Msg("Failed to migrate database")
	}

	log.Info().Msg("Migrated database successfully")
	app := fiber.New()
	app.Use(middleware.LoggingMiddleware())
	log.Info().Msg("Initializing fiber application")

	userService := setupUserService(db, &config.PPlace)
	authService := setupAuthService(db, &config.PPlace, userService)
	pixelService := setupPixelService(db, &config.PPlace, userService)
	infoService := setupInfoService(&config.PPlace)

	ws.Start()

	api := app.Group("/api")
	transport.SetupUserRoutes(api, userService)
	transport.SetupAuthRoutes(api, authService)
	transport.SetupPixelRoutes(api, pixelService, userService)
	transport.SetupInfoRoutes(api, infoService)

	log.Info().Msgf("Starting server on port %d", config.PPlace.Port)
	log.Fatal().Err(app.Listen(fmt.Sprintf(":%d", config.PPlace.Port)))
}

func setupLogger() {
	zerolog.TimeFieldFormat = zerolog.TimeFormatUnix
	zerolog.ErrorStackMarshaler = pkgerrors.MarshalStack
	log.Logger = log.Output(zerolog.ConsoleWriter{Out: os.Stderr})
	zerolog.SetGlobalLevel(zerolog.InfoLevel)
}

func setupUserService(db *gorm.DB, c *config2.PPlaceConfig) *service2.UserService {
	userDatabase := database.NewUserDatabase(db)
	userService := service2.NewUserService(userDatabase, c)
	return userService
}

func setupAuthService(db *gorm.DB, c *config2.PPlaceConfig, us *service2.UserService) *service2.AuthService {
	authService := service2.NewAuthService(us, c)
	return authService
}

func setupPixelService(db *gorm.DB, c *config2.PPlaceConfig, us *service2.UserService) *service2.PixelService {
	pixelDatabase := database.NewPixelDatabase(db)
	pixelService := service2.NewPixelService(pixelDatabase, c, us)
	return pixelService
}

func setupInfoService(c *config2.PPlaceConfig) *service2.InfoService {
	infoService := service2.NewInfoService(c)
	return infoService
}
