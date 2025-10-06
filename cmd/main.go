package main

import (
	"fmt"
	"os"
	config2 "pplace_backend/internal/config"

	"github.com/gofiber/fiber/v2"
	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
	"github.com/rs/zerolog/pkgerrors"
	"gopkg.in/yaml.v3"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
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
	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{})
	if err != nil {
		log.Fatal().Err(err).Msg("Failed to connect to database")
	}
	log.Info().Msg("Connected to database")

	err = db.AutoMigrate()
	if err != nil {
		log.Fatal().Err(err).Msg("Failed to migrate database")
	}
	log.Info().Msg("Migrated database successfully")

	app := fiber.New()
	// TODO: add http logging middleware
	log.Info().Msg("Initializing fiber application")

	// Layers

	log.Info().Msgf("Starting server on port %d", config.PPlace.Port)
	log.Fatal().Err(app.Listen(fmt.Sprintf(":%d", config.PPlace.Port)))
}

func setupLogger() {
	zerolog.TimeFieldFormat = zerolog.TimeFormatUnix
	zerolog.ErrorStackMarshaler = pkgerrors.MarshalStack
	log.Logger = log.Output(zerolog.ConsoleWriter{Out: os.Stderr})
}
