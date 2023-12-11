package config

import (
	"fmt"
	"os"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestNew(t *testing.T) {
	tempFile, err := os.CreateTemp(".", "test.config.*.yaml")
	if err != nil {
		t.Fatal(err)
	}

	defer os.Remove(tempFile.Name())

	text := []byte(`
request:
  url: https://example.com
  method: GET
  body: '{"hello": "world"}'
  headers:
      content-type: application/json
  timeout: 1
numberOfRequests: 1
maxConcurrentRequests: 1
`)

	_, err = tempFile.Write(text)
	require.NoError(t, err)

	err = tempFile.Close()
	require.NoError(t, err)

	expected := Config{
		Request: Request{
			URL:     "https://example.com",
			Method:  "GET",
			Body:    `{"hello": "world"}`,
			Headers: map[string]string{"content-type": "application/json"},
			Timeout: 1,
		},
		NumberOfRequests:      1,
		MaxConcurrentRequests: 1,
	}

	actual, err := New(fmt.Sprintf("./%s", tempFile.Name()))
	require.NoError(t, err)
	assert.Equal(t, expected, actual)
}
