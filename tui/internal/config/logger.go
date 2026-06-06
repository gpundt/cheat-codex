package config

import (
	"fmt"
	"os"
	"strconv"
	"time"

	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
)

// Helper function to intialize log string Format
// Sets up file logger and stdout logger
func InitializeLogger(component string, filepathsStruct Config.Filepaths) {
	zerolog.TimeFieldFormat = time.RFC3339
	zerolog.SetGlobalLevel(zerolog.DebugLevel)

	date := time.Now().Format("2006-01-02")

	consoleWriter := zerolog.ConsoleWriter{
		Out:        os.Stdout,
		TimeFormat: "15:04:05",
		NoColor:    false,
	}

	log.Logger = zerolog.New(consoleWriter).
		With().
		Timestamp().
		Logger()

	log.Info().Msg("Logger initialized")
}
