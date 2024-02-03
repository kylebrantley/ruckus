package progressbar

import (
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestProgressbar_Print(t *testing.T) {
	scenarios := []struct {
		name     string
		total    int
		current  int
		expected string
	}{
		{
			name:     "test progress bar at 0%",
			total:    100,
			current:  0,
			expected: "\r[                                                  ] 0%",
		},
		{
			name:     "test progress bar at 50%",
			total:    100,
			current:  50,
			expected: "\r[■■■■■■■■■■■■■■■■■■■■■■■■■                         ] 50%",
		},
		{
			name:     "test progress bar at 100%",
			total:    100,
			current:  100,
			expected: "\r[■■■■■■■■■■■■■■■■■■■■■■■■■■■■■■■■■■■■■■■■■■■■■■■■■■] 100%\n",
		},
	}

	for _, scenario := range scenarios {
		t.Run(
			scenario.name, func(t *testing.T) {
				buf := strings.Builder{}

				p := New(scenario.total, WithWriter(&buf))
				p.Print(scenario.current)

				result := buf.String()
				assert.Equal(t, scenario.expected, result)
			},
		)
	}
}
