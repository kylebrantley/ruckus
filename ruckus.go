package main

import (
	"fmt"
	"io/ioutil"
	"net/http"
	"os"
	"os/signal"
	"sync/atomic"
	"time"

	"github.com/kylebrantley/ruckus/internal/config"
	"github.com/kylebrantley/ruckus/internal/executor"
	"github.com/kylebrantley/ruckus/internal/progressbar"
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

	cfg, err := config.New("config.yaml")
	if err != nil {
		logger.Fatal().Err(err).Msg("failed to load config")
	}

	c := make(chan os.Signal, 1)
	signal.Notify(c, os.Interrupt)

	// bodyReader := bytes.NewReader([]byte(`{"message":"hello"}`))
	request, err := http.NewRequest(cfg.Request.Method, cfg.Request.URL, nil /* bodyReader */)
	if err != nil {
		// TODO: should log a message about potential config values when that is implemented
		logger.Fatal().Msg("failed to create request")
	}

	results := make(chan report.RequestResult, min(maxChannelSize, cfg.NumberOfRequests))
	r := report.New(cfg.NumberOfRequests, results)
	e := executor.New(request, cfg.NumberOfRequests, cfg.MaxConcurrentRequests, cfg.Request.Timeout, results, r)

	go func() {
		<-c
		logger.Info().Msg("the ruckus was ended prematurely")
		e.Stop()
		os.Exit(0)
	}()

	requestCounter := int32(0)

	content, _ := ioutil.ReadFile("ascii.txt")
	fmt.Println(string(content))

	progressBar := progressbar.New(cfg.NumberOfRequests)
	go func() {
		fmt.Println("")
		fmt.Println("")
		for range e.ProcessedChannel {
			atomic.AddInt32(&requestCounter, 1)
			progressBar.Print(int(atomic.LoadInt32(&requestCounter)))
		}
	}()

	fmt.Printf("URL: %s\n", cfg.Request.URL)
	fmt.Printf("method: %s\n", cfg.Request.Method)
	fmt.Printf("total number requests: %d\n", cfg.NumberOfRequests)

	e.Start()

	logger.Info().Interface("details", r.Details()).Msg("the ruckus has been brought")
}
