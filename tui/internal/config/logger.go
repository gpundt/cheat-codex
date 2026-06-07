package config

import (
	"os"
	"time"

	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
)

// Helper function to intialize log string Format
// Sets up stdout logger
func InitializeLogger() {
	zerolog.TimeFieldFormat = time.RFC3339
	zerolog.SetGlobalLevel(zerolog.DebugLevel)

	consoleWriter := zerolog.ConsoleWriter{
		Out:        os.Stdout,
		TimeFormat: "15:04:05",
		NoColor:    false,
	}

	log.Logger = zerolog.New(consoleWriter).
		With().
		Timestamp().
		Logger()

	log.Debug().Msg("Logger initialized")
}
