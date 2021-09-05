package main

import (
	"logging/logger"
	"os"

	"github.com/rs/zerolog"
)

func main() {
	if len(os.Args) < 2 {
		panic("argument missing")
	}
	level := os.Args[1]
	// Initialize global logger
	switch level {
	case "info":
		logger.Log = logger.NewLogger(zerolog.InfoLevel)
	case "debug":
		logger.Log = logger.NewLogger(zerolog.DebugLevel)
	case "trace":
		logger.Log = logger.NewLogger(zerolog.TraceLevel)
	case "error":
		logger.Log = logger.NewLogger(zerolog.ErrorLevel)
	case "fatal":
		logger.Log = logger.NewLogger(zerolog.FatalLevel)
		logger.Log.Fatal().Msgf("setting timeout to '%d'", 10)
	default:
		panic("can not determine log level")
	}
	logger.Log.Info().Msgf("setting timeout to '%d'", 10)
	logger.Log.Debug().Msgf("setting timeout to '%d'", 10)
	logger.Log.Trace().Msgf("setting timeout to '%d'", 10)
	logger.Log.Error().Msgf("setting timeout to '%d'", 10)
}
