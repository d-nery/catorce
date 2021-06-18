package main

import (
	"math/rand"
	"os"
	"time"

	"github.com/joho/godotenv"
	"github.com/rs/zerolog"

	"github.com/d-nery/catorce/pkg/bot"
)

func main() {
	logger := zerolog.New(zerolog.NewConsoleWriter()).
		With().
		Timestamp().
		Caller().
		Logger().
		Level(zerolog.InfoLevel)

	err := godotenv.Load()

	if err != nil {
		logger.Error().Err(err).Send()
		return
	}

	logger.Info().Msg("Initializing bot...")

	rand.Seed(time.Now().UnixNano())

	b, err := bot.New(os.Getenv("TELEGRAM_TOKEN"), logger)

	if err != nil {
		logger.Error().Err(err).Send()
		return
	}

	b.Load()
	b.SetupHandlers()

	// b.Dump()

	b.Start()
}
