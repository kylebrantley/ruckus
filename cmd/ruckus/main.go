package main

import (
	"net/http"
	"os"
	"os/signal"
	"time"

	"github.com/kylebrantley/ruckus/internal/executor"
	"github.com/rs/zerolog"
)

func main() {
	logger := zerolog.New(
		zerolog.ConsoleWriter{
			Out:        os.Stderr,
			TimeFormat: time.RFC3339,
		},
	).With().Timestamp().Caller().Logger()

	logger.Info().Msg("bringing da ruckus")

	c := make(chan os.Signal, 1)
	signal.Notify(c, os.Interrupt)

	// bodyReader := bytes.NewReader([]byte(`{"message":"hello"}`))
	request, err := http.NewRequest(http.MethodGet, "http://localhost:8080/hello", nil /* bodyReader */)
	if err != nil {
		// TODO: should log a message about potential config values when that is implemented
		logger.Fatal().Msg("failed to create request")
	}

	e := executor.Executor{Request: request, NumberOfRequests: 2, NumberOfThreads: 1, RequestTimeout: 10}

	go func() {
		<-c
		logger.Info().Msg("the ruckus was ended prematurely")
		e.Stop()
		os.Exit(0)
	}()

	e.Start()
	logger.Info().Msg("the ruckus has been brought")
}
