package main

import (
	"net/http"
	"os"
	"os/signal"
	"time"

	"github.com/kylebrantley/ruckus/internal/executor"
	"github.com/kylebrantley/ruckus/internal/report"
	"github.com/rs/zerolog"
)

// Arbitrary max channel size to reduce how much memory will be allocated.
// TODO: Do some tests and find an optimal value.
const maxChannelSize = 100

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

	results := make(chan report.RequestResult, min(maxChannelSize, 10))
	r := report.New(results)
	e := executor.New(request, 10, 2, 10, results, r)

	go func() {
		<-c
		logger.Info().Msg("the ruckus was ended prematurely")
		e.Stop()
		os.Exit(0)
	}()

	e.Start()

	logger.Info().Interface("details", r.Details()).Msg("the ruckus has been brought")
}
