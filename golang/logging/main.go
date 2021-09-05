package main

import (
	"logging/logger"

	"github.com/rs/zerolog"
)

var (
	info  = true
	debug = false
	trace = false
	error = false
	fatal = false
)

func main() {
	// Initialize global logger
	if info {
		logger.Log = logger.NewLogger(zerolog.InfoLevel)
	} else if debug {
		logger.Log = logger.NewLogger(zerolog.DebugLevel)
	} else if trace {
		logger.Log = logger.NewLogger(zerolog.TraceLevel)
	} else if error {
		logger.Log = logger.NewLogger(zerolog.ErrorLevel)
	} else if fatal {
		logger.Log = logger.NewLogger(zerolog.FatalLevel)
	} else {
		panic("can not determine log level")
	}
	logger.Log.Info().Msgf("setting timeout to '%d'", 10)
	logger.Log.Debug().Msgf("setting timeout to '%d'", 10)
	logger.Log.Trace().Msgf("setting timeout to '%d'", 10)
	logger.Log.Error().Msgf("setting timeout to '%d'", 10)
	logger.Log.Fatal().Msgf("setting timeout to '%d'", 10)
}
