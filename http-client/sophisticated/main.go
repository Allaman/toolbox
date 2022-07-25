package main

import (
	aux "github.com/allaman/toolbox/http-client/auxiliary"
	"github.com/allaman/toolbox/http-client/sophisticated/logger"
	"github.com/rs/zerolog"
)

// C is the client instance
var C *Client

func main() {
	debug := false
	verbose := true
	timeout := 30
	// Initialize global logger
	if debug {
		logger.Log = logger.NewLogger(zerolog.DebugLevel)
	} else if verbose {
		logger.Log = logger.NewLogger(zerolog.TraceLevel)
	} else {
		logger.Log = logger.NewLogger(zerolog.InfoLevel)
	}

	if logger.Log.Trace().Enabled() {
		C = CreateClient(
			withTimeout(timeout),
			withDumpingEnabled(),
			withBaseURL("https://postman-echo.com"),
		)
	} else {
		C = CreateClient(
			withTimeout(timeout),
			withBaseURL("https://postman-echo.com"),
		)
	}

	req, err := C.NewRequest("POST", "/post", map[string]string{"foo": "bar"})
	if err != nil {
		logger.Log.Fatal().Msgf("error occured while creating request: %v", err)
	}
	resp, err := C.Do(req)
	if err != nil {
		logger.Log.Fatal().Msgf("error while performing request: %v", err)
	}
	defer resp.Body.Close()

	aux.PrintResponse(resp)
}
