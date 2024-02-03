package config

import (
	"fmt"
	"os"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestNew(t *testing.T) {
	scenarios := []struct {
		name     string
		config   string
		expected Config
	}{
		{
			name: "valid config",
			config: `
request:
  url: https://example.com
  method: GET
  body: '{"hello": "world"}'
  headers:
      content-type: application/json
  timeout: 1
numberOfRequests: 1
maxConcurrentRequests: 1
`,
			expected: Config{
				Request: Request{
					URL:     "https://example.com",
					Method:  "GET",
					Body:    `{"hello": "world"}`,
					Headers: map[string]string{"content-type": "application/json"},
					Timeout: 1,
				},
				NumberOfRequests:      1,
				MaxConcurrentRequests: 1,
			},
		},
	}

	for _, scenario := range scenarios {
		t.Run(
			scenario.name, func(t *testing.T) {
				tempFile, err := os.CreateTemp(".", "test.config.*.yaml")
				if err != nil {
					t.Fatal(err)
				}

				defer os.Remove(tempFile.Name())

				_, err = tempFile.Write([]byte(scenario.config))
				require.NoError(t, err)

				err = tempFile.Close()
				require.NoError(t, err)

				actual, err := New(fmt.Sprintf("./%s", tempFile.Name()))
				require.NoError(t, err)
				assert.Equal(t, scenario.expected, actual)
			},
		)
	}
}
