package main

import (
	"os"

	"github.com/joho/godotenv"
	"github.com/rs/zerolog"

	"github.com/d-nery/catorce/pkg/bot"
)

func main() {
	logger := zerolog.New(zerolog.NewConsoleWriter()).
		With().
		Timestamp().
		Logger().
		Level(zerolog.InfoLevel)

	err := godotenv.Load()

	if err != nil {
		logger.Error().Err(err).Send()
		return
	}

	if _, v := os.LookupEnv("VERBOSE"); v {
		logger = logger.With().Caller().Logger().Level(zerolog.TraceLevel)
	}

	logger.Info().Msgf("Initializing bot... %s", bot.Version)

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
