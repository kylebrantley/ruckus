package executor_test

import (
	"net/http"
	"net/http/httptest"
	"sync/atomic"
	"testing"

	"github.com/kylebrantley/ruckus/internal/executor"
	"github.com/kylebrantley/ruckus/internal/report"
	"github.com/stretchr/testify/assert"
)

type mockParser struct {
	finished chan bool
}

func newMockParser() mockParser {
	return mockParser{
		finished: make(chan bool, 1),
	}
}

func (m mockParser) Run() {
	m.finished <- true
}

func (m mockParser) Details() report.Details {
	return report.Details{}
}

func (m mockParser) Finished() <-chan bool {
	return m.finished
}

func TestExecutor_Start(t *testing.T) {
	scenarios := []struct {
		name             string
		numberOfRequests int
		numberOfThreads  int
	}{
		{
			name:             "test execution with 1 thread",
			numberOfRequests: 10,
			numberOfThreads:  1,
		},
		{
			name:             "test execution with 2 threads",
			numberOfRequests: 4,
			numberOfThreads:  2,
		},
	}
	for _, scenario := range scenarios {
		t.Run(
			scenario.name, func(t *testing.T) {
				requestCounter := int64(0)

				server := httptest.NewServer(
					http.HandlerFunc(
						func(w http.ResponseWriter, r *http.Request) {
							atomic.AddInt64(&requestCounter, int64(1))
						},
					),
				)
				defer server.Close()

				request, _ := http.NewRequest(http.MethodGet, server.URL, nil)

				parser := newMockParser()
				go parser.Run()

				e := executor.New(
					request,
					scenario.numberOfRequests,
					scenario.numberOfThreads,
					10,
					make(chan report.RequestResult, 10),
					parser,
				)

				e.Start()

				assert.Equal(t, int64(scenario.numberOfRequests), requestCounter)
			},
		)
	}
}
